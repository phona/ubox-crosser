package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"ubox-crosser/models/config"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

// for opening a listener to proxy request
type ProxyServer struct {
	// generated from client

	dispatcher  *connector.Dispatcher
	controllers map[string]*Controller
	errs        chan error

	context map[string]config.ServerConfig
	// exposers    map[string]*Exposer
}

func NewProxyServer(configs map[string]config.ServerConfig) *ProxyServer {
	total := len(configs)
	dispatcher := connector.NewDispatcher(uint64(total))
	listenedAddr := make([]string, 0, total)
	server := &ProxyServer{
		dispatcher:  dispatcher,
		controllers: make(map[string]*Controller, total),
		errs:        make(chan error, 10),
		context:     configs,
	}

	for _, config_ := range configs {
		go server.initWorker(listenedAddr, config_)
	}
	return server
}

func (p *ProxyServer) initWorker(listenedAddr []string, serverConfig config.ServerConfig) {
	// deduplicate address
	for _, addr := range listenedAddr {
		if addr == serverConfig.Address {
			return
		}
	}

	var cipher *ss.Cipher
	if serverConfig.Method != "" {
		if err := ss.CheckCipherMethod(serverConfig.Method); err != nil {
			p.errs <- err
			return
		} else if cipher, err = ss.NewCipher(serverConfig.Method, serverConfig.Key); err != nil {
			p.errs <- err
			return
		}
	}

	if l, err := net.Listen("tcp", serverConfig.Address); err != nil {
		p.errs <- err
	} else {
		p.dispatcher.Add(connector.NewCipherListener(l, cipher))
	}
}

func (p *ProxyServer) Err() error {
	select {
	case err := <-p.errs:
		return err
	default:
		return p.dispatcher.Err()
	}
}

func (p *ProxyServer) Process() {
	// log.Infof("Exposer listen on %s, Controller listen on %s", p.exposerAddr, p.controllerAddr)
	for {
		conn := p.dispatcher.Conn()
		go p.handleConnection(conn)
	}
}

func (p *ProxyServer) handleConnection(conn net.Conn) {
	coordinator := connector.AsCoordinator(conn)

	var reqMsg message.Message
	if content, err := coordinator.ReadMsg(); err != nil {
		p.errs <- err
	} else if err := json.Unmarshal([]byte(content), &reqMsg); err != nil {
		p.errs <- err
	} else {
		var handleErr = func(err error) {
			p.errs <- err
			conn.Close()
		}
		switch reqMsg.Type {
		case message.LOGIN:
			if context, ok := p.context[reqMsg.ServeName]; !ok {
				handleErr(fmt.Errorf("Unknown serve %s were received", reqMsg.ServeName))
			} else if reqMsg.Password == context.LoginPass {
				controller := NewController(coordinator)
				p.controllers[reqMsg.ServeName] = controller
				controller.daemonize()
			} else {
				handleErr(fmt.Errorf("Error Receving invalid login password request %+v", reqMsg))
			}
		case message.GEN_WORKER:
			if controller, ok := p.controllers[reqMsg.ServeName]; !ok {
				handleErr(fmt.Errorf("Controller for %s does not alive", reqMsg.ServeName))
			} else {
				controller.HandleConnection(conn)
			}
		case message.AUTHENTICATION:
			p.handleAuthRequest(reqMsg.ServeName, reqMsg.Password, coordinator)
		default:
			handleErr(fmt.Errorf("Unknown type %s were received", reqMsg.Type))
		}
	}
}

func (p *ProxyServer) handleAuthRequest(serveName, authPass string, coordinator *connector.Coordinator) {
	var errFunc = func(err error) {
		var respMsg message.ResultMessage
		respMsg.Result = message.FAILED
		buf, _ := json.Marshal(respMsg)
		coordinator.SendMsg(string(buf))
		coordinator.Close()
		log.Errorf("Error handling connection in proxy server: %s", err)
	}

	if context, ok := p.context[serveName]; !ok {
		errFunc(fmt.Errorf("Unknown serve %s were received", serveName))
	} else if controller, ok := p.controllers[serveName]; !ok {
		errFunc(fmt.Errorf("Controller for %s does not alive", serveName))
	} else {
		if authPass != context.AuthPass {
			errFunc(fmt.Errorf("Invalid password %s != %s", context.AuthPass, authPass))
		} else {
			var simpleErrHandle = func(err error) {
				coordinator.Close()
				log.Errorf("Error handling connection in proxy server: %s", err)
			}

			respMsg := message.ResultMessage{message.SUCCESS}
			buf, _ := json.Marshal(respMsg)
			if err := coordinator.SendMsg(string(buf)); err != nil {
				simpleErrHandle(err)
			} else if workConn, err := controller.GetConn(); err != nil {
				simpleErrHandle(err)
			} else {
				go drillingTunnel(coordinator.Conn, workConn)
			}
		}
	}
}

func testSocks5Req(src, dst net.Conn) {
	buf := make([]byte, 10)
	dst.Write([]byte{5, 2, 0, 1})
	dst.Read(buf)
	fmt.Println(buf)
	src.Close()
	dst.Close()
}
