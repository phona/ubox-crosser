package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phona/ubox-crosser/version"
)

func TestVersionHandler_ReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	VersionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestVersionHandler_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	VersionHandler(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}

func TestVersionHandler_BodyFields(t *testing.T) {
	version.Version = "0.1.0"
	version.Commit = "abc1234"
	version.BuildTime = "2025-01-01T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	VersionHandler(rec, req)

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for _, field := range []string{"version", "commit", "build_time"} {
		if _, ok := resp[field]; !ok {
			t.Errorf("missing field %q in response", field)
		}
	}

	if resp["version"] != "0.1.0" {
		t.Errorf("version = %q, want %q", resp["version"], "0.1.0")
	}
	if resp["commit"] != "abc1234" {
		t.Errorf("commit = %q, want %q", resp["commit"], "abc1234")
	}
	if resp["build_time"] != "2025-01-01T00:00:00Z" {
		t.Errorf("build_time = %q, want %q", resp["build_time"], "2025-01-01T00:00:00Z")
	}
}

func TestVersionHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/version", nil)
	rec := httptest.NewRecorder()

	VersionHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
