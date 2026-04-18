//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// Contract tests derived from openspec/changes/req-03/contract.spec.yaml
// These validate schema-level constraints: field types, required fields,
// enum values, and boundary cases not covered by acceptance tests.

// --- Request schema: field types and required fields ---

func TestContract_ResponseStatusFieldIsString(t *testing.T) {
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
		t.Fatal("missing required field \"status\"")
	}
	if _, ok := status.(string); !ok {
		t.Fatalf("expected \"status\" to be a string, got %T", status)
	}
}

func TestContract_ResponseStatusEnumValue(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var hr struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &hr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	allowed := map[string]bool{"ok": true}
	if !allowed[hr.Status] {
		t.Fatalf("status %q not in allowed enum [ok]", hr.Status)
	}
}

// --- Response schema: canonical JSON ---

func TestContract_ResponseCanonicalJSON(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var parsed interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(body))), &parsed); err != nil {
		t.Fatalf("body is not valid JSON: %s", body)
	}

	canonical, _ := json.Marshal(parsed)
	expected := `{"status":"ok"}`
	if string(canonical) != expected {
		t.Fatalf("expected canonical JSON %s, got %s", expected, string(canonical))
	}
}

// --- Error responses: 405 for OPTIONS ---

func TestContract_OPTIONS_Returns405(t *testing.T) {
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

// --- Boundary: similar paths, subpaths, trailing slash ---

func TestContract_SimilarPath_Returns404(t *testing.T) {
	resp, err := http.Get(healthURL("/healthz"))
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestContract_HealthSubpath_Returns404(t *testing.T) {
	resp, err := http.Get(healthURL("/health/extra"))
	if err != nil {
		t.Fatalf("GET /health/extra failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestContract_TrailingSlash(t *testing.T) {
	resp, err := http.Get(healthURL("/health/"))
	if err != nil {
		t.Fatalf("GET /health/ failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 200 or 404 for /health/, got %d", resp.StatusCode)
	}
}

// --- Idempotency ---

func TestContract_GET_Idempotent(t *testing.T) {
	for i := 0; i < 3; i++ {
		resp, err := http.Get(healthURL("/health"))
		if err != nil {
			t.Fatalf("attempt %d: GET /health failed: %v", i+1, err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("attempt %d: expected 200, got %d", i+1, resp.StatusCode)
		}

		var hr struct {
			Status string `json:"status"`
		}
		if err := json.Unmarshal(body, &hr); err != nil {
			t.Fatalf("attempt %d: invalid JSON: %v", i+1, err)
		}
		if hr.Status != "ok" {
			t.Fatalf("attempt %d: status=%q, want \"ok\"", i+1, hr.Status)
		}
	}
}

// --- Content-Type header (dedicated check) ---

func TestContract_ContentTypeHeader(t *testing.T) {
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type starting with application/json, got %q", ct)
	}
}
