package server

import (
	log "github.com/sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"sync"
	"time"
)

var customLeackyBuf = ss.NewLeakyBuf(2048, 4096)

func pipeThenClose(src, dst net.Conn, bytesCounter func(uint64)) {
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
				if bytesCounter != nil {
					bytesCounter(uint64(n))
				}
				log.Debugf("pipe: %d, %x", n, buf[0:n])
			}
		}
		if err != nil {
			break
		}
	}
}

func drillingTunnel(src, dst net.Conn, collector *Collector) {
	log.Debugf("Pipe between request connection and work connection, %s -> %s", src.RemoteAddr().String(), dst.RemoteAddr().String())

	if collector != nil {
		collector.RecordTunnelStart()
		start := time.Now()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			pipeThenClose(src, dst, collector.RecordBytesIn)
			wg.Done()
		}()
		pipeThenClose(dst, src, collector.RecordBytesOut)
		wg.Wait()

		collector.RecordTunnelEnd(time.Since(start).Nanoseconds())
	} else {
		go pipeThenClose(src, dst, nil)
		pipeThenClose(dst, src, nil)
	}
}
