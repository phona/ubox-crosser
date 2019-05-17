package connector

import (
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
)

func WrapConn(conn net.Conn, cipher *shadowsocks.Cipher) net.Conn {
	if cipher != nil {
		conn = shadowsocks.NewConn(conn, cipher.Copy())
	}
	return conn
}
