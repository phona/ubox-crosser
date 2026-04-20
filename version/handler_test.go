package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsOKWithJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	var info Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
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

func TestHandler_VersionMatchesConstant(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	var info Info
	json.NewDecoder(rec.Body).Decode(&info)

	if info.Version != Version {
		t.Errorf("expected version %q, got %q", Version, info.Version)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/version", nil)
		rec := httptest.NewRecorder()

		Handler(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected status 405, got %d", method, rec.Code)
		}
	}
}
