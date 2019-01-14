package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/hashicorp/yamux"
)

//go:generate ../tools/operator.sh ../tools

type Config struct {
	serverAddr string
	localAddr  string
	tunnelId   string
}

var config Config

func parseArgs() {
	serverAddr := flag.String("server", "", "the server address")
	localAddr := flag.String("local", "", "the local address")
	tunnelId := flag.String("id", "", "the id of the tunnel")
	flag.Parse()

	config = Config{
		serverAddr: *serverAddr,
		localAddr:  *localAddr,
		tunnelId:   *tunnelId,
	}
}

func main() {

	parseArgs()
	tslConfig, err := getTlsConfig()
	if err != nil {
		log.Printf("error load certificates %v", err)
		os.Exit(1)
	}
	log.Printf("dial server %v", config.serverAddr)
	conn, err := tls.Dial("tcp", config.serverAddr, tslConfig)
	if err != nil {
		log.Printf("error dial server %v: %v", config.serverAddr, err)
		os.Exit(1)
	}

	log.Printf("connected to server %v", config.serverAddr)
	handshake := "O," + config.tunnelId + "\n"
	log.Printf("send operator handshake: %v", handshake)
	_, err = io.WriteString(conn, handshake)
	if err != nil {
		log.Printf("error write client id")
		os.Exit(1)
		return
	}

	log.Printf("waiting for server response")
	bufc := bufio.NewReader(conn)
	// TODO: this should be done with timeout
	line, err := bufc.ReadString('\n')
	if err != nil {
		log.Printf("read client id error:%v", err)
		os.Exit(1)
		return
	}

	if line != "OK\n" {
		log.Printf("server error received %v", line)
		os.Exit(1)
		return
	}

	// Setup client side yahmux
	session, err := yamux.Client(conn, nil)
	if err != nil {
		log.Printf("error create yamux session yamux.Client:", err)
		os.Exit(1)
	}

	control_stream, err := session.Open()
	if err != nil {
		log.Printf("error open control stream : %v", err)
		conn.Close()
		return
	}

	go readServerMessages(control_stream)

	listener, err := net.Listen("tcp", config.localAddr)
	if err != nil {
		log.Printf("error create listener for address %v", config.localAddr, err)
		os.Exit(1)
	}

	for {
		log.Printf("waiting for connection on %v", config.localAddr)
		local_conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}

		log.Printf("accepted connection %v\n", config.localAddr)
		// Open a new stream
		log.Printf("try to open server a stream...")
		stream, err := session.Open()
		if err != nil {
			log.Printf("session.Open() error: %v", err)
			local_conn.Close()
			return
		}
		log.Printf("server stream opened")
		go forward(local_conn, stream)
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
			log.Print("error read on control stream, we must exit or retry connect")
			// TODO: kill all local connection
			os.Exit(1)
			return
		}
		s := string(buf[0:n])
		log.Printf("message received:%v", s)
	}
}

func forward(conn net.Conn, client net.Conn) {
	log.Printf("start copy from %s->%s to %s->%s", conn.RemoteAddr(), conn.LocalAddr(), client.RemoteAddr(), client.LocalAddr())
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
