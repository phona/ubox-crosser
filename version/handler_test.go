package version_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phona/ubox-crosser/version"
)

func TestHandler_StatusAndContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	version.Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}

func TestHandler_BodyFields(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	version.Handler(rec, req)

	var info version.Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if info.Version == "" {
		t.Error("version field is empty")
	}
	if info.Commit == "" {
		t.Error("commit field is empty")
	}
	if info.BuildTime == "" {
		t.Error("build_time field is empty")
	}
}

func TestHandler_VersionValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	version.Handler(rec, req)

	var info version.Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if info.Version != "0.1.0" {
		t.Errorf("expected version 0.1.0, got %q", info.Version)
	}
}
