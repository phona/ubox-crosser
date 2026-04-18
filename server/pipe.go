package server

import (
	log "github.com/sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"time"
)

var customLeackyBuf = ss.NewLeakyBuf(2048, 4096)

func pipeThenClose(src, dst net.Conn, byteCounter func(int64)) {
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
			log.Debugf("%s -> %s size: %d, %x", src.LocalAddr().String(), dst.LocalAddr().String(), n, buf[0:n])
			log.Debugf("%s -> %s size: %d, %x", src.RemoteAddr().String(), dst.RemoteAddr().String(), n, buf[0:n])
			if _, err := dst.Write(buf[0:n]); err != nil {
				log.Errorln("write:", err)
				break
			} else {
				if byteCounter != nil {
					byteCounter(int64(n))
				}
				log.Debugf("pipe: %d, %x", n, buf[0:n])
			}
		}
		if err != nil {
			break
		}
	}
}

func drillingTunnel(src, dst net.Conn) {
	log.Debugf("Pipe between request connection and work connection, %s -> %s", src.RemoteAddr().String(), dst.RemoteAddr().String())
	go pipeThenClose(src, dst, nil)
	pipeThenClose(dst, src, nil)
}

func drillingTunnelWithStats(src, dst net.Conn, stats *Collector) {
	log.Debugf("Pipe between request connection and work connection, %s -> %s", src.RemoteAddr().String(), dst.RemoteAddr().String())
	stats.OnTunnelStart()
	start := time.Now()
	go pipeThenClose(src, dst, stats.AddBytesIn)
	pipeThenClose(dst, src, stats.AddBytesOut)
	stats.OnTunnelEnd(start)
}
