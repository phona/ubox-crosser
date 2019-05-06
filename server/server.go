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

func (p *ProxyServer) handleExposerConn(src net.Conn) {
	if p.exposerPass == "" {
		dst, err := p.controller.GetConn()
		if err != nil {
			log.Error(err)
			src.Close()
			return
		}
		go drillingTunnel(src, dst)
		return
	}

	if p.cipher != nil {
		src = ss.NewConn(src, p.cipher.Copy())
	}

	var reqMsg message.Message
	coordinator := connector.AsCoordinator(src)
	var errFunc = func(err error) {
		var respMsg message.ResultMessage
		respMsg.Result = message.FAILED
		buf, _ := json.Marshal(respMsg)
		coordinator.SendMsg(string(buf))
		src.Close()
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
			src.Close()
			log.Errorf("Error handling connection in proxy server: %s", err)
		}

		var respMsg message.ResultMessage
		respMsg.Result = message.SUCCESS
		buf, _ := json.Marshal(respMsg)
		if err := coordinator.SendMsg(string(buf)); err != nil {
			simpleErrHandle(err)
		} else if p.controller == nil {
			simpleErrHandle(fmt.Errorf("The controller of proxy server is null."))
		} else if dst, err := p.controller.GetConn(); err != nil {
			simpleErrHandle(err)
		} else {
			go drillingTunnel(src, dst)
		}
	}
}

func (p *ProxyServer) runController() {
	p.controller = NewController(p.controllerAddr, p.controllerPass, p.cipher.Copy())
	p.controller.Run()
}

func testSocks5Req(src, dst net.Conn) {
	buf := make([]byte, 10)
	dst.Write([]byte{5, 2, 0, 1})
	dst.Read(buf)
	fmt.Println(buf)
	src.Close()
	dst.Close()
}
