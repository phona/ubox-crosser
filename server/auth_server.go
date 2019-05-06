package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"os"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
)

type AuthServer struct {
	targetAddress string
	password      string
	cipher        *ss.Cipher
}

func NewAuthServer(targetAddress, password string, cipher *ss.Cipher) *AuthServer {
	return &AuthServer{
		targetAddress: targetAddress,
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
			log.Info("get a new request")
			go a.handleConnection(conn)
		}
	}
}

func (a *AuthServer) handleConnection(conn net.Conn) {
	workConn, err := a.getConn()
	if err != nil {
		log.Error(err)
		conn.Close()
		return
	}
	go a.pipe(conn, workConn)
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
			Type:     message.LOGIN,
			Password: a.password,
		}

		if a.cipher != nil {
			conn = ss.NewConn(conn, a.cipher.Copy())
		}
		coordinator := connector.AsCoordinator(conn)
		buf, _ := json.Marshal(reqMsg)

		var errFunc = func(e error) (net.Conn, error) {
			conn.Close()
			return nil, e
		}

		if err := coordinator.SendMsg(string(buf)); err != nil {
			return errFunc(fmt.Errorf("Error sending message: %s", err))
		} else if content, err := coordinator.ReadMsg(); err != nil {
			return errFunc(fmt.Errorf("Error reading message: %s", err))
		} else if err := json.Unmarshal([]byte(content), &respMsg); err != nil {
			return errFunc(fmt.Errorf("Error unmarshal data %s: %s", content, err))
		}

		switch respMsg.Result {
		case message.SUCCESS:
			return conn, nil
		case message.FAILED:
			return errFunc(fmt.Errorf("Login failure"))
		default:
			return errFunc(fmt.Errorf("Invalid status code %s", respMsg.Result))
		}
	}
}

func (a *AuthServer) pipe(conn, workConn net.Conn) {
	go pipeThenClose(workConn, conn)
	pipeThenClose(conn, workConn)
}
