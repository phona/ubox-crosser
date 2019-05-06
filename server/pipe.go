package server

import (
	log "github.com/Sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
)

var customLeackyBuf = ss.NewLeakyBuf(2048, 4096)

func pipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	buf := customLeackyBuf.Get()
	defer customLeackyBuf.Put(buf)
	for {
		ss.SetReadTimeout(src)
		n, err := src.Read(buf)
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			// Note: avoid overwrite err returned by Read.
			log.Infof("%s -> %s size: %d", src.LocalAddr().String(), dst.LocalAddr().String(), n)
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
