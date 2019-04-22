package crosser

import (
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"log"
	"net"
	"os"
)

type Tunnel struct {
	connChannel chan net.Conn
}

func NewTunnel(size int) *Tunnel {
	ch := make(chan net.Conn, size)
	return &Tunnel{connChannel: ch}
}

func (tunnel *Tunnel) OpenNorth(address string) {
	if listener, err := net.Listen("tcp", address); err != nil {
		log.Fatalln(err)
	} else {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
			}

			tunnel.connChannel <- conn
			log.Println("Add new connection to proxy pool")
		}
	}
}

func (tunnel *Tunnel) OpenSouth(address string) {
	if listener, err := net.Listen("tcp", address); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	} else {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				os.Exit(0)
			}
			log.Println("get a request from south, build a tunnel for a request that from south")
			go tunnel.handleConnection(conn)
		}
	}
}

func (tunnel *Tunnel) OpenSouthWithCipher(address, method, password string) {
	var cipher *ss.Cipher
	var err error
	if cipher == nil {
		log.Println("creating cipher for address:", address)
		cipher, err = ss.NewCipher(method, password)
		if err != nil {
			log.Printf("Error generating cipher for address: %s %v\n", address, err)
		}
	}

	if listener, err := net.Listen("tcp", address); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	} else {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				os.Exit(0)
			}
			log.Println("get a request")
			go tunnel.handleConnectionWithCipher(conn, cipher)
		}
	}
}

func (tunnel *Tunnel) handleConnectionWithCipher(conn net.Conn, cipher *ss.Cipher) {
	proxy := <-tunnel.connChannel
	newProxy := ss.NewConn(proxy, cipher.Copy())
	go ss.PipeThenClose(conn, newProxy)
	go ss.PipeThenClose(newProxy, conn)
	log.Println("done")
}

func (tunnel *Tunnel) handleConnection(conn net.Conn) {
	proxy := <-tunnel.connChannel
	go ss.PipeThenClose(conn, proxy)
	go ss.PipeThenClose(proxy, conn)
	log.Println("done")
}
