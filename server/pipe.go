package server

import (
	log "github.com/sirupsen/logrus"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"net"
	"sync"
	"time"
)

var customLeackyBuf = ss.NewLeakyBuf(2048, 4096)

func pipeThenClose(src, dst net.Conn, stats *Collector, recordIn bool) {
	defer dst.Close()
	buf := customLeackyBuf.Get()
	defer customLeackyBuf.Put(buf)
	for {
		ss.SetReadTimeout(src)
		n, err := src.Read(buf)
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			if stats != nil {
				if recordIn {
					stats.RecordBytesIn(uint64(n))
				} else {
					stats.RecordBytesOut(uint64(n))
				}
			}
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
			break
		}
	}
}

func drillingTunnel(src, dst net.Conn) {
	drillingTunnelWithStats(src, dst, nil)
}

func drillingTunnelWithStats(src, dst net.Conn, stats *Collector) {
	log.Debugf("Pipe between request connection and work connection, %s -> %s", src.RemoteAddr().String(), dst.RemoteAddr().String())
	if stats != nil {
		stats.RecordTunnelStart()
	}
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pipeThenClose(src, dst, stats, true)
		wg.Done()
	}()
	pipeThenClose(dst, src, stats, false)
	wg.Wait()
	if stats != nil {
		stats.RecordTunnelEnd(time.Since(start))
	}
}
