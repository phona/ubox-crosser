package stats

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleStatsGET(t *testing.T) {
	collector := NewCollector()
	collector.TrackConnection()
	collector.AddBytesIn(512)
	collector.AddBytesOut(1024)

	srv := NewServer(collector)

	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	rec := httptest.NewRecorder()
	srv.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var snap Snapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if snap.TotalConnections != 1 {
		t.Errorf("expected TotalConnections=1, got %d", snap.TotalConnections)
	}
	if snap.ActiveConnections != 1 {
		t.Errorf("expected ActiveConnections=1, got %d", snap.ActiveConnections)
	}
	if snap.TotalBytesIn != 512 {
		t.Errorf("expected TotalBytesIn=512, got %d", snap.TotalBytesIn)
	}
	if snap.TotalBytesOut != 1024 {
		t.Errorf("expected TotalBytesOut=1024, got %d", snap.TotalBytesOut)
	}
}

func TestHandleStatsMethodNotAllowed(t *testing.T) {
	collector := NewCollector()
	srv := NewServer(collector)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/stats", nil)
		rec := httptest.NewRecorder()
		srv.mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected status 405, got %d", method, rec.Code)
		}
	}
}

func TestHandleStatsResponseSchema(t *testing.T) {
	collector := NewCollector()
	srv := NewServer(collector)

	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	rec := httptest.NewRecorder()
	srv.mux.ServeHTTP(rec, req)

	var result map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	requiredFields := []string{
		"total_connections",
		"active_connections",
		"total_bytes_in",
		"total_bytes_out",
		"avg_tunnel_duration_ms",
	}
	for _, field := range requiredFields {
		if _, ok := result[field]; !ok {
			t.Errorf("missing required field: %s", field)
		}
	}
}
