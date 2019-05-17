package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"os"
	"ubox-crosser/models/errors"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

type AuthServer struct {
	targetAddress, serveName, password string
	cipher                             *ss.Cipher
}

func NewAuthServer(targetAddress, serveName, password string, cipher *ss.Cipher) *AuthServer {
	return &AuthServer{
		targetAddress: targetAddress,
		serveName:     serveName,
		password:      password,
		cipher:        cipher,
	}
}

func (a *AuthServer) Listen(address string) {
	if listener, err := net.Listen("tcp", address); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	} else {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				continue
			}
			log.Infof("Get a new request from remote client %s", conn.RemoteAddr().String())
			go a.handleConnection(conn)
		}
	}
}

func (a *AuthServer) handleConnection(src net.Conn) {
	dst, err := a.getConn()
	if err != nil {
		log.Error(err)
		src.Close()
		return
	}
	go drillingTunnel(src, dst)
}

func (a *AuthServer) getConn() (net.Conn, error) {
	var conn net.Conn
	if addr, err := net.ResolveTCPAddr("", a.targetAddress); err != nil {
		return nil, err
	} else if conn, err = net.DialTCP("tcp", nil, addr); err != nil {
		return nil, err
	} else {
		var respMsg message.ResultMessage
		reqMsg := message.Message{
			Type:      message.AUTHENTICATION,
			Password:  a.password,
			ServeName: a.serveName,
		}

		if a.cipher != nil {
			conn = ss.NewConn(conn, a.cipher.Copy())
		}
		coordinator := connector.AsCoordinator(conn)
		buf, _ := json.Marshal(reqMsg)

		var errFunc = func(e error) (net.Conn, error) {
			coordinator.Close()
			return nil, e
		}

		if err := coordinator.SendMsg(string(buf)); err != nil {
			return errFunc(fmt.Errorf("Error sending message: %s", err))
		} else if content, err := coordinator.ReadMsg(); err != nil {
			return errFunc(fmt.Errorf("Error reading message: %s", err))
		} else if err := json.Unmarshal([]byte(content), &respMsg); err != nil {
			return errFunc(fmt.Errorf("Error unmarshal data %s: %s", content, err))
		}

		if respMsg.Result != errors.OK {
			return errFunc(respMsg.Result)
		}
		return conn, nil
	}
}
