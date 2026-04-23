//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func getHTTPAddr() string {
	if v := os.Getenv("PROXY_HTTP_ADDR"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

type HealthzResponse struct {
	Status string `json:"status"`
	Uptime struct {
		Seconds  int64  `json:"seconds"`
		Duration string `json:"duration"`
	} `json:"uptime"`
}

// TestHealthzS1 verifies that /healthz returns a valid 200 response with expected JSON structure.
// Scenario: REQ-e2e-1776916220-S1
func TestHealthzS1(t *testing.T) {
	// Give services a moment to establish
	time.Sleep(2 * time.Second)

	baseAddr := getHTTPAddr()
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(baseAddr + "/healthz")
	if err != nil {
		t.Fatalf("failed to GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Logf("warning: expected Content-Type application/json, got %s", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var healthz HealthzResponse
	if err := json.Unmarshal(body, &healthz); err != nil {
		t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, string(body))
	}

	if healthz.Status != "healthy" {
		t.Errorf("expected status='healthy', got '%s'", healthz.Status)
	}

	if healthz.Uptime.Seconds < 0 {
		t.Errorf("expected uptime.seconds >= 0, got %d", healthz.Uptime.Seconds)
	}

	if healthz.Uptime.Duration == "" {
		t.Errorf("expected uptime.duration to be non-empty")
	}
}

// TestHealthzS2 verifies that uptime increases over time.
// Scenario: REQ-e2e-1776916220-S2
func TestHealthzS2(t *testing.T) {
	baseAddr := getHTTPAddr()
	client := &http.Client{Timeout: 5 * time.Second}

	// First request
	resp1, err := client.Get(baseAddr + "/healthz")
	if err != nil {
		t.Fatalf("failed to GET /healthz (first request): %v", err)
	}
	defer resp1.Body.Close()

	body1, _ := io.ReadAll(resp1.Body)
	var healthz1 HealthzResponse
	if err := json.Unmarshal(body1, &healthz1); err != nil {
		t.Fatalf("failed to unmarshal first response: %v", err)
	}

	// Wait 1 second
	time.Sleep(1 * time.Second)

	// Second request
	resp2, err := client.Get(baseAddr + "/healthz")
	if err != nil {
		t.Fatalf("failed to GET /healthz (second request): %v", err)
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)
	var healthz2 HealthzResponse
	if err := json.Unmarshal(body2, &healthz2); err != nil {
		t.Fatalf("failed to unmarshal second response: %v", err)
	}

	if healthz2.Uptime.Seconds <= healthz1.Uptime.Seconds {
		t.Errorf("expected uptime to increase: first=%d, second=%d", healthz1.Uptime.Seconds, healthz2.Uptime.Seconds)
	}
}

// TestHealthzS3 verifies that the duration format is a valid Go time.Duration string.
// Scenario: REQ-e2e-1776916220-S3
func TestHealthzS3(t *testing.T) {
	baseAddr := getHTTPAddr()
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(baseAddr + "/healthz")
	if err != nil {
		t.Fatalf("failed to GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var healthz HealthzResponse
	if err := json.Unmarshal(body, &healthz); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Try to parse the duration string; it should be a valid Go duration
	_, err = time.ParseDuration(healthz.Uptime.Duration)
	if err != nil {
		t.Errorf("invalid duration format '%s': %v", healthz.Uptime.Duration, err)
	}

	// Verify duration string matches seconds (with 1 second tolerance)
	if healthz.Uptime.Seconds > 0 {
		expectedDur := time.Duration(healthz.Uptime.Seconds) * time.Second
		actualDur, _ := time.ParseDuration(healthz.Uptime.Duration)

		diff := expectedDur - actualDur
		if diff < 0 {
			diff = -diff
		}
		if diff > time.Second {
			t.Errorf("duration string doesn't match seconds: %d seconds = %v, but duration string = %s",
				healthz.Uptime.Seconds, expectedDur, healthz.Uptime.Duration)
		}
	}
}
