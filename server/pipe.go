package server

import (
	log "github.com/sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"

	"github.com/phona/ubox-crosser/server/stats"
)

var customLeackyBuf = ss.NewLeakyBuf(2048, 4096)

func pipeThenClose(src, dst net.Conn, bytesCounter func(int64)) {
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
			if bytesCounter != nil {
				bytesCounter(int64(n))
			}
			if _, err := dst.Write(buf[0:n]); err != nil {
				log.Errorln("write:", err)
				break
			} else {
				log.Debugf("pipe: %d, %x", n, buf[0:n])
			}
		}
		if err != nil {
			// Always "use of closed network connection", but no easy way to
			// identify this specific error. So just leave the error along for now.
			// More info here: https://code.google.com/p/go/issues/detail?id=4373
			break
		}
	}
}

func drillingTunnel(src, dst net.Conn, collector *stats.Collector) {
	log.Debugf("Pipe between request connection and work connection, %s -> %s", src.RemoteAddr().String(), dst.RemoteAddr().String())
	var inCounter, outCounter func(int64)
	if collector != nil {
		inCounter = collector.AddBytesIn
		outCounter = collector.AddBytesOut
	}
	go pipeThenClose(src, dst, inCounter)
	pipeThenClose(dst, src, outCounter)
}
