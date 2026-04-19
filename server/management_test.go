package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phona/ubox-crosser/version"
)

func TestHandleVersion_GET(t *testing.T) {
	srv := NewManagementServer(":0")
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["version"] != version.Version {
		t.Fatalf("expected %s, got %s", version.Version, result["version"])
	}
}

func TestHandleVersion_POST(t *testing.T) {
	srv := NewManagementServer(":0")
	req := httptest.NewRequest(http.MethodPost, "/version", nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestHandleVersion_VersionInjection(t *testing.T) {
	old := version.Version
	version.Version = "v2"
	defer func() { version.Version = old }()

	srv := NewManagementServer(":0")
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, req)

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["version"] != "v2" {
		t.Fatalf("expected v2, got %s", result["version"])
	}
}

func TestHandleVersion_UnknownPath(t *testing.T) {
	srv := NewManagementServer(":0")
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandleVersion_TrailingSlash(t *testing.T) {
	srv := NewManagementServer(":0")
	req := httptest.NewRequest(http.MethodGet, "/version/", nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, req)

	if w.Code == http.StatusMovedPermanently {
		t.Fatal("should not redirect /version/ with 301")
	}
}
