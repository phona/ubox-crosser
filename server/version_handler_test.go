package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"runtime"
	"testing"
)

func TestVersionEndpointHandler(t *testing.T) {
	GitSHA = "abc123def456abc123def456abc123def456abc1"
	defer func() { GitSHA = "unknown" }()

	t.Run("handleVersion returns 200 OK for GET request", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("GET", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("handleVersion returns 405 for POST request", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("POST", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})

	t.Run("handleVersion returns 405 for DELETE request", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("DELETE", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})

	t.Run("handleVersion returns 405 for PUT request", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("PUT", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})

	t.Run("handleVersion returns correct Content-Type", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("GET", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", contentType)
		}
	})

	t.Run("handleVersion returns valid JSON with all required fields", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("GET", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		var resp VersionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.SHA != GitSHA {
			t.Errorf("expected sha %q, got %q", GitSHA, resp.SHA)
		}
		if resp.Version != "1.0.0" {
			t.Errorf("expected version 1.0.0, got %q", resp.Version)
		}
		if resp.Module != "github.com/phona/ubox-crosser" {
			t.Errorf("expected module github.com/phona/ubox-crosser, got %q", resp.Module)
		}
		if resp.GoOS != runtime.GOOS {
			t.Errorf("expected go_os %q, got %q", runtime.GOOS, resp.GoOS)
		}
		if resp.GoArch != runtime.GOARCH {
			t.Errorf("expected go_arch %q, got %q", runtime.GOARCH, resp.GoArch)
		}
		if resp.Commit != GitSHA {
			t.Errorf("expected commit %q, got %q", GitSHA, resp.Commit)
		}
	})

	t.Run("handleVersion with unknown GitSHA", func(t *testing.T) {
		GitSHA = "unknown"
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("GET", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		var resp VersionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.SHA != "unknown" {
			t.Errorf("expected sha 'unknown', got %q", resp.SHA)
		}
		if resp.Commit != "unknown" {
			t.Errorf("expected commit 'unknown', got %q", resp.Commit)
		}
	})

	t.Run("handleVersion SHA format validation", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("GET", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		var resp VersionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		validPattern := regexp.MustCompile(`^([0-9a-f]{40}|unknown)$`)
		if !validPattern.MatchString(resp.SHA) {
			t.Errorf("sha %q doesn't match expected format (40-char hex or 'unknown')", resp.SHA)
		}
	})

	t.Run("handleVersion SHA and Commit fields match", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("GET", "/version", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		var resp VersionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.SHA != resp.Commit {
			t.Errorf("sha %q doesn't match commit %q", resp.SHA, resp.Commit)
		}
	})

	t.Run("handleVersion query parameters are ignored", func(t *testing.T) {
		ps := &ProxyServer{version: "1.0.0"}
		req := httptest.NewRequest("GET", "/version?foo=bar&baz=qux", nil)
		w := httptest.NewRecorder()

		ps.handleVersion(w, req)

		var resp VersionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.SHA != GitSHA {
			t.Errorf("query parameters affected response: got sha %q", resp.SHA)
		}
	})
}
