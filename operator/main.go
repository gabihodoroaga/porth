package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"io"
	"net"
	"os"

	"github.com/hashicorp/yamux"
	"github.com/jcelliott/lumber"
)

//go:generate ../tools/operator.sh ../tools

type Config struct {
	serverAddr string
	localAddr  string
	tunnelId   string
	logFile    string
	logConsole bool
	logLevel   string
}

var config Config
var log lumber.Logger

func main() {
	// parse command line arguments
	parseArgs()

	// init the logger
	initLogger()

	tslConfig, err := getTlsConfig()
	if err != nil {
		log.Error("error load certificates %v", err)
		os.Exit(1)
	}
	log.Debug("dial server %v", config.serverAddr)
	conn, err := tls.Dial("tcp", config.serverAddr, tslConfig)
	if err != nil {
		log.Error("error dial server %v: %v", config.serverAddr, err)
		os.Exit(1)
	}

	log.Info("connected to server %v", config.serverAddr)
	handshake := "O," + config.tunnelId + "\n"
	log.Debug("send operator handshake: %v", handshake)
	_, err = io.WriteString(conn, handshake)
	if err != nil {
		log.Error("error write client id")
		os.Exit(1)
		return
	}

	log.Debug("waiting for server response")
	bufc := bufio.NewReader(conn)
	// TODO: this should be done with timeout
	line, err := bufc.ReadString('\n')
	if err != nil {
		log.Error("read client id error:%v", err)
		os.Exit(1)
		return
	}

	if line != "OK\n" {
		log.Error("server error received %v", line)
		os.Exit(1)
		return
	}

	log.Debug("server replied OK")
	// Setup client side yahmux
	session, err := yamux.Client(conn, nil)
	if err != nil {
		log.Error("error create yamux session yamux.Client:", err)
		os.Exit(1)
	}

	control_stream, err := session.Open()
	if err != nil {
		log.Error("error open control stream : %v", err)
		conn.Close()
		os.Exit(1)
	}

	log.Debug("control stream created")

	go readServerMessages(control_stream)

	listener, err := net.Listen("tcp", config.localAddr)
	if err != nil {
		log.Error("error create listener for address %v", config.localAddr, err)
		os.Exit(1)
	}

	for {
		log.Debug("waiting for connection on %v", config.localAddr)
		local_conn, err := listener.Accept()
		if err != nil {
			log.Error("ERROR: failed to accept listener: %v", err)
		}

		log.Debug("accepted connection %v\n", config.localAddr)
		// Open a new stream
		log.Debug("try to open server a stream...")
		stream, err := session.Open()
		if err != nil {
			log.Error("session.Open() error: %v", err)
			local_conn.Close()
			continue
		}
		log.Debug("server stream opened - start forward")
		go forward(local_conn, stream)
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

func parseArgs() {
	serverAddr := flag.String("server", "", "the server address")
	localAddr := flag.String("local", "", "the local address")
	tunnelId := flag.String("id", "", "the id of the tunnel")
	logFile := flag.String("log-file", "", "Full log file path")
	logConsole := flag.Bool("log-console", true, "Send logs to console")
	logLevel := flag.String("log-level", "TRACE", "Set the file log level")
	flag.Parse()

	if *serverAddr == "" || *localAddr == "" || *tunnelId == "" {
		flag.Usage()
		os.Exit(1)
		return
	}

	config = Config{
		serverAddr: *serverAddr,
		localAddr:  *localAddr,
		tunnelId:   *tunnelId,
		logFile:    *logFile,
		logConsole: *logConsole,
		logLevel:   *logLevel,
	}
}

func getTlsConfig() (*tls.Config, error) {
	// read certificates
	cert_client, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootCA))
	if !ok {
		return nil, errors.New("error load rooCA")
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert_client},
		InsecureSkipVerify: false,
		RootCAs:            roots,
	}, nil
}

func readServerMessages(conn net.Conn) {
	var buf [64]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			log.Error("error read on control stream, we must exit or retry connect")
			// TODO: kill all local connection
			os.Exit(1)
			return
		}
		s := string(buf[0:n])
		log.Debug("message received:%v", s)
	}
}

func forward(conn net.Conn, client net.Conn) {
	log.Debug("start copy from %s->%s to %s->%s", conn.RemoteAddr(), conn.LocalAddr(), client.RemoteAddr(), client.LocalAddr())
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()

	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}
