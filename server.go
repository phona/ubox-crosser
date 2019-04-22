package crosser

import (
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"io"
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
		os.Exit(0)
	} else {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				os.Exit(0)
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
		var proxy net.Conn
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				os.Exit(0)
			}
			log.Println("get a request from south, build a tunnel for a request that from south")
			for {
				proxy = <-tunnel.connChannel
				if r, err := IsUsableConnection(proxy); err != nil || !r {
					if err != nil {
						log.Fatal("Get unusable connection from north connection pool")
					}
					proxy.Close()
					proxy = <-tunnel.connChannel
				} else {
					break
				}
			}
			go ss.PipeThenClose(conn, proxy)
			go ss.PipeThenClose(proxy, conn)
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
		var proxy net.Conn
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				os.Exit(0)
			}
			log.Println("get a request")
			for {
				proxy = <-tunnel.connChannel
				if r, err := IsUsableConnection(proxy); err != nil || !r {
					if err != nil {
						log.Fatal("Get unusable connection from north connection pool")
					}
					proxy.Close()
					proxy = <-tunnel.connChannel
				} else {
					break
				}
			}
			newProxy := ss.NewConn(proxy, cipher.Copy())
			go ss.PipeThenClose(conn, newProxy)
			go ss.PipeThenClose(newProxy, conn)
		}
	}
}

func IsUsableConnection(c net.Conn) (bool, error) {
	ss.SetReadTimeout(c)
	if _, err := c.Read([]byte{}); err == io.EOF {
		return false, err
	} else {
		return true, err
	}
}
