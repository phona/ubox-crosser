//go:build contract

package contract

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var healthBaseURL = getEnv("HEALTH_BASE_URL", "http://localhost:8080")

func healthURL(path string) string {
	return healthBaseURL + path
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

// ---------------------------------------------------------------------------
// GET /health — 200 OK
// ---------------------------------------------------------------------------

func TestGetHealth_Returns200(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetHealth_ContentTypeJSON(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}

func TestGetHealth_ResponseSchema(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v (body: %s)", err, string(body))
	}

	status, ok := result["status"]
	if !ok {
		t.Fatal("response JSON missing required field 'status'")
	}

	statusStr, ok := status.(string)
	if !ok {
		t.Fatalf("field 'status' is not a string, got %T", status)
	}

	if statusStr != "ok" {
		t.Fatalf("expected status='ok', got %q", statusStr)
	}
}

func TestGetHealth_NoExtraFields(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected exactly 1 field in response, got %d: %v", len(result), result)
	}
}

func TestGetHealth_BodyExactMatch(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("body is not valid JSON: %v", err)
	}

	expected := map[string]interface{}{"status": "ok"}
	reEncoded, _ := json.Marshal(parsed)
	expectedEncoded, _ := json.Marshal(expected)
	if string(reEncoded) != string(expectedEncoded) {
		t.Fatalf("body mismatch: got %s, want %s", string(reEncoded), string(expectedEncoded))
	}
}

// ---------------------------------------------------------------------------
// Non-GET methods on /health — 405 Method Not Allowed
// ---------------------------------------------------------------------------

func TestPostHealth_Returns405(t *testing.T) {
	resp, err := httpClient.Post(healthURL("/health"), "application/json", nil)
	if err != nil {
		t.Fatalf("POST /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestPostHealth_AllowHeader(t *testing.T) {
	resp, err := httpClient.Post(healthURL("/health"), "application/json", nil)
	if err != nil {
		t.Fatalf("POST /health failed: %v", err)
	}
	defer resp.Body.Close()

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

func TestPutHealth_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, healthURL("/health"), nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestPutHealth_AllowHeader(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, healthURL("/health"), nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /health failed: %v", err)
	}
	defer resp.Body.Close()

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

func TestDeleteHealth_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, healthURL("/health"), nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestDeleteHealth_AllowHeader(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, healthURL("/health"), nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /health failed: %v", err)
	}
	defer resp.Body.Close()

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

func TestPatchHealth_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, healthURL("/health"), nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestPatchHealth_AllowHeader(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, healthURL("/health"), nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH /health failed: %v", err)
	}
	defer resp.Body.Close()

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

// ---------------------------------------------------------------------------
// Unknown paths — 404 Not Found
// ---------------------------------------------------------------------------

func TestRootPath_Returns404(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/"))
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestMetricsPath_Returns404(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/metrics"))
	if err != nil {
		t.Fatalf("GET /metrics failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestRandomPath_Returns404(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/foo/bar/baz"))
	if err != nil {
		t.Fatalf("GET /foo/bar/baz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestHealthzPath_Returns404(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/healthz"))
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// Boundary / edge cases
// ---------------------------------------------------------------------------

func TestGetHealthTrailingSlash_Returns404(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/health/"))
	if err != nil {
		t.Fatalf("GET /health/ failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("/health/ should return 404 (only exact /health is valid), got %d", resp.StatusCode)
	}
}

func TestGetHealthCaseSensitive_Returns404(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/Health"))
	if err != nil {
		t.Fatalf("GET /Health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("/Health should return 404 (path is case-sensitive), got %d", resp.StatusCode)
	}
}

func TestHeadHealth_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodHead, healthURL("/health"), nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("HEAD /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 for HEAD, got %d", resp.StatusCode)
	}
}

func TestOptionsHealth_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodOptions, healthURL("/health"), nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("OPTIONS /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 for OPTIONS, got %d", resp.StatusCode)
	}
}

func TestPostUnknownPath_Returns404(t *testing.T) {
	resp, err := httpClient.Post(healthURL("/unknown"), "application/json", nil)
	if err != nil {
		t.Fatalf("POST /unknown failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestGetHealth_ResponseNotEmpty(t *testing.T) {
	resp, err := httpClient.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("response body should not be empty")
	}
}
