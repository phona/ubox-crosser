package server

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

type controller struct {
	conn     *connector.Coordinator
	workConn chan net.Conn
}

func newController(conn *connector.Coordinator) *controller {
	return &controller{
		conn:     conn,
		workConn: make(chan net.Conn, DEFAULT_SERVES),
	}
}

func (c *controller) getConn() (net.Conn, error) {
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

func (c *controller) daemonize() {
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
				log.Errorf("Unknown type %d were received", reqMsg.Type)
			}
		}
	}
}
