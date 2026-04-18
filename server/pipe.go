package server

import (
	log "github.com/sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"time"
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
			log.Debugf("%s -> %s size: %d, %x", src.LocalAddr().String(), dst.LocalAddr().String(), n, buf[0:n])
			log.Debugf("%s -> %s size: %d, %x", src.RemoteAddr().String(), dst.RemoteAddr().String(), n, buf[0:n])
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
			/*
				if bool(Debug) && err != io.EOF {
					Debug.Println("read:", err)
				}
			*/
			break
		}
	}
}

func drillingTunnel(src, dst net.Conn) {
	log.Debugf("Pipe between request connection and work connection, %s -> %s", src.RemoteAddr().String(), dst.RemoteAddr().String())
	go pipeThenClose(src, dst)
	pipeThenClose(dst, src)
}

func pipeThenCloseWithStats(src, dst net.Conn, collector *Collector, recordBytes func(uint64)) {
	defer dst.Close()
	buf := customLeackyBuf.Get()
	defer customLeackyBuf.Put(buf)
	for {
		ss.SetReadTimeout(src)
		n, err := src.Read(buf)
		if n > 0 {
			log.Debugf("%s -> %s size: %d, %x", src.LocalAddr().String(), dst.LocalAddr().String(), n, buf[0:n])
			log.Debugf("%s -> %s size: %d, %x", src.RemoteAddr().String(), dst.RemoteAddr().String(), n, buf[0:n])
			if _, err := dst.Write(buf[0:n]); err != nil {
				log.Errorln("write:", err)
				break
			} else {
				recordBytes(uint64(n))
				log.Debugf("pipe: %d, %x", n, buf[0:n])
			}
		}
		if err != nil {
			break
		}
	}
}

func drillingTunnelWithStats(src, dst net.Conn, collector *Collector) {
	log.Debugf("Pipe between request connection and work connection, %s -> %s", src.RemoteAddr().String(), dst.RemoteAddr().String())
	collector.RecordTunnelStart()
	start := time.Now()
	go pipeThenCloseWithStats(src, dst, collector, collector.RecordBytesIn)
	pipeThenCloseWithStats(dst, src, collector, collector.RecordBytesOut)
	collector.RecordTunnelEnd(time.Since(start).Nanoseconds())
}
