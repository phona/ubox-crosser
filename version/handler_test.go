package version

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_GET_ReturnsOKWithCommit(t *testing.T) {
	Commit = "abc123"
	defer func() { Commit = "" }()

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected application/json, got %q", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var vr Response
	if err := json.Unmarshal(body, &vr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if vr.Commit != "abc123" {
		t.Errorf("expected commit \"abc123\", got %q", vr.Commit)
	}
}

func TestHandler_GET_EmptyCommitDefaultsToUnknown(t *testing.T) {
	Commit = ""

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	body, _ := io.ReadAll(w.Result().Body)
	var vr Response
	if err := json.Unmarshal(body, &vr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if vr.Commit != "unknown" {
		t.Errorf("expected \"unknown\", got %q", vr.Commit)
	}
}

func TestHandler_SchemaValidation(t *testing.T) {
	Commit = "deadbeef"
	defer func() { Commit = "" }()

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	body, _ := io.ReadAll(w.Result().Body)
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("not valid JSON: %v", err)
	}

	commitVal, ok := raw["commit"]
	if !ok {
		t.Fatal("missing required field \"commit\"")
	}
	if _, ok := commitVal.(string); !ok {
		t.Errorf("\"commit\" should be string, got %T", commitVal)
	}
	if len(raw) != 1 {
		t.Errorf("expected exactly 1 field, got %d: %v", len(raw), raw)
	}
}

func TestHandler_RejectsNonGET(t *testing.T) {
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
