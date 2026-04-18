package server

import (
	"sync/atomic"
	"time"
)

type StatsResponse struct {
	TotalConnections     uint64  `json:"total_connections"`
	ActiveConnections    uint64  `json:"active_connections"`
	TotalBytesIn         uint64  `json:"total_bytes_in"`
	TotalBytesOut        uint64  `json:"total_bytes_out"`
	AvgTunnelDurationMs  float64 `json:"avg_tunnel_duration_ms"`
}

type Collector struct {
	totalConnections     uint64
	activeConnections    uint64
	totalBytesIn         uint64
	totalBytesOut        uint64
	totalTunnelDurationNs uint64
	completedTunnels     uint64
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) RecordConnection() {
	atomic.AddUint64(&c.totalConnections, 1)
}

func (c *Collector) RecordTunnelStart() {
	atomic.AddUint64(&c.activeConnections, 1)
}

func (c *Collector) RecordTunnelEnd(duration time.Duration) {
	atomic.AddUint64(&c.totalTunnelDurationNs, uint64(duration.Nanoseconds()))
	atomic.AddUint64(&c.completedTunnels, 1)
	atomic.AddUint64(&c.activeConnections, ^uint64(0)) // decrement
}

func (c *Collector) RecordBytesIn(n uint64) {
	atomic.AddUint64(&c.totalBytesIn, n)
}

func (c *Collector) RecordBytesOut(n uint64) {
	atomic.AddUint64(&c.totalBytesOut, n)
}

func (c *Collector) Snapshot() StatsResponse {
	completed := atomic.LoadUint64(&c.completedTunnels)
	var avgMs float64
	if completed > 0 {
		totalNs := atomic.LoadUint64(&c.totalTunnelDurationNs)
		avgMs = float64(totalNs) / float64(completed) / 1e6
	}
	return StatsResponse{
		TotalConnections:    atomic.LoadUint64(&c.totalConnections),
		ActiveConnections:   atomic.LoadUint64(&c.activeConnections),
		TotalBytesIn:        atomic.LoadUint64(&c.totalBytesIn),
		TotalBytesOut:       atomic.LoadUint64(&c.totalBytesOut),
		AvgTunnelDurationMs: avgMs,
	}
}
