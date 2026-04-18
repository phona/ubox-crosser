package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStatsHandlerGET(t *testing.T) {
	c := NewCollector()
	c.OnTunnelStart()
	c.AddBytesIn(1024)
	c.AddBytesOut(512)
	start := time.Now().Add(-200 * time.Millisecond)
	c.OnTunnelEnd(start)

	handler := NewStatsHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var resp StatsResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.TotalConnections != 1 {
		t.Errorf("expected TotalConnections=1, got %d", resp.TotalConnections)
	}
	if resp.ActiveConnections != 0 {
		t.Errorf("expected ActiveConnections=0, got %d", resp.ActiveConnections)
	}
	if resp.TotalBytesIn != 1024 {
		t.Errorf("expected TotalBytesIn=1024, got %d", resp.TotalBytesIn)
	}
	if resp.TotalBytesOut != 512 {
		t.Errorf("expected TotalBytesOut=512, got %d", resp.TotalBytesOut)
	}
	if resp.AvgTunnelDurationMs < 100 {
		t.Errorf("expected AvgTunnelDurationMs >= 100, got %f", resp.AvgTunnelDurationMs)
	}
}

func TestStatsHandlerMethodNotAllowed(t *testing.T) {
	c := NewCollector()
	handler := NewStatsHandler(c)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/stats", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected status 405, got %d", method, rec.Code)
		}

		allow := rec.Header().Get("Allow")
		if allow != http.MethodGet {
			t.Errorf("%s: expected Allow header GET, got %s", method, allow)
		}
	}
}

func TestStatsHandlerEmptyCollector(t *testing.T) {
	c := NewCollector()
	handler := NewStatsHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp StatsResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.TotalConnections != 0 || resp.ActiveConnections != 0 ||
		resp.TotalBytesIn != 0 || resp.TotalBytesOut != 0 ||
		resp.AvgTunnelDurationMs != 0 {
		t.Errorf("expected all zeros for empty collector, got %+v", resp)
	}
}

func TestStatsHandlerNotFoundPath(t *testing.T) {
	c := NewCollector()
	handler := NewStatsHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/api/other", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for unknown path, got %d", rec.Code)
	}
}
