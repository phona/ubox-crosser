package healthz

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/healthz", nil)
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

	var hr Response
	if err := json.Unmarshal(body, &hr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if hr.Status != "ok" {
		t.Errorf("expected status \"ok\", got %q", hr.Status)
	}
}

func TestHandler_SchemaValidation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/healthz", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	body, _ := io.ReadAll(w.Result().Body)
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("not valid JSON: %v", err)
	}

	statusVal, ok := raw["status"]
	if !ok {
		t.Fatal("missing required field \"status\"")
	}
	if _, ok := statusVal.(string); !ok {
		t.Errorf("\"status\" should be string, got %T", statusVal)
	}
	if len(raw) != 1 {
		t.Errorf("expected exactly 1 field, got %d: %v", len(raw), raw)
	}
}

func TestHandler_RejectsNonGet(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/healthz", nil)
			w := httptest.NewRecorder()

			Handler(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s, got %d", method, w.Code)
			}
		})
	}
}
