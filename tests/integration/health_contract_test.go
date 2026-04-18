//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

var healthAddr = getEnv("HEALTH_ADDR", "http://localhost:8080")

func healthURL(path string) string {
	return strings.TrimRight(healthAddr, "/") + path
}

// --- Response schema ---

type healthResponse struct {
	Status string `json:"status"`
}

// --- GET /health: 200 OK ---

func TestHealthEndpoint_GET_Returns200(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_GET_ContentTypeJSON(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}

func TestHealthEndpoint_GET_ResponseSchema(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var hr healthResponse
	if err := json.Unmarshal(body, &hr); err != nil {
		t.Fatalf("response is not valid JSON: %v (body: %s)", err, body)
	}

	if hr.Status != "ok" {
		t.Fatalf("expected status field to be %q, got %q", "ok", hr.Status)
	}
}

func TestHealthEndpoint_GET_ResponseHasRequiredFields(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if _, ok := raw["status"]; !ok {
		t.Fatal("response JSON missing required field \"status\"")
	}
}

func TestHealthEndpoint_GET_StatusFieldIsString(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	status, ok := raw["status"]
	if !ok {
		t.Fatal("missing \"status\" field")
	}
	if _, ok := status.(string); !ok {
		t.Fatalf("expected \"status\" to be a string, got %T", status)
	}
}

func TestHealthEndpoint_GET_StatusFieldEnumValue(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var hr healthResponse
	if err := json.Unmarshal(body, &hr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	allowed := map[string]bool{"ok": true}
	if !allowed[hr.Status] {
		t.Fatalf("status %q is not in allowed enum values [ok]", hr.Status)
	}
}

func TestHealthEndpoint_GET_NoExtraFields(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	expected := map[string]bool{"status": true}
	for k := range raw {
		if !expected[k] {
			t.Errorf("unexpected field %q in response", k)
		}
	}
}

// --- Method not allowed: 405 ---

func TestHealthEndpoint_POST_Returns405(t *testing.T) {
	resp, err := http.Post(healthURL("/health"), "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("POST /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_PUT_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, healthURL("/health"), strings.NewReader("{}"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_DELETE_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, healthURL("/health"), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_PATCH_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, healthURL("/health"), strings.NewReader("{}"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

// --- Unknown path: 404 ---

func TestHealthEndpoint_UnknownPath_Returns404(t *testing.T) {
	resp, err := http.Get(healthURL("/unknown"))
	if err != nil {
		t.Fatalf("GET /unknown failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_RootPath_Returns404(t *testing.T) {
	resp, err := http.Get(healthURL("/"))
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_SimilarPath_Returns404(t *testing.T) {
	resp, err := http.Get(healthURL("/healthz"))
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_HealthSubpath_Returns404(t *testing.T) {
	resp, err := http.Get(healthURL("/health/extra"))
	if err != nil {
		t.Fatalf("GET /health/extra failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

// --- Boundary / edge cases ---

func TestHealthEndpoint_HEAD_Returns200or405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodHead, healthURL("/health"), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HEAD /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 200 or 405 for HEAD, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_OPTIONS_Returns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodOptions, healthURL("/health"), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("OPTIONS /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_GET_ResponseBodyExactJSON(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	trimmed := strings.TrimSpace(string(body))
	var parsed interface{}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		t.Fatalf("body is not valid JSON: %s", trimmed)
	}

	re, _ := json.Marshal(parsed)
	expected := `{"status":"ok"}`
	if string(re) != expected {
		t.Fatalf("expected canonical JSON %s, got %s", expected, string(re))
	}
}

func TestHealthEndpoint_GET_Idempotent(t *testing.T) {
	for i := 0; i < 3; i++ {
		resp, err := http.Get(healthURL("/health"))
		if err != nil {
			t.Fatalf("GET /health attempt %d failed: %v", i+1, err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("attempt %d: expected 200, got %d", i+1, resp.StatusCode)
		}

		var hr healthResponse
		if err := json.Unmarshal(body, &hr); err != nil {
			t.Fatalf("attempt %d: invalid JSON: %v", i+1, err)
		}
		if hr.Status != "ok" {
			t.Fatalf("attempt %d: status=%q, want \"ok\"", i+1, hr.Status)
		}
	}
}

func TestHealthEndpoint_GET_TrailingSlash(t *testing.T) {
	resp, err := http.Get(healthURL("/health/"))
	if err != nil {
		t.Fatalf("GET /health/ failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 200 or 404 for /health/, got %d", resp.StatusCode)
	}
}
