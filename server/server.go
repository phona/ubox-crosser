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
}

func NewProxyServer() *ProxyServer {
	return &ProxyServer{}
}

func (p *ProxyServer) Listen(southAddress, northAddress string) {
	log.Infof("South bridge listen on %s, North bridge listen on %s", southAddress, northAddress)
	go p.openSouth(southAddress)
	p.openNorth(northAddress)
}

func (p *ProxyServer) openSouth(address string) {
	if listener, err := net.Listen("tcp", address); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	} else {
		p.listener = &listener
		for {
			rawConn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				os.Exit(0)
			}
			log.Info("get a new request")
			go p.pipe(rawConn)
		}
	}
}

func (p *ProxyServer) openNorth(address string) {
	p.controller = NewController(address)
	p.controller.Run()
}

func (p *ProxyServer) pipe(conn net.Conn) {
	for {
		workConn, err := p.controller.GetConn()
		if err != nil {
			log.Println("Listener for incoming connections from client closed")
			log.Error("Error pipe:", err)
			return
		} else {
			go ss.PipeThenClose(conn, workConn)
			go ss.PipeThenClose(workConn, conn)
			break
		}
	}
}
