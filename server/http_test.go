package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/phona/ubox-crosser/server"
)

func TestHealthzHandlerReturns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	server.HealthzHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHealthzHandlerReturnsValidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	server.HealthzHandler(w, req)

	var resp struct {
		Status        string `json:"status"`
		UptimeSeconds int64  `json:"uptime_seconds"`
		Timestamp     int64  `json:"timestamp"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.Status != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", resp.Status)
	}
	if resp.UptimeSeconds < 0 {
		t.Errorf("uptime_seconds should be non-negative, got %d", resp.UptimeSeconds)
	}
	if resp.Timestamp == 0 {
		t.Error("timestamp should be a non-zero Unix value")
	}
}

func TestBuildInfoHandlerReturns200(t *testing.T) {
	handler := server.BuildInfoHandler("abc1234")
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBuildInfoHandlerReturnsThreeFields(t *testing.T) {
	t.Setenv("BUILD_ID", "test-build-42")

	handler := server.BuildInfoHandler("abc1234")
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	var resp struct {
		GitSHA    string `json:"git_sha"`
		BuildID   string `json:"build_id"`
		GoVersion string `json:"go_version"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.GitSHA != "abc1234" {
		t.Errorf("expected git_sha 'abc1234', got '%s'", resp.GitSHA)
	}
	if resp.BuildID != "test-build-42" {
		t.Errorf("expected build_id 'test-build-42', got '%s'", resp.BuildID)
	}
	if resp.GoVersion != "go1.23" {
		t.Errorf("expected go_version 'go1.23', got '%s'", resp.GoVersion)
	}
}

func TestBuildInfoHandlerDefaultsBuildIDToDev(t *testing.T) {
	os.Unsetenv("BUILD_ID")

	handler := server.BuildInfoHandler("abc1234")
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	var resp struct {
		BuildID string `json:"build_id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.BuildID != "dev" {
		t.Errorf("expected default build_id 'dev', got '%s'", resp.BuildID)
	}
}

func TestBuildInfoHandlerContentTypeJSON(t *testing.T) {
	handler := server.BuildInfoHandler("abc1234")
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}
}
