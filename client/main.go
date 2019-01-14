package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"io"
	"net"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/hashicorp/yamux"
	"github.com/jcelliott/lumber"
)

//go:generate ../tools/client.sh ../tools

type Config struct {
	serverAddr string
	localAddr  string
	tunnelId   string
	logFile    string
	logConsole bool
	logLevel   string
}

var config Config
var session *yamux.Session
var log lumber.Logger

func main() {
	// parse command line arguments
	parseArgs()
	// init logger
	initLogger()

	go func() {
		tlsConfig, err := getTlsConfig()
		if err != nil {
			log.Error("load certificates error: %v", err)
			os.Exit(1)
		}
		// connect to the server
		log.Debug("dial server %v", config.serverAddr)
		conn, err := tls.Dial("tcp", config.serverAddr, tlsConfig)
		if err != nil {
			log.Error("cannot connect to %v: %v", config.serverAddr, err)
			os.Exit(1)
			return
		}
		log.Info("connected to server %v", config.serverAddr)

		//send the client type C and identification
		handshake := "C," + config.tunnelId + "\n"
		log.Debug("send client id %v", handshake)
		_, err = io.WriteString(conn, handshake)
		if err != nil {
			log.Error("write client id error: %v", err)
			os.Exit(1)
			return
		}

		log.Debug("start yamux.Server()")
		session, err = yamux.Server(conn, nil)
		if err != nil {
			log.Error("start yamux.Server(): %v", err)
			os.Exit(1)
			return
		}

		log.Debug("waiting for server to create the control stream")
		sessionAcceptCh := make(chan bool, 1)
		go func() {
			controlStream, err := session.Accept()
			if err != nil {
				log.Error("error accept control stream")
				session.Close()
				conn.Close()
				os.Exit(1)
				return
			}
			log.Debug("control stream created")
			go readServerMessages(controlStream)
			sessionAcceptCh <- true
		}()

		select {
		case <-sessionAcceptCh:
		case <-time.After(time.Second * 10):
			log.Error("timeout waiting for control stream to be created")
			session.Close()
			conn.Close()
			os.Exit(1)
			return
		}

		for {
			// Accept a stream
			log.Debug("waiting for tunnel streams session.Accept()")
			stream, err := session.Accept()
			if err != nil {
				log.Error("session.Accept() ERROR: %v", err)
				os.Exit(1)
				return
			}
			log.Debug("stream accepted - start forward")
			go forward(stream, config.localAddr)
		}
	}()

	// wait for CTRL + C
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	// TODO: cleanup connections here
	session.Close()
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
	serverAddr := flag.String("server", "", "Server address")
	localAddr := flag.String("local", "", "local address")
	tunnelId := flag.String("id", "", "the id of the tunnel")
	logFile := flag.String("log-file", "", "Full log file path")
	logConsole := flag.Bool("log-console", true, "Send logs to console")
	logLevel := flag.String("log-level", "TRACE", "Set the file log level")
	flag.Parse()
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
			session.Close()
			os.Exit(1)
			return
		}
		s := string(buf[0:n])
		log.Debug("message received:%v", s)
	}
}

func forward(conn net.Conn, localAddress string) {
	client, err := net.Dial("tcp", localAddress)
	if err != nil {
		log.Error("dial local address failed %v: %v", localAddress, err)
		return
	}
	log.Info("connected to %v start forward", localAddress)
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
