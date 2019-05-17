package mode

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"ubox-crosser/utils"
	"ubox-crosser/utils/protocols/ccp"
)

const (
	BrigdeMode                   = "brigde-mode"
	P2pFullConeNatMode           = "P2pFullConeNAT-mode"
	P2pAddrRestrictedConeNatMode = "P2pAddrRestrictedConeNAT-Mode"
)

// This interface describe how to get a work connection for a client
type ConnectMode interface {
	Mode() string
	GetConn(addr string) (net.Conn, error)
}

type brigdeMode struct {
}

func newBrigdeMode() *brigdeMode {
	return &brigdeMode{}
}

func (*brigdeMode) Mode() string {
	return BrigdeMode
}

func (*brigdeMode) GetConn(addr string) (net.Conn, error) {
	var conn net.Conn
	if addr, err := net.ResolveTCPAddr("", addr); err != nil {
		return nil, nil
	} else if conn, err = net.DialTCP("tcp", nil, addr); err != nil {
		return nil, nil
	} else {
		return conn, nil
	}
}

type p2pFullConeNatMode struct {
}

func newP2pFullConeNatMode(address, password string, cipher *shadowsocks.Cipher) (*p2pFullConeNatMode, error) {
	if err := handShake(address, password, cipher, nil); err != nil {
		return nil, err
	} else {
		return &p2pFullConeNatMode{}, nil
	}
}

func (*p2pFullConeNatMode) Mode() string {
	return P2pFullConeNatMode
}

func (*p2pFullConeNatMode) GetConn(addr string) (net.Conn, error) {
	if l, err := net.Listen("tcp", addr); err != nil {
		return nil, err
	} else if conn, err := l.Accept(); err != nil {
		return nil, err
	} else if err := l.Close(); err != nil {
		return nil, err
	} else {
		return conn, nil
	}
}

type p2pAddrRestrictedConeNatMode struct {
}

func newP2pAddrRestrictedConeNatMode(address, password string, cipher *shadowsocks.Cipher) (mode *p2pAddrRestrictedConeNatMode, err error) {
	var sub = func(currentLocalAddr string, response ccp.AuthResponse) error {
		if localAddr, err := net.ResolveUDPAddr("udp", currentLocalAddr); err != nil {
			return err
		} else if remoteAddr, err := net.ResolveUDPAddr("udp", response.Address); err != nil {
			return err
		} else if udpConn, err := net.DialUDP("udp", localAddr, remoteAddr); err != nil {
			return err
		} else if _, err := udpConn.Write([]byte{0}); err != nil {
			return err
		} else if err := udpConn.Close(); err != nil {
			return err
		} else {
			return nil
		}
	}
	if err := handShake(address, password, cipher, sub); err != nil {
		return nil, err
	} else {
		return &p2pAddrRestrictedConeNatMode{}, nil
	}
}

func (*p2pAddrRestrictedConeNatMode) Mode() string {
	return P2pAddrRestrictedConeNatMode
}

func (*p2pAddrRestrictedConeNatMode) GetConn(addr string) (net.Conn, error) {
	if conn, err := net.Dial("udp", addr); err != nil {
		return nil, err
	} else {
		conn.Close()
	}

	if l, err := net.Listen("tcp", addr); err != nil {
		return nil, err
	} else if conn, err := l.Accept(); err != nil {
		return nil, err
	} else if err := l.Close(); err != nil {
		return nil, err
	} else {
		return conn, nil
	}
}

/**
 * remote_add is remote upd server address
 */
func GetConnectMode(address, password string, cipher *shadowsocks.Cipher) (ConnectMode, error) {
	if mode2, err := newP2pFullConeNatMode(address, password, cipher); err != nil {
		log.Errorf("Full-cone NAT doesn't supported: %s", err)
	} else {
		return mode2, nil
	}

	if mode1, err := newP2pAddrRestrictedConeNatMode(address, password, cipher); err != nil {
		log.Errorf("(Address)-restricted-cone NAT doesn't supported: %s", err)
	} else {
		return mode1, nil
	}
	return newBrigdeMode(), nil
}

// Received one remote ccp server address
// Target: do authentication, get remote udp conn address, try to get request from this remote address
func handShake(address, password string,
	cipher *shadowsocks.Cipher,
	subroutine func(currentLocalAddr string, msg ccp.AuthResponse) error) (err error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}

	ccpProtocol := ccp.NewProtocol(conn, cipher)
	if err = ccpProtocol.Send(ccp.Message{
		ccp.REQUEST_AUTH,
		utils.MakeJsonBuf(ccp.AuthRequest{password}),
		"",
	}); err != nil {
		return err
	} else if respMsg, err := ccpProtocol.Receive(); err != nil {
		return err
	} else {
		if respMsg.Err != "" {
			err = fmt.Errorf(respMsg.Err)
			return err
		}

		var respAuthMsg ccp.AuthResponse
		if err = json.Unmarshal(respMsg.Data, &respAuthMsg); err != nil {
			return err
		} else {
			if subroutine != nil {
				if err := subroutine(conn.LocalAddr().String(), respAuthMsg); err != nil {
					return err
				}
			}
			if err := ccpProtocol.Send(ccp.Message{
				ccp.RESPONSE_SUCCESS,
				[]byte{},
				"",
			}); err != nil {
				return err
			} else if err := conn.Close(); err != nil {
				return err
			}

			pc, err := net.ListenPacket("udp", conn.LocalAddr().String())
			log.Infof("handShake:> Listen udp on %s", conn.LocalAddr().String())
			if err != nil {
				return err
			} else if _, addr, err := pc.ReadFrom([]byte{0}); err != nil {
				return err
			} else {
				defer pc.Close()
				log.Infof("Received packet from %s", addr.String())
				if _, err := pc.WriteTo([]byte{0}, addr); err != nil {
					return err
				}
				return nil
			}
		}
	}
}
