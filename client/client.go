package client

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/armon/go-socks5"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"ubox-crosser/client/mode"
	messageLib "ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

type Client struct {
	controller   *Controller
	mode         mode.ConnectMode
	cipher       *shadowsocks.Cipher
	sessionLayer *socks5.Server
}

func NewClient(mode mode.ConnectMode, cipher *shadowsocks.Cipher) *Client {
	return &Client{cipher: cipher, mode: mode}
}

func (cli *Client) Connect(address, name, password string) error {
	conf := &socks5.Config{}
	if server, err := socks5.New(conf); err != nil {
		return err
	} else {
		cli.sessionLayer = server
		cli.controller = NewController(address, name, password, cli.cipher)
		cli.controller.Run()
		return nil
	}
}

func (cli *Client) handleMessage() {
	for {
		message := <-cli.controller.messages
		switch message.Type {
		case messageLib.GEN_WORKER:
			// TODO: apply a work connection
		default:
			log.Errorf("Unknown type %s were received", message.Type)
		}
	}
}

func (cli *Client) NewWorkConn(address string) {
	if workConn, err := cli.mode.GetConn(address); err != nil {
		log.Error("Error generating a worker ", err)
	} else {
		defer workConn.Close()
		reqMsg := messageLib.Message{
			Type:      messageLib.GEN_WORKER,
			Password:  cli.controller.Password,
			ServeName: cli.controller.ServeName,
		}
		if cli.cipher != nil {
			workConn = shadowsocks.NewConn(workConn, cli.cipher.Copy())
		}
		coordinator := connector.AsCoordinator(workConn)
		buf, _ := json.Marshal(reqMsg)
		// add this connection to server workers pool
		if err := coordinator.SendMsg(string(buf)); err != nil {
			log.Infof("Error sending work message to %s in a work connection: %s", address, err)
		} else {
			log.Infof("Create a new socks5 work connection, %s", workConn.LocalAddr().String())
			if err := cli.sessionLayer.ServeConn(workConn); err != nil {
				log.Errorf("Error serving a work connection with socks5 protocol: ", err)
			}
		}
	}
}
