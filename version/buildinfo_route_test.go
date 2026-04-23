package version_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phona/ubox-crosser/version"
)

func newAdminMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", version.Handler)
	mux.HandleFunc("GET /buildinfo", version.Handler)
	return mux
}

func TestBuildinfo_StatusAndContentType(t *testing.T) {
	mux := newAdminMux()
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}

func TestBuildinfo_BodyFields(t *testing.T) {
	mux := newAdminMux()
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

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

func TestBuildinfo_ExactJSONSchema(t *testing.T) {
	mux := newAdminMux()
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

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

func TestBuildinfo_MethodNotAllowed(t *testing.T) {
	mux := newAdminMux()
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/buildinfo", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s /buildinfo, got %d", method, rec.Code)
			}
		})
	}
}

func TestBuildinfo_MatchesVersion(t *testing.T) {
	mux := newAdminMux()

	reqV := httptest.NewRequest(http.MethodGet, "/version", nil)
	recV := httptest.NewRecorder()
	mux.ServeHTTP(recV, reqV)

	reqB := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	recB := httptest.NewRecorder()
	mux.ServeHTTP(recB, reqB)

	if recV.Code != recB.Code {
		t.Fatalf("status mismatch: /version=%d, /buildinfo=%d", recV.Code, recB.Code)
	}
	if recV.Body.String() != recB.Body.String() {
		t.Errorf("body mismatch:\n  /version:   %s\n  /buildinfo: %s", recV.Body.String(), recB.Body.String())
	}
}
