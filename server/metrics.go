package server

import (
	"sync/atomic"
	"time"
)

// Metrics holds runtime counters for the proxy server.
type Metrics struct {
	startTime        time.Time
	activeConnections int64
	totalConnections  int64
	activeControllers int64
	totalLogins       int64
	totalErrors       int64
}

// MetricsSnapshot is the JSON-serializable form of Metrics.
type MetricsSnapshot struct {
	UptimeSeconds     float64 `json:"uptime_seconds"`
	ActiveConnections int64   `json:"active_connections"`
	TotalConnections  int64   `json:"total_connections"`
	ActiveControllers int64   `json:"active_controllers"`
	TotalLogins       int64   `json:"total_logins"`
	TotalErrors       int64   `json:"total_errors"`
}

// NewMetrics creates a new Metrics instance with the current time as start.
func NewMetrics() *Metrics {
	return &Metrics{
		startTime: time.Now(),
	}
}

// Snapshot returns a point-in-time copy of the metrics.
func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		UptimeSeconds:     time.Since(m.startTime).Seconds(),
		ActiveConnections: atomic.LoadInt64(&m.activeConnections),
		TotalConnections:  atomic.LoadInt64(&m.totalConnections),
		ActiveControllers: atomic.LoadInt64(&m.activeControllers),
		TotalLogins:       atomic.LoadInt64(&m.totalLogins),
		TotalErrors:       atomic.LoadInt64(&m.totalErrors),
	}
}

func (m *Metrics) AddConnection() {
	atomic.AddInt64(&m.activeConnections, 1)
	atomic.AddInt64(&m.totalConnections, 1)
}

func (m *Metrics) RemoveConnection() {
	atomic.AddInt64(&m.activeConnections, -1)
}

func (m *Metrics) AddController() {
	atomic.AddInt64(&m.activeControllers, 1)
	atomic.AddInt64(&m.totalLogins, 1)
}

func (m *Metrics) RemoveController() {
	atomic.AddInt64(&m.activeControllers, -1)
}

func (m *Metrics) AddError() {
	atomic.AddInt64(&m.totalErrors, 1)
}
