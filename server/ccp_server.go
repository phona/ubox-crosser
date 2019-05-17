package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"time"
	"ubox-crosser/models/config"
	"ubox-crosser/utils"
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
	if ccpConfig.Timeout == 0 {
		ccpConfig.Timeout = 30
	}

	return &CCPServer{ccpConfig: ccpConfig}
}

func (server *CCPServer) Run() error {
	content, _ := json.Marshal(server.ccpConfig)
	log.Infof("CCP server running with %s", content)
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
	ccpProtocol := ccp.NewProtocol(conn, server.cipher)
	var handleErr = func(errMessage string) {
		err := ccp.SendErrMsg(ccpProtocol, errMessage)
		if err != nil {
			log.Errorf("(*CCPServer).handleConnection:> SendErrMsg failed %s", err)
		}
	}

	if msg, err := ccpProtocol.Receive(); err != nil {
		log.Errorf("(*CCPServer).handleConnection:> Error in receive: %s", err)
	} else {
		switch msg.Class {
		case ccp.REQUEST_AUTH:
			var authMsg ccp.AuthRequest
			if err := json.Unmarshal(msg.Data, &authMsg); err != nil {
				log.Errorf("(*CCPServer).handleConnection:> %s", err)
				handleErr("Invalid data format")
			} else if authMsg.Password != server.ccpConfig.LoginPass {
				log.Errorf("(*CCPServer).handleConnection:> Invalid password %s != %s",
					authMsg.Password, server.ccpConfig.LoginPass)
				handleErr("Invalid password")
			} else {
				if udpAddr, err := net.ResolveUDPAddr("udp", conn.RemoteAddr().String()); err != nil {
					handleErr(fmt.Sprintf("Invalid address %s", conn.RemoteAddr().String()))
				} else if udpConn, err := net.DialUDP("udp", nil, udpAddr); err != nil {
					handleErr("Dial udp failure")
				} else {
					defer udpConn.Close()
					err := ccpProtocol.Send(ccp.Message{
						ccp.RESPONSE_SUCCESS,
						utils.MakeJsonBuf(ccp.AuthResponse{udpConn.LocalAddr().String()}),
						"",
					})
					if err != nil {
						log.Error("(*CCPServer).handleConnection:> ccp protocol send failed")
					} else if msg, err = ccpProtocol.Receive(); err != nil {
						handleErr("receive message failed")
					} else if msg.Class != ccp.RESPONSE_SUCCESS {
						log.Errorf("(*CCPServer).handleConnection:> ccp protocol received %#X", ccp.RESPONSE_SUCCESS)
					} else if err := conn.Close(); err != nil {
						log.Error("(*CCPServer).handleConnection:> close connection failied")
					} else {
						var i uint8 = 0
						for ; i < server.ccpConfig.MaxTry; i++ {
							err = udpConn.SetReadDeadline(time.Now().Add(time.Duration(server.ccpConfig.Timeout * time.Second)))
							log.Debugf("(*CCPServer).handleConnection:> send udp for %s", udpConn.RemoteAddr().String())
							if _, err := udpConn.Write([]byte{0}); err != nil {
								log.Errorf("(*CCPServer).handleConnection:> send udp failed: %s, already try %d times", err, i+1)
							} else if _, addr, err := udpConn.ReadFrom([]byte{0}); err != nil {
								log.Errorf("(*CCPServer).handleConnection:> received udp failed: %s, already try %d times", err, i+1)
							} else {
								log.Infof("(*CCPServer).handleConnection:> ckeck capability is successful for %s", addr.String())
								return
							}
						}
					}
				}
			}
		default:
			log.Errorf("(*CCPServer).handleConnection:> Unknown class %#X", msg.Class)
			handleErr("Invalid class")
		}
	}
}
