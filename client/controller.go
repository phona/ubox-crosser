package client

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/armon/go-socks5"
	"io"
	"net"
	"sync"
	"time"
	"ubox-crosser/models"
	"ubox-crosser/utils/conn"
)

// control channel for communicate with proxy server
type Controller struct {
	Address        string
	coordinator    *conn.Coordinator
	heartBeatTimer *time.Timer
	count          uint
	sessionLayer   *socks5.Server

	mutex sync.Mutex
}

func NewController(address string, server *socks5.Server) *Controller {
	return &Controller{
		Address:        address,
		coordinator:    nil,
		heartBeatTimer: nil,
		sessionLayer:   server,
	}
}

func (c *Controller) Run() {
	var sleepTime time.Duration = 1
	var err error
	for {
		if c.coordinator == nil || c.coordinator.IsTerminate() || err == io.EOF {
			// login for construct a control channel
			// always run with the control channel and do heart beat
			log.Info("Generate a control for ", c.Address)
			for {
				err := c.login()
				if err == nil {
					break
				} else {
					log.Error("Error in login: %s", err)
					if sleepTime < 60 {
						sleepTime = sleepTime * 2
					}
					time.Sleep(sleepTime * time.Second)
				}
			}
		} else {
			c.handleMessage()
		}
	}
}

func (c *Controller) handleMessage() {
	for {
		if c.coordinator.IsTerminate() {
			break
		}

		var message models.Message
		if content, err := c.coordinator.ReadMsg(); err != nil {
			log.Error("Error reading message: ", err)
		} else if err := json.Unmarshal([]byte(content), &message); err != nil {
			log.Error("Error unmarshal message: ", err)
		} else {
			// distribute message to different handler
			log.Infof("Received content: %s", content)
			switch message.Type {
			case models.GEN_WORKER:
				// get a generating worker request
				// open a new tcp connection
				go c.newWorkConn()
			case models.HEART_BEAT:
				// a heart beat
				log.Infof("Received a heart beat from %s", c.coordinator.Conn.RemoteAddr().String())
			default:
				log.Errorf("Unknown type %s, message % were received", message.Type, message.Msg)
			}
		}
	}
}

func (c *Controller) newWorkConn() {
	var message models.Message
	if workConn, err := c.getConn(); err != nil {
		log.Error("Error generating a worker ", err)
	} else {
		message.Type = models.GEN_WORKER
		message.Msg = ""
		buf, _ := json.Marshal(message)
		if err := workConn.SendMsg(string(buf)); err != nil {
			log.Infof("Error sending work message to %s in a work connection: %s",
				c.coordinator.Conn.RemoteAddr().String(), err)
		} else {
			// temp, err := workConn.ReadMsg()
			// log.Info("Work connection received content ", temp, err)
			log.Info("Create a new socks5 work connection")
			if err := c.sessionLayer.ServeConn(workConn.Conn); err != nil {
				log.Errorf("Error serving a work connection with socks5 protocol: ", err)
			}
		}
	}
}

func (c *Controller) getConn() (*conn.Coordinator, error) {
	if rawConn, err := net.Dial("tcp", c.Address); err != nil {
		return nil, err
	} else {
		return conn.AsCoordinator(rawConn, 0), nil
	}
}

func (c *Controller) login() error {
	// for get a control connection
	controlConn, err := c.getConn()
	if c.coordinator != nil {
		c.coordinator.Close()
	}
	c.coordinator = controlConn
	log.Debug("Get new connection for login to proxy server")
	if err != nil {
		log.Errorf("Error getting new connection for login to proxy server: %s", err)
		return err
	}

	reqMsg := models.Message{Type: models.LOGIN, Name: fmt.Sprintf("%p#%d", c, c.count)}
	c.count++
	buf, _ := json.Marshal(reqMsg)
	err = controlConn.SendMsg(string(buf))
	if err != nil {
		log.Errorf("Control connection construct failed: %s", err)
		return err
	}
	// login success
	go c.startHeartBeat(controlConn)
	return nil
}

func (c *Controller) startHeartBeat(conn *conn.Coordinator) {
	f := func() {
		log.Error("HeartBeat timeout!")
		if conn != nil {
			conn.Close()
		}
	}
	c.heartBeatTimer = time.AfterFunc(time.Duration(HeartBeatTimeout)*time.Second, f)
	defer c.heartBeatTimer.Stop()

	reqMsg := models.Message{Type: models.HEART_BEAT}
	buf, err := json.Marshal(reqMsg)
	if err != nil {
		log.Warn("Serialize clientCtlReq err! Err: %v", err)
	}

	log.Infof("Start to send heartbeat send %+v", reqMsg)
	for {
		time.Sleep(time.Duration(HeartBeatInterval) * time.Second)
		if c != nil && !conn.IsTerminate() {
			err = conn.SendMsg(string(buf))
			log.Info("Send hearbeat to server")
			if err != nil {
				log.Error("Send hearbeat to server failed! Err:%v", err)
				continue
			}
		} else {
			break
		}
	}
	log.Debug("Heartbeat exit")
}
