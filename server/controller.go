package server

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

type Controller struct {
	conn     *connector.Coordinator
	workConn chan net.Conn
}

func NewController(conn *connector.Coordinator) *Controller {
	return &Controller{
		conn:     conn,
		workConn: make(chan net.Conn, 10),
	}
}

func (c *Controller) GetConn() (net.Conn, error) {
	reqMessage := message.Message{Type: message.GEN_WORKER}
	buf, _ := json.Marshal(reqMessage)
	for {
		if err := c.conn.SendMsg(string(buf)); err != nil {
			return nil, err
		} else {
			if rawConn, ok := <-c.workConn; ok {
				return rawConn, nil
			}
		}
	}
}

func (c *Controller) daemonize() {
	// a connected control connection only can receive a heartbeat
	var reqMsg message.Message
	for {
		if c.conn.IsTerminate() {
			break
		} else if content, err := c.conn.ReadMsg(); err != nil {
			log.Error("Error receiving content in daemonize: ", err)
		} else if err := json.Unmarshal([]byte(content), &reqMsg); err != nil {
			log.Errorf("Error Unmarshal content in daemonize %s: %s", content, err)
		} else {
			switch reqMsg.Type {
			case message.HEART_BEAT:
				if buf, err := json.Marshal(reqMsg); err != nil {
					log.Error("Error Sending heartbeat: ", err)
				} else if err := c.conn.SendMsg(string(buf)); err != nil {
					log.Error("Error Sending heartbeat: ", err)
				}
				log.Debug("Received a heartbeat")
			default:
				log.Errorf("Unknown type %s were received", reqMsg.Type)
			}
		}
	}
}

func (c *Controller) HandleConnection(conn net.Conn) {
	coordinator := connector.AsCoordinator(conn)
	var reqMessage message.Message
	if content, err := coordinator.ReadMsg(); err != nil {
		log.Error("Error receiving content: ", err)
		coordinator.Close()
	} else if err := json.Unmarshal([]byte(content), &reqMessage); err != nil {
		log.Errorf("Error Unmarshal content %s: %s", content, err)
	} else {
		log.Debugf("Received content: %s", content)
		switch reqMessage.Type {
		case message.GEN_WORKER:
			c.workConn <- coordinator.Conn
			log.Debug("Add new work connection")
		default:
			log.Errorf("Unknown type %s were received", reqMessage.Type)
		}
	}
}
