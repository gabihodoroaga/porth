package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hashicorp/yamux"
	"github.com/jcelliott/lumber"
)

//go:generate ../tools/server.sh ../tools

type Config struct {
	serverAddr string
	logFile    string
	logConsole bool
	logLevel   string
	httpAddr   string
}

type ClientList map[string]*ClientPair

type Client struct {
	id        string
	conn      net.Conn
	session   *yamux.Session
	connected bool
}

type Operator struct {
	conn    net.Conn
	session *yamux.Session
}

type ClientPair struct {
	cl *Client
	op []Operator
}

var config Config
var clients ClientList = ClientList{}

var log lumber.Logger

func main() {
	// parse command line arguments
	parseArgs()
	// init logger
	initLogger()
	// load certificates
	cert_server, err := tls.X509KeyPair([]byte(serverCert), []byte(serverKey))
	if err != nil {
		log.Error("load server ceritificates error: %v", err)
		os.Exit(1)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootCA))
	if !ok {
		panic("failed to parse root certificate")
	}
	// start tunnel server
	go startTunnelServer(cert_server, roots)
	// start http server
	go startHttpServer()

	// wait for crtl+c or any process interups
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Press CTRL+C to stop the program")

	// wait for signal intrerupt here
	<-done
}

func parseArgs() {
	serverAddr := flag.String("addr", ":2671", "server address and port")
	logFile := flag.String("log-file", "", "Full log file path")
	logConsole := flag.Bool("log-console", true, "Send logs to console")
	logLevel := flag.String("log-level", "TRACE", "Set the file log level")
	httpAddr := flag.String("http-addr", ":2672", "http server addredd")
	flag.Parse()
	config = Config{
		serverAddr: *serverAddr,
		logFile:    *logFile,
		logConsole: *logConsole,
		logLevel:   *logLevel,
		httpAddr:   *httpAddr,
	}
}

func initLogger() {
	var logger = lumber.NewMultiLogger()

	if config.logFile != "" {
		log, err := lumber.NewRotateLogger(config.logFile, 5000, 9)
		if err != nil {
			panic(err.Error())
		}
		log.Level(lumber.LvlInt(config.logLevel))
		logger.AddLoggers(log)
	}

	if config.logConsole {
		log := lumber.NewConsoleLogger(lumber.TRACE)
		logger.AddLoggers(log)
	}

	log = logger
}

func startHttpServer() {
	apiserver := &ApiServer{}
	log.Info("start http server on %v", config.httpAddr)
	err := http.ListenAndServe(config.httpAddr, apiserver)
	if err != nil {
		log.Error("failed to start HTTP: %s", err)
		os.Exit(1)
		return
	}
	log.Info("http server stoped %v", config.httpAddr)
}

func startTunnelServer(cert tls.Certificate, roots *x509.CertPool) {
	// setup ssl
	tlsConfig := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          roots,
	}

	// create server listner
	listener, err := tls.Listen("tcp", config.serverAddr, &tlsConfig)
	if err != nil {
		log.Fatal("Failed to setup listener: %v", err)
	}

	// start the server
	log.Info("start tunnel server on address %v", config.serverAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			// TODO handle accept error - fatal=exit process / continue
			log.Fatal("accept connection error: %v", err)
		}
		log.Info("[%v] connection accepted", conn.RemoteAddr().String())
		go handdleConnection(conn)
	}
}

func forward(conn net.Conn, client net.Conn) {

	go func() {
		defer client.Close()
		defer conn.Close()
		log.Debug("start copy from %s->%s to %s->%s", conn.RemoteAddr(), conn.LocalAddr(), client.RemoteAddr(), client.LocalAddr())
		nr, err := io.Copy(client, conn)
		log.Debug("stop copy from %s->%s to %s->%s(bytes %v)", conn.RemoteAddr(), conn.LocalAddr(), client.RemoteAddr(), client.LocalAddr(), nr)
		if err != nil {
			log.Debug("error copy from %s->%s to %s->%s (%v)", conn.RemoteAddr(), conn.LocalAddr(), client.RemoteAddr(), client.LocalAddr(), err)
		}
	}()

	go func() {
		defer client.Close()
		defer conn.Close()
		log.Debug("start copy from %s->%s to %s->%s", client.RemoteAddr(), client.LocalAddr(), conn.RemoteAddr(), conn.LocalAddr())
		nr, err := io.Copy(conn, client)
		log.Debug("stop copy from %s->%s to %s->%s(bytes: %v)", client.RemoteAddr(), client.LocalAddr(), conn.RemoteAddr(), conn.LocalAddr(), nr)
		if err != nil {
			log.Debug("error copy from %s->%s to %s->%s(bytes: %v)", client.RemoteAddr(), client.LocalAddr(), conn.RemoteAddr(), conn.LocalAddr(), err)
		}
	}()
}

func handdleConnection(conn net.Conn) {
	var logContext = fmt.Sprintf("[%v]", conn.RemoteAddr().String())
	// TODO: read client id and type - this should be done with timeout
	log.Debug("%v waiting for client type and id", logContext)
	bufc := bufio.NewReader(conn)
	client_address, err := bufc.ReadString('\n')
	if err != nil {
		log.Error("read client id error:%v", err)
		conn.Close()
		return
	}

	client_parts := strings.Split(client_address, ",")
	client_type := client_parts[0]
	tunnel_id := client_parts[1][:len(client_parts[1])-1]

	switch client_type {
	default:
		log.Error("%v invalid handshake received: %v", logContext, client_address)
		conn.Close()
		return
	case "C":
		handleClient(conn, tunnel_id)
	case "O":
		handleOperator(conn, tunnel_id)
	}
}

func handleClient(conn net.Conn, tunnel_id string) {
	// Setup client side of yamux
	var logContext = fmt.Sprintf("[%v][%v][%v]", conn.RemoteAddr().String(), "C", tunnel_id)
	log.Debug("%v begin client_session create", logContext)
	client_session, err := yamux.Client(conn, nil)
	if err != nil {
		log.Error("%v create yamux client error:%v", logContext, err)
		conn.Close()
		return
	}

	// open stream here dows not return error if the connection is broken
	// it olny checks if the config is ok and the maximum number of streams is reached
	control_client_stream, err := client_session.Open()
	if err != nil {
		log.Error("%v error open client control stream:%v", logContext, err)
		client_session.Close()
		conn.Close()
		return
	}

	var cp *ClientPair
	if t, ok := clients[tunnel_id]; ok {
		cp = t
	} else {
		cp = &ClientPair{}
		clients[tunnel_id] = cp
	}

	cl := &Client{id: tunnel_id, session: client_session, conn: conn, connected: true}
	cp.cl = cl
	log.Info("%v tunnel_id %v was added to the tunnel list", logContext, tunnel_id)
	go readClientMessages(control_client_stream, cp, logContext)
}

func checkAndRemoveClientPair(cp *ClientPair) {
	if cp.cl.connected == false && len(cp.op) == 0 {
		delete(clients, cp.cl.id)
	}
}

func readClientMessages(conn net.Conn, cp *ClientPair, logContext string) {
	var buf [64]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			log.Debug("%v error read on client control stream:%v", logContext, err)
			cp.cl.connected = false
			checkAndRemoveClientPair(cp)
			return
		}
		s := string(buf[0:n])
		log.Debug("%v, message received:%v", logContext, s)
	}
}

func removeOperator(s []Operator, r Operator) []Operator {
	// TODO: mutex -> this list can be modified by multiple threads
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func readOperatorMessages(conn net.Conn, logContext string) {
	var buf [64]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			log.Debug("%v error read on operator control stream:%v", logContext, err)
			return
		}
		s := string(buf[0:n])
		log.Debug("%v message received:%v", logContext, s)
	}
}

func handleOperator(conn net.Conn, tunnel_id string) {
	var logContext = fmt.Sprintf("[%v][%v][%v]", conn.RemoteAddr().String(), "O", tunnel_id)
	// search for client_id
	cp, ok := clients[tunnel_id]
	// if client_id does not exists then drop connection
	if ok != true {
		log.Info("%v client id not found %v", logContext, tunnel_id)
		_, err := io.WriteString(conn, "client with id "+tunnel_id+" not found\n")
		if err != nil {
			log.Error("%v error write response to operator: %v", logContext, err)
		}
		conn.Close()
		return
	}

	_, err := io.WriteString(conn, "OK\n")
	if err != nil {
		log.Error("%v error write response to operator: %v", logContext, err)
		conn.Close()
		return
	}

	log.Debug("%v open operator server session", logContext)
	operator_session, err := yamux.Server(conn, nil)
	if err != nil {
		log.Error("%v error create operator_session yamux.Server: %v", logContext, err)
		conn.Close()
		return
	}

	log.Debug("%v waiting for operator control stream", logContext)
	operator_control_stream, err := operator_session.Accept()
	if err != nil {
		log.Error("%v error accept operator_control_stream: %v", logContext, err)
		conn.Close()
		return
	}
	// start reading messages from operator connection
	go readOperatorMessages(operator_control_stream, logContext)

	op := Operator{conn: conn, session: operator_session}
	cp.op = append(cp.op, op)

	defer func() {
		conn.Close()
		cp.op = removeOperator(cp.op, op)
		checkAndRemoveClientPair(cp)
	}()

	for {
		log.Debug("%v waiting for operator streams operator_session.Accept()", logContext)
		operator_stream, err := operator_session.Accept()
		if err != nil {
			log.Error("%v error operator_session.Accept(): %v", logContext, err)
			return
		}

		log.Debug("%v operator_stream accepted", logContext)
		client_stream, err := cp.cl.session.Open()
		if err != nil {
			log.Error("%v client_stream session.Open() error: %v", logContext, err)
			operator_stream.Close()
			continue
		}
		log.Debug("client_stream created")
		go forward(operator_stream, client_stream)
	}
}
