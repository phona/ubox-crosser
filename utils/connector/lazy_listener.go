package connector

import (
	"github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
)

type LazyListener struct {
	listener         net.Listener
	cipher           *shadowsocks.Cipher
	network, address string
}

func NewLazyListener(network, address string, cipher *shadowsocks.Cipher) *LazyListener {
	return &LazyListener{network: network, address: address, listener: nil, cipher: cipher}
}

func (l *LazyListener) doListen() (err error) {
	if l.listener == nil {
		if l.listener, err = net.Listen(l.network, l.address); err != nil {
			return
		}
	}
	return
}

func (l *LazyListener) Accept() (net.Conn, error) {
	if err := l.doListen(); err != nil {
		return nil, err
	} else {
		conn, err := l.listener.Accept()
		if l.cipher != nil {
			conn = shadowsocks.NewConn(conn, l.cipher.Copy())
		}
		return conn, err
	}
}

func (l *LazyListener) Close() error {
	if err := l.doListen(); err != nil {
		return err
	} else {
		return l.listener.Close()
	}
}

func (l *LazyListener) Addr() net.Addr {
	if err := l.doListen(); err != nil {
		return simpleAddr{network: l.network, address: l.address}
	} else {
		return l.listener.Addr()
	}
}

type simpleAddr struct {
	network, address string
}

func (s simpleAddr) Network() string {
	return s.network
}

func (s simpleAddr) String() string {
	return s.address
}
