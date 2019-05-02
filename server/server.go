package server

import (
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"os"
)

// for opening a listener to proxy request
type ProxyServer struct {
	// generated from client
	listener   *net.Listener
	controller *Controller

	cipher *ss.Cipher
}

func NewProxyServer(cipher *ss.Cipher) *ProxyServer {
	return &ProxyServer{cipher: cipher}
}

func (p *ProxyServer) Listen(southAddress, northAddress, loginPass string) {
	log.Infof("South bridge listen on %s, North bridge listen on %s", southAddress, northAddress)
	go p.serve(southAddress)
	p.openController(northAddress, loginPass)
}

func (p *ProxyServer) serve(address string) {
	if listener, err := net.Listen("tcp", address); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	} else {
		p.listener = &listener
		for {
			rawConn, err := listener.Accept()
			if err != nil {
				log.Fatalln(err)
				continue
			}
			log.Info("get a new request")
			go p.pipe(rawConn)
		}
	}
}

func (p *ProxyServer) openController(address, loginPass string) {
	p.controller = NewController(address, loginPass, p.cipher)
	p.controller.Run()
}

func (p *ProxyServer) pipe(conn net.Conn) {
	if p.controller == nil {
		log.Error("The controller of proxy server is null.")
		return
	}

	workConn, err := p.controller.GetConn()
	if err != nil {
		log.Error(err)
		conn.Close()
		return
	}

	log.Debugf("Pipe between request connection and work connection, %s -> %s", conn.RemoteAddr().String(), workConn.RemoteAddr().String())
	if err != nil {
		log.Println("Listener for incoming connections from client closed")
		log.Error("Error pipe:", err)
	} else {
		//go ss.PipeThenClose(conn, workConn)
		//ss.PipeThenClose(workConn, conn)
		go pipeThenClose(workConn, conn)
		pipeThenClose(conn, workConn)
	}
}

var customLeackyBuf = ss.NewLeakyBuf(2048, 4096)

func pipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	buf := customLeackyBuf.Get()
	defer customLeackyBuf.Put(buf)
	for {
		ss.SetReadTimeout(src)
		n, err := src.Read(buf)
		// log.Infof("%s -> %s", src.LocalAddr().String(), dst.LocalAddr().String())
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			// Note: avoid overwrite err returned by Read.
			if _, err := dst.Write(buf[0:n]); err != nil {
				log.Println("write:", err)
				break
			}
		}
		if err != nil {
			// Always "use of closed network connection", but no easy way to
			// identify this specific error. So just leave the error along for now.
			// More info here: https://code.google.com/p/go/issues/detail?id=4373
			/*
				if bool(Debug) && err != io.EOF {
					Debug.Println("read:", err)
				}
			*/
			break
		}
	}
}
