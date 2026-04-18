package stats

import (
	"sync/atomic"
	"time"
)

// Collector tracks proxy connection statistics using atomic counters
// for zero-lock overhead on the data path.
type Collector struct {
	totalConnections atomic.Int64
	activeConnections atomic.Int64
	totalBytesIn     atomic.Int64
	totalBytesOut    atomic.Int64
	totalDurationMs  atomic.Int64
	closedConnections atomic.Int64
}

// NewCollector creates a new stats Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// TrackConnection records a new connection being opened.
func (c *Collector) TrackConnection() {
	c.totalConnections.Add(1)
	c.activeConnections.Add(1)
}

// ReleaseConnection records a connection being closed with its duration.
func (c *Collector) ReleaseConnection(startTime time.Time) {
	c.activeConnections.Add(-1)
	durationMs := time.Since(startTime).Milliseconds()
	c.totalDurationMs.Add(durationMs)
	c.closedConnections.Add(1)
}

// AddBytesIn records inbound bytes transferred.
func (c *Collector) AddBytesIn(n int64) {
	c.totalBytesIn.Add(n)
}

// AddBytesOut records outbound bytes transferred.
func (c *Collector) AddBytesOut(n int64) {
	c.totalBytesOut.Add(n)
}

// Snapshot holds a point-in-time copy of all stats.
type Snapshot struct {
	TotalConnections  int64 `json:"total_connections"`
	ActiveConnections int64 `json:"active_connections"`
	TotalBytesIn      int64 `json:"total_bytes_in"`
	TotalBytesOut     int64 `json:"total_bytes_out"`
	AvgTunnelDurationMs int64 `json:"avg_tunnel_duration_ms"`
}

// Snapshot returns a point-in-time copy of the stats.
func (c *Collector) Snapshot() Snapshot {
	closed := c.closedConnections.Load()
	var avgDuration int64
	if closed > 0 {
		avgDuration = c.totalDurationMs.Load() / closed
	}
	return Snapshot{
		TotalConnections:    c.totalConnections.Load(),
		ActiveConnections:   c.activeConnections.Load(),
		TotalBytesIn:        c.totalBytesIn.Load(),
		TotalBytesOut:       c.totalBytesOut.Load(),
		AvgTunnelDurationMs: avgDuration,
	}
}
