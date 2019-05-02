package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"os"
	"sync"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

type Controller struct {
	Address string
	cipher  *shadowsocks.Cipher

	ctlConn   *connector.Coordinator
	workConn  chan net.Conn
	mutex     sync.Mutex
	LoginPass string
}

func NewController(address, loginPass string, cipher *shadowsocks.Cipher) *Controller {
	return &Controller{
		Address:   address,
		cipher:    cipher,
		LoginPass: loginPass,
		workConn:  make(chan net.Conn, 10),
	}
}

func (c *Controller) Run() {
	// open a connection as controller
	// 1. received heartbeat
	// 2. sending request to get a worker connection
	l, err := net.Listen("tcp", c.Address)
	if err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}
	for {
		rawConn, err := l.Accept()
		if err != nil {
			log.Error("Error accepting new connection: ", err)
		}
		if c.cipher != nil {
			rawConn = shadowsocks.NewConn(rawConn, c.cipher.Copy())
		}
		coordinator := connector.AsCoordinator(rawConn)
		go c.handleConnection(coordinator)
	}
}

func (c *Controller) GetConn() (net.Conn, error) {
	if c.ctlConn == nil {
		return nil, fmt.Errorf("The controller is not running")
	}

	reqMessage := message.Message{Type: message.GEN_WORKER}
	buf, _ := json.Marshal(reqMessage)
	for {
		if err := c.ctlConn.SendMsg(string(buf)); err != nil {
			return nil, err
		} else {
			if rawConn, ok := <-c.workConn; ok {
				return rawConn, nil
			}
		}
	}
}

func (c *Controller) login(reqMsg message.Message, coordinator *connector.Coordinator) {
	if reqMsg.Password != c.LoginPass {
		log.Errorf("Invalid login password for user %s: %+v", reqMsg.Name, reqMsg)
		return
	}
	if reqMsg.Name == "" {
		log.Errorf("Username can't be empty when login to server: %+v", reqMsg)
		return
	}
	c.mutex.Lock()
	if c.ctlConn != nil {
		c.ctlConn.Close()
	}
	c.ctlConn = coordinator
	c.mutex.Unlock()
	go c.daemonize(coordinator)
}

func (c *Controller) daemonize(coordinator *connector.Coordinator) {
	// a connected control connection only can receive a heartbeat
	var reqMsg message.Message
	for {
		if coordinator.IsTerminate() {
			break
		} else if content, err := coordinator.ReadMsg(); err != nil {
			log.Error("Error receiving content in daemonize: ", err)
		} else if err := json.Unmarshal([]byte(content), &reqMsg); err != nil {
			log.Errorf("Error Unmarshal content in daemonize %s: %s", content, err)
		} else {
			switch reqMsg.Type {
			case message.HEART_BEAT:
				if buf, err := json.Marshal(reqMsg); err != nil {
					log.Error("Error Sending heartbeat: ", err)
				} else if err := coordinator.SendMsg(string(buf)); err != nil {
					log.Error("Error Sending heartbeat: ", err)
				}
				log.Debug("Received a heartbeat")
			default:
				log.Errorf("Unknown type %s, message % were received", reqMsg.Type, reqMsg.Msg)
			}
		}
	}
}

func (c *Controller) handleConnection(coordinator *connector.Coordinator) {
	var reqMessage message.Message
	if content, err := coordinator.ReadMsg(); err != nil {
		log.Error("Error receiving content: ", err)
		coordinator.Close()
	} else if err := json.Unmarshal([]byte(content), &reqMessage); err != nil {
		log.Errorf("Error Unmarshal content %s: %s", content, err)
	} else {
		log.Debugf("Received content: %s", content)
		switch reqMessage.Type {
		case message.LOGIN:
			go c.login(reqMessage, coordinator)
		case message.GEN_WORKER:
			c.workConn <- coordinator.Conn
			log.Debug("Add new work connection")
		default:
			log.Errorf("Unknown type %s, message % were received", reqMessage.Type, reqMessage.Msg)
		}
	}
}
