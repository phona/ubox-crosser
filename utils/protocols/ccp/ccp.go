package ccp

import (
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"ubox-crosser/models/errors"
	"ubox-crosser/utils/protocols"
)

/*
	CCP-Client --------- CCP-Server
		|                    |
           send anth request
		|      -------->     |
              get address
		|      <--------     |
			      OK
		|      -------->     |
       send with another address
      listen   <--------     |
*/

type Class byte

const (
	REQUEST_AUTH     Class = 0x00
	RESPONSE_SUCCESS Class = 0x01
	RESPONSE_ECHO    Class = 0x02
	RESPONSE_ERROR   Class = 0x03
)

type AuthRequest struct {
	Password string `json:"password"`
}

type AuthResponse struct {
	Address string `json:"address"`
}

type EchoMessage struct {
	Result errors.Error `json:"result"`
}

type Message struct {
	Class Class
	Data  []byte
	Err   string
}

type CcpProtocol struct {
	protocols.LineBaseProtocol
	cipher *shadowsocks.Cipher
}

func NewProtocol(conn net.Conn, cipher *shadowsocks.Cipher) *CcpProtocol {
	if cipher == nil {
		cipher, _ = shadowsocks.NewCipher("chacha20", "UBoxtech")
	}
	return &CcpProtocol{
		LineBaseProtocol: protocols.MakeLineBaseProtocol(conn),
		cipher:           cipher,
	}
}

func (p *CcpProtocol) Send(message Message) error {
	buf := make([]byte, 0, 1+len(message.Data))
	buf[0] = byte(message.Class)
	for i, _ := range message.Data {
		buf[i+1] = message.Data[i]
	}
	if err := p.SendMsg(buf); err != nil {
		return err
	}
	return nil
}

func (p *CcpProtocol) Receive() (Message, error) {
	msg := Message{}
	buf, err := p.ReadMsg()
	if err != nil {
		return msg, err
	}

	msg.Class = Class(buf[0])
	if msg.Class != RESPONSE_ERROR {
		msg.Class = Class(buf[0])
		msg.Data = buf[1:]
	} else {
		msg.Data = make([]byte, 0)
		msg.Err = string(buf[1:])
	}
	return msg, nil
}
