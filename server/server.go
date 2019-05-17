package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"ubox-crosser/models/config"
	"ubox-crosser/models/errors"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

// for opening a listener to proxy request
type ProxyServer struct {
	// generated from client

	dispatcher  *connector.Dispatcher
	controllers map[string]*controller
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
		controllers: make(map[string]*controller, total),
		errs:        make(chan error, DEFAULT_SERVES),
		context:     configs,
	}

	for _, config_ := range configs {
		go server.initWorker(&listenedAddr, config_)
	}
	return server
}

func (p *ProxyServer) initWorker(pListenedAddr *[]string, serverConfig config.ServerConfig) {
	// deduplicate address
	listenedAddr := *pListenedAddr
	for _, addr := range listenedAddr {
		if addr == serverConfig.Address {
			return
		}
	}

	listenedAddr = append(listenedAddr, serverConfig.Address)
	*pListenedAddr = listenedAddr

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
		log.Infof("Add new listener on address %s", serverConfig.Address)
		p.dispatcher.Add(connector.NewCipherListener(l, cipher))
	}
}

func (p *ProxyServer) Err() error {
	select {
	case err := <-p.errs:
		return err
	case err := <-p.dispatcher.Errs:
		return err
	}
}

func (p *ProxyServer) Process() {
	for {
		conn := <-p.dispatcher.Conns
		log.Infoln("Received a new connection")
		go p.handleConnection(conn)
	}
}

func (p *ProxyServer) handleConnection(conn net.Conn) {
	log.Infof("Remote address %s connect to center server", conn.RemoteAddr().String())
	coordinator := connector.AsCoordinator(conn)

	var reqMsg message.Message
	if content, err := coordinator.ReadMsg(); err != nil {
		p.errs <- err
	} else if err := json.Unmarshal([]byte(content), &reqMsg); err != nil {
		p.errs <- err
	} else {
		log.Infof("Received content: %s", content)
		switch reqMsg.Type {
		case message.LOGIN:
			p.handleLoginRequest(reqMsg.ServeName, reqMsg.Password, coordinator)
		case message.GEN_WORKER:
			if controller, ok := p.controllers[reqMsg.ServeName]; !ok {
				p.handleConnErr(coordinator, fmt.Errorf("controller for %s does not alive", reqMsg.ServeName), errors.INVALID_SERVE_NAME)
			} else {
				controller.workConn <- conn
			}
		case message.AUTHENTICATION:
			p.handleAuthRequest(reqMsg.ServeName, reqMsg.Password, coordinator)
		default:
			p.handleConnErr(coordinator, fmt.Errorf("Unknown type %d were received", reqMsg.Type), errors.UNKNOWN_CODE)
		}
	}
}

func (p *ProxyServer) handleLoginRequest(serveName, loginPass string, coordinator *connector.Coordinator) {
	if context, ok := p.context[serveName]; !ok {
		p.handleConnErr(coordinator, fmt.Errorf("Unknown serve %s were received", serveName), errors.INVALID_SERVE_NAME)
	} else if loginPass == context.LoginPass {
		respMsg := message.ResultMessage{errors.OK}
		content, _ := json.Marshal(respMsg)
		if err := coordinator.SendMsg(string(content)); err != nil {
			p.errs <- err
			coordinator.Close()
		} else {
			controller := newController(coordinator)
			p.controllers[serveName] = controller
			controller.daemonize()
		}
	} else {
		p.handleConnErr(coordinator, fmt.Errorf("Invalid password for login %s != %s", context.LoginPass, loginPass), errors.INVALID_PASSWORD)
	}
}

func (p *ProxyServer) handleAuthRequest(serveName, authPass string, coordinator *connector.Coordinator) {
	if context, ok := p.context[serveName]; !ok {
		p.handleConnErr(coordinator, fmt.Errorf("Unknown serve %s were received", serveName), errors.INVALID_SERVE_NAME)
	} else if controller, ok := p.controllers[serveName]; !ok {
		p.handleConnErr(coordinator, fmt.Errorf("controller for %s does not alive", serveName), errors.INVALID_SERVE_NAME)
	} else {
		if authPass != context.AuthPass {
			p.handleConnErr(coordinator, fmt.Errorf("Invalid password for authenticate %s != %s", context.AuthPass, authPass), errors.INVALID_PASSWORD)
		} else {
			var simpleErrHandle = func(err error) {
				coordinator.Close()
				log.Errorf("Error handling connection in proxy server: %s", err)
				p.errs <- err
			}

			respMsg := message.ResultMessage{errors.OK}
			buf, _ := json.Marshal(respMsg)
			if err := coordinator.SendMsg(string(buf)); err != nil {
				simpleErrHandle(err)
			} else if workConn, err := controller.getConn(); err != nil {
				simpleErrHandle(err)
			} else {
				go drillingTunnel(coordinator.Conn, workConn)
			}
		}
	}
}

func (p *ProxyServer) handleConnErr(coordinator *connector.Coordinator, err error, cErr errors.Error) {
	p.errs <- err
	respMsg := message.ResultMessage{cErr}
	content, _ := json.Marshal(respMsg)
	coordinator.SendMsg(string(content))
	coordinator.Close()
}

func testSocks5Req(src, dst net.Conn) {
	buf := make([]byte, 10)
	dst.Write([]byte{5, 2, 0, 1})
	dst.Read(buf)
	fmt.Println(buf)
	src.Close()
	dst.Close()
}
