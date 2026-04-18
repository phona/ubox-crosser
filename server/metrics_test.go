package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()
	if m == nil {
		t.Fatal("NewMetrics returned nil")
	}
	snap := m.Snapshot()
	if snap.ActiveConnections != 0 {
		t.Errorf("expected ActiveConnections=0, got %d", snap.ActiveConnections)
	}
	if snap.TotalConnections != 0 {
		t.Errorf("expected TotalConnections=0, got %d", snap.TotalConnections)
	}
	if snap.ActiveControllers != 0 {
		t.Errorf("expected ActiveControllers=0, got %d", snap.ActiveControllers)
	}
	if snap.TotalLogins != 0 {
		t.Errorf("expected TotalLogins=0, got %d", snap.TotalLogins)
	}
	if snap.TotalErrors != 0 {
		t.Errorf("expected TotalErrors=0, got %d", snap.TotalErrors)
	}
}

func TestMetricsUptime(t *testing.T) {
	m := NewMetrics()
	time.Sleep(10 * time.Millisecond)
	snap := m.Snapshot()
	if snap.UptimeSeconds <= 0 {
		t.Errorf("expected positive uptime, got %f", snap.UptimeSeconds)
	}
}

func TestMetricsAddRemoveConnection(t *testing.T) {
	m := NewMetrics()
	m.AddConnection()
	m.AddConnection()
	snap := m.Snapshot()
	if snap.ActiveConnections != 2 {
		t.Errorf("expected ActiveConnections=2, got %d", snap.ActiveConnections)
	}
	if snap.TotalConnections != 2 {
		t.Errorf("expected TotalConnections=2, got %d", snap.TotalConnections)
	}

	m.RemoveConnection()
	snap = m.Snapshot()
	if snap.ActiveConnections != 1 {
		t.Errorf("expected ActiveConnections=1, got %d", snap.ActiveConnections)
	}
	if snap.TotalConnections != 2 {
		t.Errorf("expected TotalConnections still 2, got %d", snap.TotalConnections)
	}
}

func TestMetricsAddRemoveController(t *testing.T) {
	m := NewMetrics()
	m.AddController()
	m.AddController()
	m.AddController()
	snap := m.Snapshot()
	if snap.ActiveControllers != 3 {
		t.Errorf("expected ActiveControllers=3, got %d", snap.ActiveControllers)
	}
	if snap.TotalLogins != 3 {
		t.Errorf("expected TotalLogins=3, got %d", snap.TotalLogins)
	}

	m.RemoveController()
	snap = m.Snapshot()
	if snap.ActiveControllers != 2 {
		t.Errorf("expected ActiveControllers=2, got %d", snap.ActiveControllers)
	}
	if snap.TotalLogins != 3 {
		t.Errorf("expected TotalLogins still 3, got %d", snap.TotalLogins)
	}
}

func TestMetricsAddError(t *testing.T) {
	m := NewMetrics()
	m.AddError()
	m.AddError()
	snap := m.Snapshot()
	if snap.TotalErrors != 2 {
		t.Errorf("expected TotalErrors=2, got %d", snap.TotalErrors)
	}
}

func TestMetricsHandlerGET(t *testing.T) {
	m := NewMetrics()
	m.AddConnection()
	m.AddController()
	m.AddError()

	handler := metricsHandler(m)
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type application/json; charset=utf-8, got %s", ct)
	}

	var snap MetricsSnapshot
	if err := json.NewDecoder(w.Body).Decode(&snap); err != nil {
		t.Fatalf("failed to decode response: %s", err)
	}
	if snap.ActiveConnections != 1 {
		t.Errorf("expected ActiveConnections=1, got %d", snap.ActiveConnections)
	}
	if snap.TotalConnections != 1 {
		t.Errorf("expected TotalConnections=1, got %d", snap.TotalConnections)
	}
	if snap.ActiveControllers != 1 {
		t.Errorf("expected ActiveControllers=1, got %d", snap.ActiveControllers)
	}
	if snap.TotalLogins != 1 {
		t.Errorf("expected TotalLogins=1, got %d", snap.TotalLogins)
	}
	if snap.TotalErrors != 1 {
		t.Errorf("expected TotalErrors=1, got %d", snap.TotalErrors)
	}
	if snap.UptimeSeconds <= 0 {
		t.Errorf("expected positive uptime, got %f", snap.UptimeSeconds)
	}
}

func TestMetricsHandlerPOST(t *testing.T) {
	m := NewMetrics()
	handler := metricsHandler(m)
	req := httptest.NewRequest(http.MethodPost, "/api/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestMetricsHandlerResponseFields(t *testing.T) {
	m := NewMetrics()
	handler := metricsHandler(m)
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var raw map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %s", err)
	}

	requiredFields := []string{
		"uptime_seconds",
		"active_connections",
		"total_connections",
		"active_controllers",
		"total_logins",
		"total_errors",
	}
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("missing required field: %s", field)
		}
	}
}
