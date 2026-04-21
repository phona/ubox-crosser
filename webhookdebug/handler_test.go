package webhookdebug

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_PostWithJSONBody(t *testing.T) {
	body := `{"event":"push"}`
	req := httptest.NewRequest(http.MethodPost, "/webhook-debug", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected application/json, got %q", ct)
	}

	var info RequestInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if info.Method != "POST" {
		t.Errorf("expected method POST, got %q", info.Method)
	}
	if info.Path != "/webhook-debug" {
		t.Errorf("expected path /webhook-debug, got %q", info.Path)
	}
	if info.Body != body {
		t.Errorf("expected body %q, got %q", body, info.Body)
	}
	if _, ok := info.Headers["Content-Type"]; !ok {
		t.Error("expected Content-Type in headers")
	}
}

func TestHandler_GetWithQueryParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/webhook-debug?foo=bar&baz=qux", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	var info RequestInfo
	if err := json.NewDecoder(w.Result().Body).Decode(&info); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if info.Method != "GET" {
		t.Errorf("expected GET, got %q", info.Method)
	}
	if vals, ok := info.Query["foo"]; !ok || len(vals) == 0 || vals[0] != "bar" {
		t.Errorf("expected query foo=bar, got %v", info.Query)
	}
	if vals, ok := info.Query["baz"]; !ok || len(vals) == 0 || vals[0] != "qux" {
		t.Errorf("expected query baz=qux, got %v", info.Query)
	}
	if info.Body != "" {
		t.Errorf("expected empty body, got %q", info.Body)
	}
}

func TestHandler_PutWithFormBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/webhook-debug", strings.NewReader("key=val"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	Handler(w, req)

	var info RequestInfo
	if err := json.NewDecoder(w.Result().Body).Decode(&info); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if info.Method != "PUT" {
		t.Errorf("expected PUT, got %q", info.Method)
	}
	if info.Body != "key=val" {
		t.Errorf("expected body key=val, got %q", info.Body)
	}
}

func TestHandler_DeleteEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/webhook-debug", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	var info RequestInfo
	if err := json.NewDecoder(w.Result().Body).Decode(&info); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if info.Method != "DELETE" {
		t.Errorf("expected DELETE, got %q", info.Method)
	}
	if info.Body != "" {
		t.Errorf("expected empty body, got %q", info.Body)
	}
}

func TestHandler_ResponseSchema(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/webhook-debug", nil)
	w := httptest.NewRecorder()

	Handler(w, req)

	body, _ := io.ReadAll(w.Result().Body)
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("not valid JSON: %v", err)
	}

	for _, field := range []string{"method", "path", "query", "headers", "body"} {
		if _, ok := raw[field]; !ok {
			t.Errorf("missing required field %q", field)
		}
	}

	var method, path, bodyStr string
	var query, headers map[string][]string
	if err := json.Unmarshal(raw["method"], &method); err != nil {
		t.Errorf("method not string: %v", err)
	}
	if err := json.Unmarshal(raw["path"], &path); err != nil {
		t.Errorf("path not string: %v", err)
	}
	if err := json.Unmarshal(raw["query"], &query); err != nil {
		t.Errorf("query not map[string][]string: %v", err)
	}
	if err := json.Unmarshal(raw["headers"], &headers); err != nil {
		t.Errorf("headers not map[string][]string: %v", err)
	}
	if err := json.Unmarshal(raw["body"], &bodyStr); err != nil {
		t.Errorf("body not string: %v", err)
	}
}
