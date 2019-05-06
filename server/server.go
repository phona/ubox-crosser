package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"os"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

// for opening a listener to proxy request
type ProxyServer struct {
	// generated from client
	listener   *net.Listener
	controller *Controller

	exposerPass, controllerPass string
	exposerAddr, controllerAddr string

	cipher *ss.Cipher
}

func NewProxyServer(exposerAddr, exposerPass, controllerAddr, controllerPass string, cipher *ss.Cipher) *ProxyServer {
	return &ProxyServer{
		exposerPass:    exposerPass,
		exposerAddr:    exposerAddr,
		controllerPass: controllerPass,
		controllerAddr: controllerAddr,
		cipher:         cipher,
	}
}

func (p *ProxyServer) Run() {
	log.Infof("Exposer listen on %s, Controller listen on %s", p.exposerAddr, p.controllerAddr)
	go p.runExposer()
	p.runController()
}

func (p *ProxyServer) runExposer() {
	if listener, err := net.Listen("tcp", p.exposerAddr); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	} else {
		p.listener = &listener
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				continue
			}
			log.Info("get a new request")
			go p.handleExposerConn(conn)
		}
	}
}

func (p *ProxyServer) handleExposerConn(conn net.Conn) {
	newConn := conn
	if p.cipher != nil {
		newConn = ss.NewConn(newConn, p.cipher.Copy())
	}

	var reqMsg message.Message
	coordinator := connector.AsCoordinator(newConn)
	var errFunc = func(err error) {
		var respMsg message.ResultMessage
		respMsg.Result = message.FAILED
		buf, _ := json.Marshal(respMsg)
		coordinator.SendMsg(string(buf))
		conn.Close()
		log.Errorf("Error handling connection in proxy server: %s", err)
	}

	if content, err := coordinator.ReadMsg(); err != nil {
		errFunc(err)
	} else if err := json.Unmarshal([]byte(content), &reqMsg); err != nil {
		errFunc(err)
	} else if reqMsg.Type != message.LOGIN {
		errFunc(fmt.Errorf("Invalid type %s != %s", message.LOGIN, reqMsg.Type))
	} else if reqMsg.Password != p.exposerPass {
		errFunc(fmt.Errorf("Invalid password %s != %s", p.exposerPass, reqMsg.Password))
	} else {
		var simpleErrHandle = func(err error) {
			conn.Close()
			log.Errorf("Error handling connection in proxy server: %s", err)
		}
		var respMsg message.ResultMessage
		respMsg.Result = message.SUCCESS
		buf, _ := json.Marshal(respMsg)
		if err := coordinator.SendMsg(string(buf)); err != nil {
			simpleErrHandle(err)
		} else {
			//go p.pipe(newConn)
			go p.pipe(conn)
		}
	}
}

func (p *ProxyServer) runController() {
	var IsNeedEncrypt bool
	if p.exposerPass == "" {
		IsNeedEncrypt = true
	} else {
		IsNeedEncrypt = false
	}
	p.controller = NewController(p.controllerAddr, p.controllerPass, IsNeedEncrypt, p.cipher)
	p.controller.Run()
}

func (p *ProxyServer) pipe(conn net.Conn) {
	if p.controller == nil {
		log.Error("The controller of proxy server is null.")
		return
	}

	workConn, err := p.controller.GetConn()
	if err != nil {
		log.Error(err)
		conn.Close()
		return
	}
	//
	//if p.cipher != nil {
	//	workConn = ss.NewConn(workConn, p.cipher.Copy())
	//}

	log.Debugf("Pipe between request connection and work connection, %s -> %s", conn.RemoteAddr().String(), workConn.RemoteAddr().String())
	go pipeThenClose(workConn, conn)
	pipeThenClose(conn, workConn)
}
