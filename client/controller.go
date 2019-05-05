package client

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/armon/go-socks5"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"io"
	"net"
	"sync"
	"time"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

var (
	HeartBeatInterval int64 = 5
	HeartBeatTimeout  int64 = 30
)

// control channel for communicate with proxy server
type Controller struct {
	Address        string
	ctlConn        *connector.Coordinator
	heartBeatTimer *time.Timer

	Password string

	// this property can be abstracted
	sessionLayer *socks5.Server
	cipher       *shadowsocks.Cipher
	mutex        sync.Mutex
}

func NewController(address string, server *socks5.Server, cipher *shadowsocks.Cipher, password string) *Controller {
	return &Controller{
		Address:        address,
		ctlConn:        nil,
		heartBeatTimer: nil,
		sessionLayer:   server,
		cipher:         cipher,
		Password:       password,
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
			case message.GEN_WORKER:
				// get a generating worker request
				// open a new tcp connection
				go c.newWorkConn()
			case message.HEART_BEAT:
				// a heart beat
				log.Infof("Received a heart beat from %s", c.ctlConn.Conn.RemoteAddr().String())
			default:
				log.Errorf("Unknown type %s were received", respMsg.Type)
			}
		}
	}
}

func (c *Controller) newWorkConn() {
	var reqMsg message.Message
	if workConn, err := c.getConn(); err != nil {
		log.Error("Error generating a worker ", err)
	} else {
		defer workConn.Close()
		reqMsg.Type = message.GEN_WORKER
		buf, _ := json.Marshal(reqMsg)
		// add this connection to server workers pool
		if err := workConn.SendMsg(string(buf)); err != nil {
			log.Infof("Error sending work message to %s in a work connection: %s",
				c.ctlConn.Conn.RemoteAddr().String(), err)
		} else {
			// temp, err := workConn.ReadMsg()
			// log.Info("Work connection received content ", temp, err)
			log.Infof("Create a new socks5 work connection, %s", workConn.Conn.LocalAddr())
			if err := c.sessionLayer.ServeConn(workConn.Conn); err != nil {
				log.Errorf("Error serving a work connection with socks5 protocol: ", err)
			}
		}
	}
}

func (c *Controller) getConn() (*connector.Coordinator, error) {
	if rawConn, err := net.Dial("tcp", c.Address); err != nil {
		return nil, err
	} else {
		if c.cipher != nil {
			rawConn = shadowsocks.NewConn(rawConn, c.cipher.Copy())
		}
		return connector.AsCoordinator(rawConn), nil
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

	reqMsg := message.Message{Type: message.LOGIN, Password: c.Password}
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
	buf, err := json.Marshal(reqMsg)
	if err != nil {
		log.Warn("Serialize reqMsg err! Err: %v", err)
	}

	log.Infof("Start to send heartbeat send %+v", reqMsg)
	for {
		time.Sleep(time.Duration(HeartBeatInterval) * time.Second)
		if c != nil && !coordinator.IsTerminate() {
			err = coordinator.SendMsg(string(buf))
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
