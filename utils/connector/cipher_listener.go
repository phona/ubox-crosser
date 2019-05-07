package connector

import (
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
)

type CipherListener struct {
	listener net.Listener
	cipher   *shadowsocks.Cipher
}

func NewCipherListener(listener net.Listener, cipher *shadowsocks.Cipher) *CipherListener {
	return &CipherListener{
		listener: listener,
		cipher:   cipher,
	}
}

func (l *CipherListener) Accept() (net.Conn, error) {
	conn, err := l.listener.Accept()
	if l.cipher != nil {
		conn = shadowsocks.NewConn(conn, l.cipher.Copy())
	}
	return conn, err
}

func (l *CipherListener) Close() error {
	return l.listener.Close()
}

func (l *CipherListener) Addr() net.Addr {
	return l.listener.Addr()
}
