package server

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"ubox-crosser/models/config"
	"ubox-crosser/models/errors"
	"ubox-crosser/models/message"
	"ubox-crosser/utils/connector"
	"ubox-crosser/utils/protocols/ccp"
)

// Checking Capability
type CCPServer struct {
	ccpConfig config.CCPConfig
	cipher    *shadowsocks.Cipher
}

func NewCCPServer(ccpConfig config.CCPConfig) *CCPServer {
	if ccpConfig.MaxTry == 0 {
		ccpConfig.MaxTry = DEFAULT_TRY_TIMES
	}

	return &CCPServer{ccpConfig: ccpConfig}
}

func (server *CCPServer) Run() error {
	if l, err := net.Listen("tcp", server.ccpConfig.Address); err != nil {
		return err
	} else {
		if server.ccpConfig.Method != "" {
			if err := shadowsocks.CheckCipherMethod(server.ccpConfig.Method); err != nil {
				return err
			} else if server.cipher, err = shadowsocks.NewCipher(server.ccpConfig.Method, server.ccpConfig.Key); err != nil {
				return err
			}
		}
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Errorf("Error Accepting a connection in udp server: %s", err)
			}
			go server.handleConnection(conn)
		}
	}
}

func (server *CCPServer) handleConnection(conn net.Conn) {
	ccpProtocol := ccp.NewProtocol(conn, server.cipher.Copy())
	if msg, err := ccpProtocol.Receive(); err != nil {
		log.Errorf("(*CCPServer).handleConnection:> %s", err)
	} else if msg.Class == ccp.REQUEST_AUTH {
		var authMsg ccp.AuthRequest
		if err := json.Unmarshal(msg.Data, authMsg); err != nil {
			log.Errorf("(*CCPServer).handleConnection:> %s", err)
		} else if authMsg.Password == server.ccpConfig.LoginPass {
			log.Errorf("(*CCPServer).handleConnection:> Invalid password %s != %s",
				authMsg.Password, server.ccpConfig.LoginPass)
		} else {

		}
	} else {

	}
}

func (server *CCPServer) response(address string) error {
	if reqAddr, err := net.ResolveUDPAddr("udp", address); err != nil {
		return err
	} else if bindAddr, err := net.ResolveUDPAddr("udp", server.ccpConfig.ReqBindAddress); err != nil {
		return err
	} else {
		var conn net.Conn
		// try max times to send response
		for i := 0; uint8(i) < server.ccpConfig.MaxTry; i++ {
			if conn, err = net.DialUDP("udp", bindAddr, reqAddr); err != nil {
				log.Errorf("Error response Retry %d times: %s", i, err)
				continue
			} else {
				if server.cipher != nil {
					conn = shadowsocks.NewConn(conn, server.cipher.Copy())
				}
				coordinator := connector.AsCoordinator(conn)
				buf, _ := json.Marshal(message.ResultMessage{errors.OK})
				if err := coordinator.SendMsg(string(buf)); err != nil {
					log.Errorf("Error response Retry %d times: %s", i, err)
					conn.Close()
					continue
				} else if content, err := coordinator.ReadMsg(); err != nil {
					log.Errorf("Error response Retry %d times: %s", i, err)
					conn.Close()
					continue
				} else {
					var msg message.ResultMessage
					if err := json.Unmarshal([]byte(content), &msg); err != nil {
						log.Errorf("Error response Retry %d times: %s", i, err)
						conn.Close()
						return err
					}
					log.Infof("Send response to %s is success.", address)
				}
			}
		}
		return err
	}
}

func (server *CCPServer) handleConnErr(coordinator *connector.Coordinator, code errors.Error) {
	msg := message.ResultMessage{code}
	buf, _ := json.Marshal(msg)

	if err := coordinator.SendMsg(string(buf)); err != nil {
		log.Errorf("Error sending message in CCP server send request: %s", err)
	}
}
