package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_GET_ReturnsCommit(t *testing.T) {
	original := Commit
	Commit = "abc123def456"
	defer func() { Commit = original }()

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()
	Handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	var resp versionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Commit != "abc123def456" {
		t.Errorf("expected commit abc123def456, got %q", resp.Commit)
	}
}

func TestHandler_GET_EmptyCommitReturnsUnknown(t *testing.T) {
	original := Commit
	Commit = ""
	defer func() { Commit = original }()

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()
	Handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp versionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Commit != "unknown" {
		t.Errorf("expected commit unknown, got %q", resp.Commit)
	}
}

func TestHandler_NonGET_Returns405(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/version", nil)
			w := httptest.NewRecorder()
			Handler(w, req)
			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s, got %d", method, w.Code)
			}
		})
	}
}

func TestHandler_GET_ResponseHasOnlyCommitField(t *testing.T) {
	original := Commit
	Commit = "test"
	defer func() { Commit = original }()

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()
	Handler(w, req)

	var raw map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	for key := range raw {
		if key != "commit" {
			t.Errorf("unexpected field %q in response", key)
		}
	}
}
