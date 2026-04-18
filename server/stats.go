package server

import (
	"sync/atomic"
	"time"
)

// Collector tracks proxy connection statistics using atomic counters
// for lock-free, concurrent-safe access from the data path.
type Collector struct {
	totalConnections  atomic.Int64
	activeConnections atomic.Int64
	totalBytesIn      atomic.Int64
	totalBytesOut     atomic.Int64
	totalTunnelTimeMs atomic.Int64
	completedTunnels  atomic.Int64
}

// StatsResponse is the JSON response for GET /api/stats.
type StatsResponse struct {
	TotalConnections    int64   `json:"total_connections"`
	ActiveConnections   int64   `json:"active_connections"`
	TotalBytesIn        int64   `json:"total_bytes_in"`
	TotalBytesOut       int64   `json:"total_bytes_out"`
	AvgTunnelDurationMs float64 `json:"avg_tunnel_duration_ms"`
}

// NewCollector creates a new stats Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// Snapshot returns a point-in-time copy of the stats.
func (c *Collector) Snapshot() StatsResponse {
	completed := c.completedTunnels.Load()
	var avgDuration float64
	if completed > 0 {
		avgDuration = float64(c.totalTunnelTimeMs.Load()) / float64(completed)
	}
	return StatsResponse{
		TotalConnections:    c.totalConnections.Load(),
		ActiveConnections:   c.activeConnections.Load(),
		TotalBytesIn:        c.totalBytesIn.Load(),
		TotalBytesOut:       c.totalBytesOut.Load(),
		AvgTunnelDurationMs: avgDuration,
	}
}

// OnTunnelStart records a new tunnel being established.
func (c *Collector) OnTunnelStart() {
	c.totalConnections.Add(1)
	c.activeConnections.Add(1)
}

// OnTunnelEnd records a tunnel closing and its duration.
func (c *Collector) OnTunnelEnd(start time.Time) {
	c.activeConnections.Add(-1)
	c.completedTunnels.Add(1)
	c.totalTunnelTimeMs.Add(time.Since(start).Milliseconds())
}

// AddBytesIn adds to the total inbound byte count.
func (c *Collector) AddBytesIn(n int64) {
	c.totalBytesIn.Add(n)
}

// AddBytesOut adds to the total outbound byte count.
func (c *Collector) AddBytesOut(n int64) {
	c.totalBytesOut.Add(n)
}
