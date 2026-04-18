package server

import (
	"sync/atomic"
)

// Collector tracks proxy server connection and traffic statistics using atomic counters.
type Collector struct {
	totalConnections      uint64
	activeConnections     uint64
	totalBytesIn          uint64
	totalBytesOut         uint64
	totalTunnelDurationNs uint64
	completedTunnels      uint64
}

// NewCollector creates a new stats Collector with all counters at zero.
func NewCollector() *Collector {
	return &Collector{}
}

// StatsResponse is the JSON response returned by GET /api/stats.
type StatsResponse struct {
	TotalConnections     uint64  `json:"total_connections"`
	ActiveConnections    uint64  `json:"active_connections"`
	TotalBytesIn         uint64  `json:"total_bytes_in"`
	TotalBytesOut        uint64  `json:"total_bytes_out"`
	AvgTunnelDurationMs  float64 `json:"avg_tunnel_duration_ms"`
}

// RecordConnection increments the total connection counter.
func (c *Collector) RecordConnection() {
	atomic.AddUint64(&c.totalConnections, 1)
}

// RecordTunnelStart increments the active connection counter.
func (c *Collector) RecordTunnelStart() {
	atomic.AddUint64(&c.activeConnections, 1)
}

// RecordTunnelEnd decrements the active connection counter and records tunnel duration.
func (c *Collector) RecordTunnelEnd(durationNs int64) {
	atomic.AddUint64(&c.activeConnections, ^uint64(0)) // decrement
	atomic.AddUint64(&c.totalTunnelDurationNs, uint64(durationNs))
	atomic.AddUint64(&c.completedTunnels, 1)
}

// RecordBytesIn adds to the total bytes read from source connections.
func (c *Collector) RecordBytesIn(n uint64) {
	atomic.AddUint64(&c.totalBytesIn, n)
}

// RecordBytesOut adds to the total bytes read from destination connections.
func (c *Collector) RecordBytesOut(n uint64) {
	atomic.AddUint64(&c.totalBytesOut, n)
}

// Snapshot returns a point-in-time copy of the stats.
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
