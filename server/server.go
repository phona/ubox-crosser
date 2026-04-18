package server

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"

	"github.com/phona/ubox-crosser/models/config"
	"github.com/phona/ubox-crosser/models/errors"
	"github.com/phona/ubox-crosser/models/message"
	"github.com/phona/ubox-crosser/server/stats"
	"github.com/phona/ubox-crosser/utils/connector"
)

// for opening a listener to proxy request
type ProxyServer struct {
	// generated from client

	dispatcher  *connector.Dispatcher
	controllers map[string]*controller
	errs        chan error

	context   map[string]config.ServerConfig
	collector *stats.Collector
}

func NewProxyServer(configs map[string]config.ServerConfig) *ProxyServer {
	total := len(configs)
	dispatcher := connector.NewDispatcher(uint64(total))
	listenedAddr := make([]string, 0, total)
	collector := stats.NewCollector()
	server := &ProxyServer{
		dispatcher:  dispatcher,
		controllers: make(map[string]*controller, total),
		errs:        make(chan error, 10),
		context:     configs,
		collector:   collector,
	}

	// Start stats HTTP server if configured
	for _, cfg := range configs {
		if cfg.StatsAddress != "" {
			statsServer := stats.NewServer(collector)
			go func(addr string) {
				if err := statsServer.ListenAndServe(addr); err != nil {
					log.Errorf("Stats server error: %s", err)
				}
			}(cfg.StatsAddress)
			break // only one stats server needed
		}
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
		respMsg := message.ResultMessage{Result: message.SUCCESS, Reason: errors.OK}
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

			respMsg := message.ResultMessage{Result: message.SUCCESS, Reason: errors.OK}
			buf, _ := json.Marshal(respMsg)
			if err := coordinator.SendMsg(string(buf)); err != nil {
				simpleErrHandle(err)
			} else if workConn, err := controller.getConn(); err != nil {
				simpleErrHandle(err)
			} else {
				p.collector.TrackConnection()
				startTime := time.Now()
				go func() {
					drillingTunnel(coordinator.Conn, workConn, p.collector)
					p.collector.ReleaseConnection(startTime)
				}()
			}
		}
	}
}

func (p *ProxyServer) handleConnErr(coordinator *connector.Coordinator, err error, cErr errors.Error) {
	p.errs <- err
	respMsg := message.ResultMessage{Result: message.FAILED, Reason: cErr}
	content, _ := json.Marshal(respMsg)
	_ = coordinator.SendMsg(string(content))
	coordinator.Close()
}

