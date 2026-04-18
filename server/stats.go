package server

import (
	"sync/atomic"
)

// StatsResponse is the JSON response for GET /api/stats.
type StatsResponse struct {
	TotalConnections     uint64  `json:"total_connections"`
	ActiveConnections    uint64  `json:"active_connections"`
	TotalBytesIn         uint64  `json:"total_bytes_in"`
	TotalBytesOut        uint64  `json:"total_bytes_out"`
	AvgTunnelDurationMs  float64 `json:"avg_tunnel_duration_ms"`
}

// Collector tracks proxy connection statistics using atomic counters.
type Collector struct {
	totalConnections     atomic.Uint64
	activeConnections    atomic.Uint64
	totalBytesIn         atomic.Uint64
	totalBytesOut        atomic.Uint64
	totalTunnelDurationNs atomic.Int64
	completedTunnels     atomic.Uint64
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) RecordConnection() {
	c.totalConnections.Add(1)
}

func (c *Collector) RecordTunnelStart() {
	c.activeConnections.Add(1)
}

func (c *Collector) RecordTunnelEnd(durationNs int64) {
	c.activeConnections.Add(^uint64(0)) // decrement by 1
	c.totalTunnelDurationNs.Add(durationNs)
	c.completedTunnels.Add(1)
}

func (c *Collector) RecordBytesIn(n uint64) {
	c.totalBytesIn.Add(n)
}

func (c *Collector) RecordBytesOut(n uint64) {
	c.totalBytesOut.Add(n)
}

func (c *Collector) Snapshot() StatsResponse {
	completed := c.completedTunnels.Load()
	var avgMs float64
	if completed > 0 {
		totalNs := c.totalTunnelDurationNs.Load()
		avgMs = float64(totalNs) / float64(completed) / 1e6
	}
	return StatsResponse{
		TotalConnections:    c.totalConnections.Load(),
		ActiveConnections:   c.activeConnections.Load(),
		TotalBytesIn:        c.totalBytesIn.Load(),
		TotalBytesOut:       c.totalBytesOut.Load(),
		AvgTunnelDurationMs: avgMs,
	}
}
