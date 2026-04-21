package contract

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phona/ubox-crosser/webhookdebug"
)

type RequestInfo struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func mustDecodeResponse(t *testing.T, resp *http.Response) RequestInfo {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var info RequestInfo
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	return info
}

func validateSchema(t *testing.T, info RequestInfo) {
	t.Helper()
	if info.Method == "" {
		t.Error("method field is empty")
	}
	if info.Path == "" {
		t.Error("path field is empty")
	}
	if info.Query == nil {
		t.Error("query field is nil")
	}
	if info.Headers == nil {
		t.Error("headers field is nil")
	}
}

// REQ-945-S1: POST /webhook-debug with JSON body returns 200 + correct JSON
func TestPostWithJSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(webhookdebug.Handler))
	defer srv.Close()

	reqBody := `{"event":"push"}`
	resp, err := http.Post(srv.URL+"/webhook-debug", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	info := mustDecodeResponse(t, resp)
	validateSchema(t, info)

	if info.Method != "POST" {
		t.Errorf("expected method POST, got %q", info.Method)
	}
	if info.Body != reqBody {
		t.Errorf("expected body %q, got %q", reqBody, info.Body)
	}
	if ct, ok := info.Headers["Content-Type"]; !ok || len(ct) == 0 {
		t.Error("expected Content-Type in echoed headers")
	}
}

// REQ-945-S2: GET /webhook-debug?foo=bar returns 200 + JSON with query params
func TestGetWithQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(webhookdebug.Handler))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/webhook-debug?foo=bar&baz=qux")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	info := mustDecodeResponse(t, resp)
	validateSchema(t, info)

	if info.Method != "GET" {
		t.Errorf("expected method GET, got %q", info.Method)
	}
	if vals, ok := info.Query["foo"]; !ok || len(vals) == 0 || vals[0] != "bar" {
		t.Errorf("expected query foo=bar, got %v", info.Query)
	}
	if vals, ok := info.Query["baz"]; !ok || len(vals) == 0 || vals[0] != "qux" {
		t.Errorf("expected query baz=qux, got %v", info.Query)
	}
}

// REQ-945-S3: PUT /webhook-debug returns 200 + JSON with body content
func TestPutWithBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(webhookdebug.Handler))
	defer srv.Close()

	reqBody := "key=val"
	req, _ := http.NewRequest("PUT", srv.URL+"/webhook-debug", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	info := mustDecodeResponse(t, resp)
	validateSchema(t, info)

	if info.Method != "PUT" {
		t.Errorf("expected method PUT, got %q", info.Method)
	}
	if info.Body != reqBody {
		t.Errorf("expected body %q, got %q", reqBody, info.Body)
	}
}

// REQ-945-S4: DELETE /webhook-debug returns 200 + JSON with empty body field
func TestDeleteEmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(webhookdebug.Handler))
	defer srv.Close()

	req, _ := http.NewRequest("DELETE", srv.URL+"/webhook-debug", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	info := mustDecodeResponse(t, resp)
	validateSchema(t, info)

	if info.Method != "DELETE" {
		t.Errorf("expected method DELETE, got %q", info.Method)
	}
	if info.Body != "" {
		t.Errorf("expected empty body, got %q", info.Body)
	}
}

// REQ-945-S5: Response JSON schema validation — all required fields present and typed correctly
func TestResponseSchemaValidation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(webhookdebug.Handler))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/webhook-debug")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	requiredFields := []string{"method", "path", "query", "headers", "body"}
	for _, f := range requiredFields {
		if _, ok := raw[f]; !ok {
			t.Errorf("required field %q missing from response", f)
		}
	}

	var method string
	if err := json.Unmarshal(raw["method"], &method); err != nil {
		t.Errorf("method is not a string: %v", err)
	}

	var path string
	if err := json.Unmarshal(raw["path"], &path); err != nil {
		t.Errorf("path is not a string: %v", err)
	}

	var query map[string][]string
	if err := json.Unmarshal(raw["query"], &query); err != nil {
		t.Errorf("query is not map[string][]string: %v", err)
	}

	var headers map[string][]string
	if err := json.Unmarshal(raw["headers"], &headers); err != nil {
		t.Errorf("headers is not map[string][]string: %v", err)
	}

	var bodyStr string
	if err := json.Unmarshal(raw["body"], &bodyStr); err != nil {
		t.Errorf("body is not a string: %v", err)
	}
}
