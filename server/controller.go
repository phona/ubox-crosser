package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net"
	"os"
	"sync"
	"ubox-crosser/models"
	"ubox-crosser/utils/conn"
)

type Controller struct {
	Address string
	ctlConn *conn.Coordinator

	workConn chan net.Conn
	listened bool
	mutex    sync.Mutex
}

func NewController(address string) *Controller {
	return &Controller{
		Address:  address,
		ctlConn:  nil,
		listened: false,
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
		coordinator := conn.AsCoordinator(rawConn, 0)
		go c.handleConnection(coordinator)
	}
}

func (c *Controller) GetConn() (net.Conn, error) {
	if c.ctlConn == nil {
		return nil, fmt.Errorf("The controller is not running")
	}

	reqMessage := models.Message{Type: models.GEN_WORKER}
	buf, _ := json.Marshal(reqMessage)
	for {
		if err := c.ctlConn.SendMsg(string(buf)); err != nil {
			return nil, err
		} else {
			if c, ok := <-c.workConn; ok {
				return c, nil
			}
		}
	}
}

func (c *Controller) handleConnection(coordinator *conn.Coordinator) {
	var reqMessage models.Message
	var respMessage models.Message
	if content, err := coordinator.ReadMsg(); err != nil {
		coordinator.Close()
		log.Error("Error receiving content: ", err)
	} else if err := json.Unmarshal([]byte(content), &reqMessage); err != nil {
		log.Errorf("Error Unmarshal content %s: %s", content, err)
	} else {
		log.Infof("Received content: %s", content)
		switch reqMessage.Type {
		case models.LOGIN:
			c.mutex.Lock()
			if c.ctlConn != nil {
				log.Info("Control connection was replaced by new connection")
				c.ctlConn.Close()
			}
			// TODO: dangerous
			c.ctlConn = coordinator
			log.Debug("Setting new control connection")
			go func() {
				for {
					if !c.ctlConn.IsTerminate() {
						log.Infof("%p", c.ctlConn)
						c.handleConnection(c.ctlConn)
					} else {
						break
					}
				}
			}()
			c.mutex.Unlock()
		case models.HEART_BEAT:
			respMessage.Type = models.HEART_BEAT
			buf, _ := json.Marshal(&respMessage)
			if err := c.ctlConn.SendMsg(string(buf)); err != nil {
				log.Error("Error Sending heartbeat: ", err)
			}
			log.Debug("Received a heartbeat")
		case models.GEN_WORKER:
			c.workConn <- coordinator.Conn
			log.Infof("Add new work connection, now size is %d", len(c.workConn))
		default:
			log.Errorf("Unknown type %s, message % were received", reqMessage.Type, reqMessage.Msg)
		}
	}
}
