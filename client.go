package crosser

import (
	"github.com/armon/go-socks5"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"log"
	"net"
	"time"
)

// a connection it never closed just as a server
type Connector struct {
	conn         net.Conn
	socks5Server *socks5.Server

	address string
}

func NewConnector(server *socks5.Server, address string) *Connector {
	return &Connector{
		socks5Server: server,
		address:      address,
	}
}

func (connector *Connector) RunWithCipher(method, password string) {
	var cipher *ss.Cipher
	var err error
	if cipher == nil {
		log.Println("creating cipher for address:", connector.address)
		cipher, err = ss.NewCipher(method, password)
		if err != nil {
			log.Printf("Error generating cipher for address: %s %v\n", connector.address, err)
			return;
		}
	}
	for {
		conn, err := net.Dial("tcp", connector.address);
		if err != nil {
			log.Printf("Error dialing to %s: %v\n", connector.address, err)
			time.Sleep(time.Second * 100)
		}
		newConn := ss.NewConn(conn, cipher.Copy())
		connector.conn = newConn
		if err := connector.socks5Server.ServeConn(newConn); err != nil {
			log.Printf("Error sending socks5 request for connection: %s %v\n", connector.address, err)
			time.Sleep(time.Second * 100)
		}
	}
}

func (connector *Connector) Run() {
	for {
		conn, err := net.Dial("tcp", connector.address);
		connector.conn = conn
		if err != nil {
			log.Printf("Error dialing to %s: %v\n", connector.address, err)
		} else if err := connector.socks5Server.ServeConn(conn); err != nil {
			log.Printf("Error sending socks5 request for connection: %s %v\n", connector.address, err)
		}
	}
}

func (connector *Connector) Close() {
	connector.conn.Close()
}

