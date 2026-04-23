package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBuildInfoHandlerReturns200(t *testing.T) {
	hc := NewHealthCheckServer("localhost:8080", "abc1234")
	req := httptest.NewRequest("GET", "/buildinfo", nil)
	w := httptest.NewRecorder()

	handler := hc.makeBuildInfoHandler("abc1234")
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestBuildInfoHandlerReturnsValidJSON(t *testing.T) {
	hc := NewHealthCheckServer("localhost:8080", "abc1234")
	req := httptest.NewRequest("GET", "/buildinfo", nil)
	w := httptest.NewRecorder()

	handler := hc.makeBuildInfoHandler("abc1234")
	handler.ServeHTTP(w, req)

	var info BuildInfo
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if info.GitSha != "abc1234" {
		t.Errorf("Expected git_sha 'abc1234', got '%s'", info.GitSha)
	}
	if info.GoVersion != "go1.23" {
		t.Errorf("Expected go_version 'go1.23', got '%s'", info.GoVersion)
	}
}

func TestBuildInfoHandlerWithBuildIDEnv(t *testing.T) {
	os.Setenv("BUILD_ID", "test-build-123")
	defer os.Unsetenv("BUILD_ID")

	hc := NewHealthCheckServer("localhost:8080", "abc1234")
	req := httptest.NewRequest("GET", "/buildinfo", nil)
	w := httptest.NewRecorder()

	handler := hc.makeBuildInfoHandler("abc1234")
	handler.ServeHTTP(w, req)

	var info BuildInfo
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if info.BuildID != "test-build-123" {
		t.Errorf("Expected build_id 'test-build-123', got '%s'", info.BuildID)
	}
}

func TestBuildInfoHandlerWithoutBuildIDEnv(t *testing.T) {
	os.Unsetenv("BUILD_ID")

	hc := NewHealthCheckServer("localhost:8080", "abc1234")
	req := httptest.NewRequest("GET", "/buildinfo", nil)
	w := httptest.NewRecorder()

	handler := hc.makeBuildInfoHandler("abc1234")
	handler.ServeHTTP(w, req)

	var info BuildInfo
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if info.BuildID != "dev" {
		t.Errorf("Expected build_id 'dev', got '%s'", info.BuildID)
	}
}

func TestBuildInfoHandlerMethodNotAllowed(t *testing.T) {
	hc := NewHealthCheckServer("localhost:8080", "abc1234")
	req := httptest.NewRequest("POST", "/buildinfo", nil)
	w := httptest.NewRecorder()

	handler := hc.makeBuildInfoHandler("abc1234")
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHealthzHandlerReturns200(t *testing.T) {
	hc := NewHealthCheckServer("localhost:8080", "abc1234")
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	hc.handleHealthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHealthzHandlerReturnsValidJSON(t *testing.T) {
	hc := NewHealthCheckServer("localhost:8080", "abc1234")
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	hc.handleHealthz(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if status, ok := response["status"]; !ok || status != "healthy" {
		t.Errorf("Expected status field with value 'healthy'")
	}
	if _, ok := response["uptime_seconds"]; !ok {
		t.Errorf("Expected uptime_seconds field")
	}
	if _, ok := response["timestamp"]; !ok {
		t.Errorf("Expected timestamp field")
	}
}
