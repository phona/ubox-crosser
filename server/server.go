package server

import (
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"os"
)

// for opening a listener to proxy request
type ProxyServer struct {
	// generated from client
	listener   *net.Listener
	controller *Controller

	cipher *ss.Cipher
}

func NewProxyServer(cipher *ss.Cipher) *ProxyServer {
	return &ProxyServer{cipher: cipher}
}

func (p *ProxyServer) Listen(southAddress, northAddress string) {
	log.Infof("South bridge listen on %s, North bridge listen on %s", southAddress, northAddress)
	go p.serve(southAddress)
	p.openController(northAddress)
}

func (p *ProxyServer) serve(address string) {
	if listener, err := net.Listen("tcp", address); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	} else {
		p.listener = &listener
		for {
			rawConn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				continue
			}
			log.Info("get a new request")
			go p.pipe(rawConn)
		}
	}
}

func (p *ProxyServer) openController(address string) {
	p.controller = NewController(address, p.cipher)
	p.controller.Run()
}

func (p *ProxyServer) pipe(conn net.Conn) {
	workConn, err := p.controller.GetConn()
	log.Debugf("Pipe between request connection and work connection, %s -> %s", conn.RemoteAddr().String(), workConn.RemoteAddr().String())
	if err != nil {
		log.Println("Listener for incoming connections from client closed")
		log.Error("Error pipe:", err)
	} else {
		go ss.PipeThenClose(conn, workConn)
		ss.PipeThenClose(workConn, conn)
	}
}
