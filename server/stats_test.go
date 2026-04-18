package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollector_Snapshot_ZeroValues(t *testing.T) {
	c := NewCollector()
	snap := c.Snapshot()

	assert.Equal(t, uint64(0), snap.TotalConnections)
	assert.Equal(t, uint64(0), snap.ActiveConnections)
	assert.Equal(t, uint64(0), snap.TotalBytesIn)
	assert.Equal(t, uint64(0), snap.TotalBytesOut)
	assert.Equal(t, float64(0), snap.AvgTunnelDurationMs)
}

func TestCollector_RecordConnection(t *testing.T) {
	c := NewCollector()
	c.RecordConnection()
	c.RecordConnection()
	c.RecordConnection()

	snap := c.Snapshot()
	assert.Equal(t, uint64(3), snap.TotalConnections)
}

func TestCollector_TunnelLifecycle(t *testing.T) {
	c := NewCollector()

	c.RecordTunnelStart()
	c.RecordTunnelStart()
	snap := c.Snapshot()
	assert.Equal(t, uint64(2), snap.ActiveConnections)

	c.RecordTunnelEnd(100_000_000) // 100ms in nanoseconds
	snap = c.Snapshot()
	assert.Equal(t, uint64(1), snap.ActiveConnections)
	assert.InDelta(t, 100.0, snap.AvgTunnelDurationMs, 0.001)

	c.RecordTunnelEnd(200_000_000) // 200ms
	snap = c.Snapshot()
	assert.Equal(t, uint64(0), snap.ActiveConnections)
	assert.InDelta(t, 150.0, snap.AvgTunnelDurationMs, 0.001) // avg of 100+200
}

func TestCollector_RecordBytes(t *testing.T) {
	c := NewCollector()
	c.RecordBytesIn(1024)
	c.RecordBytesIn(2048)
	c.RecordBytesOut(512)

	snap := c.Snapshot()
	assert.Equal(t, uint64(3072), snap.TotalBytesIn)
	assert.Equal(t, uint64(512), snap.TotalBytesOut)
}

func TestCollector_ConcurrentAccess(t *testing.T) {
	c := NewCollector()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.RecordConnection()
			c.RecordTunnelStart()
			c.RecordBytesIn(10)
			c.RecordBytesOut(20)
			c.RecordTunnelEnd(1_000_000)
		}()
	}
	wg.Wait()

	snap := c.Snapshot()
	assert.Equal(t, uint64(100), snap.TotalConnections)
	assert.Equal(t, uint64(0), snap.ActiveConnections)
	assert.Equal(t, uint64(1000), snap.TotalBytesIn)
	assert.Equal(t, uint64(2000), snap.TotalBytesOut)
	assert.InDelta(t, 1.0, snap.AvgTunnelDurationMs, 0.001)
}

func TestStatsHandler_GET(t *testing.T) {
	c := NewCollector()
	c.RecordConnection()
	c.RecordBytesIn(100)

	handler := newStatsHandler(c)
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp StatsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), resp.TotalConnections)
	assert.Equal(t, uint64(100), resp.TotalBytesIn)
}

func TestStatsHandler_POST_MethodNotAllowed(t *testing.T) {
	c := NewCollector()
	handler := newStatsHandler(c)
	req := httptest.NewRequest(http.MethodPost, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestStatsHandler_PUT_MethodNotAllowed(t *testing.T) {
	c := NewCollector()
	handler := newStatsHandler(c)
	req := httptest.NewRequest(http.MethodPut, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestStatsHandler_DELETE_MethodNotAllowed(t *testing.T) {
	c := NewCollector()
	handler := newStatsHandler(c)
	req := httptest.NewRequest(http.MethodDelete, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestStatsHandler_ZeroCountersOnFreshCollector(t *testing.T) {
	c := NewCollector()
	handler := newStatsHandler(c)
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp StatsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, uint64(0), resp.TotalConnections)
	assert.Equal(t, uint64(0), resp.ActiveConnections)
	assert.Equal(t, uint64(0), resp.TotalBytesIn)
	assert.Equal(t, uint64(0), resp.TotalBytesOut)
	assert.Equal(t, float64(0), resp.AvgTunnelDurationMs)
}

func TestStatsResponse_JSONFieldNames(t *testing.T) {
	resp := StatsResponse{
		TotalConnections:    10,
		ActiveConnections:   2,
		TotalBytesIn:        1024,
		TotalBytesOut:       2048,
		AvgTunnelDurationMs: 42.5,
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	require.NoError(t, err)

	assert.Contains(t, m, "total_connections")
	assert.Contains(t, m, "active_connections")
	assert.Contains(t, m, "total_bytes_in")
	assert.Contains(t, m, "total_bytes_out")
	assert.Contains(t, m, "avg_tunnel_duration_ms")
	assert.Len(t, m, 5)
}
