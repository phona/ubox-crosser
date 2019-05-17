package client

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"io"
	"net"
	"sync"
	"time"
	"ubox-crosser/models/errors"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

var (
	HeartBeatInterval int64 = 10
	HeartBeatTimeout  int64 = 120
)

// control channel for communicate with proxy server
type Controller struct {
	Address        string
	ctlConn        *connector.Coordinator
	heartBeatTimer *time.Timer

	ServeName, Password string

	// this property can be abstracted
	cipher *ss.Cipher
	mutex  sync.Mutex

	messages chan message.Message
}

func NewController(address, serveName, password string, cipher *ss.Cipher) *Controller {
	return &Controller{
		Address:        address,
		ctlConn:        nil,
		heartBeatTimer: nil,
		cipher:         cipher,
		Password:       password,
		ServeName:      serveName,
		messages:       make(chan message.Message, 10),
	}
}

func (c *Controller) Run() {
	var sleepTime time.Duration = 1
	var err error
	for {
		if c.ctlConn == nil || c.ctlConn.IsTerminate() || err == io.EOF {
			// login for construct a control channel
			// always run with the control channel and do heart beat
			log.Info("Generate a control for ", c.Address)
			for {
				err := c.login()
				if err == nil {
					break
				} else {
					log.Errorf("Error in login: %s", err)
					if _, ok := err.(errors.Error); ok {
						return
					}

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
		if c.ctlConn.IsTerminate() {
			break
		}

		var respMsg message.Message
		if content, err := c.ctlConn.ReadMsg(); err != nil {
			log.Error("Error reading respMsg: ", err)
		} else if err := json.Unmarshal([]byte(content), &respMsg); err != nil {
			log.Error("Error unmarshal respMsg: ", err)
			log.Error(content)
		} else {
			// distribute respMsg to different handler
			log.Infof("Received content: %s", content)
			switch respMsg.Type {
			case message.HEART_BEAT:
				// a heart beat
				log.Infof("Received a heart beat from %s", c.ctlConn.Conn.RemoteAddr().String())
			default:
				c.messages <- respMsg
			}
		}
	}
}

func (c *Controller) getConn() (*connector.Coordinator, error) {
	var conn net.Conn
	if addr, err := net.ResolveTCPAddr("", c.Address); err != nil {
		return nil, err
	} else if conn, err = net.DialTCP("tcp", nil, addr); err != nil {
		return nil, err
	} else {
		if c.cipher != nil {
			conn = ss.NewConn(conn, c.cipher.Copy())
		}
		return connector.AsCoordinator(conn), nil
	}
}

func (c *Controller) login() error {
	// for get a control connection
	controlConn, err := c.getConn()
	c.mutex.Lock()
	if c.ctlConn != nil {
		c.ctlConn.Close()
	}
	c.ctlConn = controlConn
	c.mutex.Unlock()
	log.Debug("Get new connection for login to proxy server")
	if err != nil {
		log.Errorf("Error getting new connection for login to proxy server: %s", err)
		return err
	}

	var respMsg message.ResultMessage
	reqMsg := message.Message{Type: message.LOGIN, Password: c.Password, ServeName: c.ServeName}
	buf, _ := json.Marshal(reqMsg)
	if err = controlConn.SendMsg(string(buf)); err != nil {
		return err
	} else if content, err := controlConn.ReadMsg(); err != nil {
		return err
	} else if err := json.Unmarshal([]byte(content), &respMsg); err != nil {
		return err
	} else if respMsg.Result != message.SUCCESS {
		return respMsg.Reason
	}
	// login success
	go c.startHeartBeat(controlConn)
	return nil
}

func (c *Controller) startHeartBeat(coordinator *connector.Coordinator) {
	f := func() {
		log.Error("HeartBeat timeout!")
		if coordinator != nil {
			coordinator.Close()
		}
	}
	c.heartBeatTimer = time.AfterFunc(time.Duration(HeartBeatTimeout)*time.Second, f)
	defer c.heartBeatTimer.Stop()

	reqMsg := message.Message{Type: message.HEART_BEAT}
	buf, _ := json.Marshal(reqMsg)
	log.Infof("Start to send heartbeat send %+v", reqMsg)
	for {
		time.Sleep(time.Duration(HeartBeatInterval) * time.Second)
		if c != nil && !coordinator.IsTerminate() {
			err := coordinator.SendMsg(string(buf))
			log.Info("Send heartbeat to server")
			if err != nil {
				log.Error("Send heartbeat to server failed! Err:%v", err)
				continue
			}
		} else {
			break
		}
	}
	log.Info("Heartbeat exit")
}
