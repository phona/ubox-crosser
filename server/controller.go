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
	"ubox-crosser/utils/conn"
)

type Controller struct {
	Address string
	ctlConn *conn.Coordinator
	cipher  *shadowsocks.Cipher

	workConn chan net.Conn
	listened bool
	mutex    sync.Mutex
}

func NewController(address string, cipher *shadowsocks.Cipher) *Controller {
	return &Controller{
		Address:  address,
		ctlConn:  nil,
		listened: false,
		cipher:   cipher,
		workConn: make(chan net.Conn, 10),
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
		coordinator := conn.AsCoordinator(rawConn)
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

func (c *Controller) handleConnection(coordinator *conn.Coordinator) {
	var reqMessage message.Message
	var respMessage message.Message
	if content, err := coordinator.ReadMsg(); err != nil {
		log.Error("Error receiving content: ", err)
		coordinator.Close()
	} else if err := json.Unmarshal([]byte(content), &reqMessage); err != nil {
		log.Errorf("Error Unmarshal content %s: %s", content, err)
	} else {
		log.Infof("Received content: %s", content)
		switch reqMessage.Type {
		case message.LOGIN:
			c.mutex.Lock()
			if c.ctlConn != nil {
				log.Info("Control connection was replaced by new connection")
				c.ctlConn.Close()
			}
			// TODO: dangerous
			c.ctlConn = coordinator
			c.mutex.Unlock()
			log.Debug("Setting new control connection")
			go func() {
				for {
					if !c.ctlConn.IsTerminate() {
						c.handleConnection(c.ctlConn)
					} else {
						break
					}
				}
			}()
		case message.HEART_BEAT:
			respMessage.Type = message.HEART_BEAT
			buf, _ := json.Marshal(&respMessage)
			if err := c.ctlConn.SendMsg(string(buf)); err != nil {
				log.Error("Error Sending heartbeat: ", err)
			}
			log.Debug("Received a heartbeat")
		case message.GEN_WORKER:
			c.workConn <- coordinator.Conn
			log.Debug("Add new work connection")
		default:
			log.Errorf("Unknown type %s, message % were received", reqMessage.Type, reqMessage.Msg)
		}
	}
}
