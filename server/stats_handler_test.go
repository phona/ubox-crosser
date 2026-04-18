package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatsHandler_GET(t *testing.T) {
	collector := NewCollector()
	collector.RecordConnection()
	collector.RecordBytesIn(1024)

	handler := newStatsHandler(collector)
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var resp StatsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.TotalConnections != 1 {
		t.Errorf("expected TotalConnections=1, got %d", resp.TotalConnections)
	}
	if resp.TotalBytesIn != 1024 {
		t.Errorf("expected TotalBytesIn=1024, got %d", resp.TotalBytesIn)
	}
}

func TestStatsHandler_POST_MethodNotAllowed(t *testing.T) {
	collector := NewCollector()
	handler := newStatsHandler(collector)
	req := httptest.NewRequest(http.MethodPost, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestStatsHandler_ZeroCounters(t *testing.T) {
	collector := NewCollector()
	handler := newStatsHandler(collector)
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var resp StatsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.TotalConnections != 0 || resp.ActiveConnections != 0 ||
		resp.TotalBytesIn != 0 || resp.TotalBytesOut != 0 ||
		resp.AvgTunnelDurationMs != 0 {
		t.Errorf("expected all zero counters on fresh collector, got %+v", resp)
	}
}
