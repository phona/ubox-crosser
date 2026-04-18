package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCollector_Snapshot_ZeroState(t *testing.T) {
	c := NewCollector()
	snap := c.Snapshot()

	if snap.TotalConnections != 0 {
		t.Errorf("TotalConnections = %d, want 0", snap.TotalConnections)
	}
	if snap.ActiveConnections != 0 {
		t.Errorf("ActiveConnections = %d, want 0", snap.ActiveConnections)
	}
	if snap.TotalBytesIn != 0 {
		t.Errorf("TotalBytesIn = %d, want 0", snap.TotalBytesIn)
	}
	if snap.TotalBytesOut != 0 {
		t.Errorf("TotalBytesOut = %d, want 0", snap.TotalBytesOut)
	}
	if snap.AvgTunnelDurationMs != 0 {
		t.Errorf("AvgTunnelDurationMs = %f, want 0", snap.AvgTunnelDurationMs)
	}
}

func TestCollector_RecordConnection(t *testing.T) {
	c := NewCollector()
	c.RecordConnection()
	c.RecordConnection()
	c.RecordConnection()

	snap := c.Snapshot()
	if snap.TotalConnections != 3 {
		t.Errorf("TotalConnections = %d, want 3", snap.TotalConnections)
	}
}

func TestCollector_TunnelLifecycle(t *testing.T) {
	c := NewCollector()

	c.RecordTunnelStart()
	snap := c.Snapshot()
	if snap.ActiveConnections != 1 {
		t.Errorf("ActiveConnections after start = %d, want 1", snap.ActiveConnections)
	}

	c.RecordTunnelStart()
	snap = c.Snapshot()
	if snap.ActiveConnections != 2 {
		t.Errorf("ActiveConnections after 2 starts = %d, want 2", snap.ActiveConnections)
	}

	c.RecordTunnelEnd(100 * time.Millisecond)
	snap = c.Snapshot()
	if snap.ActiveConnections != 1 {
		t.Errorf("ActiveConnections after 1 end = %d, want 1", snap.ActiveConnections)
	}
	if snap.AvgTunnelDurationMs < 99 || snap.AvgTunnelDurationMs > 101 {
		t.Errorf("AvgTunnelDurationMs = %f, want ~100", snap.AvgTunnelDurationMs)
	}

	c.RecordTunnelEnd(200 * time.Millisecond)
	snap = c.Snapshot()
	if snap.ActiveConnections != 0 {
		t.Errorf("ActiveConnections after all end = %d, want 0", snap.ActiveConnections)
	}
	if snap.AvgTunnelDurationMs < 149 || snap.AvgTunnelDurationMs > 151 {
		t.Errorf("AvgTunnelDurationMs = %f, want ~150", snap.AvgTunnelDurationMs)
	}
}

func TestCollector_RecordBytes(t *testing.T) {
	c := NewCollector()
	c.RecordBytesIn(100)
	c.RecordBytesIn(200)
	c.RecordBytesOut(50)

	snap := c.Snapshot()
	if snap.TotalBytesIn != 300 {
		t.Errorf("TotalBytesIn = %d, want 300", snap.TotalBytesIn)
	}
	if snap.TotalBytesOut != 50 {
		t.Errorf("TotalBytesOut = %d, want 50", snap.TotalBytesOut)
	}
}

func TestStatsHandler_GET(t *testing.T) {
	c := NewCollector()
	c.RecordConnection()
	c.RecordBytesIn(1024)

	handler := NewStatsHandler(c)
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/stats status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var resp StatsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.TotalConnections != 1 {
		t.Errorf("TotalConnections = %d, want 1", resp.TotalConnections)
	}
	if resp.TotalBytesIn != 1024 {
		t.Errorf("TotalBytesIn = %d, want 1024", resp.TotalBytesIn)
	}
}

func TestStatsHandler_POST_MethodNotAllowed(t *testing.T) {
	handler := NewStatsHandler(NewCollector())
	req := httptest.NewRequest(http.MethodPost, "/api/stats", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /api/stats status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestStatsHandler_JSONSchema(t *testing.T) {
	handler := NewStatsHandler(NewCollector())
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var raw map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("decode: %v", err)
	}

	expectedKeys := []string{"total_connections", "active_connections", "total_bytes_in", "total_bytes_out", "avg_tunnel_duration_ms"}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("missing JSON key %q", key)
		}
	}
	if len(raw) != len(expectedKeys) {
		t.Errorf("expected %d keys, got %d: %v", len(expectedKeys), len(raw), raw)
	}
}
