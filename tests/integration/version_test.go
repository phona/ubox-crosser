//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"
	"time"
)

// versionResponse represents the expected JSON response from /version endpoint
type versionResponse struct {
	SHA string `json:"sha"`
}

// TestVersionEndpointReturnsJSON verifies FEATURE-A1: /version returns valid JSON with sha field
func TestVersionEndpointReturnsJSON(t *testing.T) {
	proxyURL := fmt.Sprintf("http://%s/version", proxyAddr)

	// Give proxy server a moment to start
	time.Sleep(1 * time.Second)

	resp, err := http.Get(proxyURL)
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer resp.Body.Close()

	// FEATURE-A1: Verify status code 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Parse JSON
	var vr versionResponse
	if err := json.Unmarshal(body, &vr); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	// Verify sha field exists and is non-empty
	if vr.SHA == "" {
		t.Fatal("sha field is empty")
	}

	t.Logf("Received version response: %s", string(body))
}

// TestVersionEndpointSHAFormat verifies FEATURE-A1: sha is valid hex format
func TestVersionEndpointSHAFormat(t *testing.T) {
	proxyURL := fmt.Sprintf("http://%s/version", proxyAddr)

	resp, err := http.Get(proxyURL)
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var vr versionResponse
	if err := json.Unmarshal(body, &vr); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Verify sha is either 40-char hex (full SHA) or "unknown" (fallback)
	isValidSHA := false

	// Check for 40-character hexadecimal (full git SHA)
	if matched, _ := regexp.MatchString(`^[0-9a-f]{40}$`, vr.SHA); matched {
		isValidSHA = true
	}

	// Check for fallback "unknown"
	if vr.SHA == "unknown" {
		isValidSHA = true
	}

	if !isValidSHA {
		t.Fatalf("sha format is invalid: %s (expected 40-char hex or 'unknown')", vr.SHA)
	}

	t.Logf("SHA format valid: %s", vr.SHA)
}

// TestVersionEndpointIdempotent verifies FEATURE-A4: multiple requests return identical responses
func TestVersionEndpointIdempotent(t *testing.T) {
	proxyURL := fmt.Sprintf("http://%s/version", proxyAddr)

	const requestCount = 10
	responses := make([]string, requestCount)

	for i := 0; i < requestCount; i++ {
		resp, err := http.Get(proxyURL)
		if err != nil {
			t.Fatalf("failed to GET /version (request %d): %v", i, err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("failed to read response body (request %d): %v", i, err)
		}

		responses[i] = string(body)
	}

	// Verify all responses are identical
	for i := 1; i < requestCount; i++ {
		if responses[i] != responses[0] {
			t.Fatalf("response %d differs from response 0:\n%s\nvs\n%s", i, responses[i], responses[0])
		}
	}

	t.Logf("All %d requests returned identical responses: %s", requestCount, responses[0])
}

// TestVersionEndpointPerformance verifies FEATURE-A4: response time is fast (no I/O)
func TestVersionEndpointPerformance(t *testing.T) {
	proxyURL := fmt.Sprintf("http://%s/version", proxyAddr)

	const maxResponseTime = 100 * time.Millisecond

	start := time.Now()
	resp, err := http.Get(proxyURL)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	resp.Body.Close()

	if elapsed > maxResponseTime {
		t.Fatalf("response time %v exceeds threshold %v", elapsed, maxResponseTime)
	}

	t.Logf("Response time: %v (threshold: %v)", elapsed, maxResponseTime)
}

// TestVersionEndpointNoExtraFields verifies response contains only expected fields
func TestVersionEndpointNoExtraFields(t *testing.T) {
	proxyURL := fmt.Sprintf("http://%s/version", proxyAddr)

	resp, err := http.Get(proxyURL)
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Parse into generic object to check fields
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Verify exactly one field "sha" exists
	if len(data) != 1 {
		t.Fatalf("expected 1 field in response, got %d: %v", len(data), data)
	}

	if _, ok := data["sha"]; !ok {
		t.Fatalf("expected 'sha' field not found in response: %v", data)
	}

	t.Logf("Response structure valid: 1 field 'sha'")
}
