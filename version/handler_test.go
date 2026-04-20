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

func TestHandler_DefaultValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	version.Handler(rec, req)

	var info version.Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if info.Commit != "unknown" {
		t.Errorf("expected default commit %q, got %q", "unknown", info.Commit)
	}
	if info.BuildTime != "unknown" {
		t.Errorf("expected default build_time %q, got %q", "unknown", info.BuildTime)
	}
}

func TestHandler_ExactJSONSchema(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	version.Handler(rec, req)

	var raw map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	expected := []string{"version", "commit", "build_time"}
	if len(raw) != len(expected) {
		t.Fatalf("expected %d keys, got %d: %v", len(expected), len(raw), raw)
	}
	for _, key := range expected {
		val, ok := raw[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if _, isStr := val.(string); !isStr {
			t.Errorf("key %q is not a string", key)
		}
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", version.Handler)
	return mux
}

func TestMux_MethodNotAllowed(t *testing.T) {
	mux := newMux()
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/version", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s /version, got %d", method, rec.Code)
			}
		})
	}
}

func TestMux_NotFound(t *testing.T) {
	mux := newMux()
	for _, path := range []string{"/", "/metrics", "/health"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusNotFound {
				t.Errorf("expected 404 for GET %s, got %d", path, rec.Code)
			}
		})
	}
}
