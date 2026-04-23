//go:build integration

package contract

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// HealthzResponse represents the expected API contract for /healthz endpoint
type HealthzResponse struct {
	Status       string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp    int64  `json:"timestamp"`
}

func getHealthCheckAddr() string {
	// Read from environment variable set by Docker Compose
	addr := os.Getenv("HEALTH_CHECK_ADDR")
	if addr == "" {
		addr = "localhost:8080"
	}
	return addr
}

// TestREQe2efixS1HealthCheckEndpointReturns200OK validates endpoint returns 200 OK
func TestREQe2efixS1HealthCheckEndpointReturns200OK(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("http://%s/healthz", getHealthCheckAddr())
	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request to /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var healthz HealthzResponse
	if err := json.Unmarshal(body, &healthz); err != nil {
		t.Errorf("Response body is not valid JSON: %v", err)
	}
}

// TestREQe2efixS2EndpointReportsServiceUptimeInSeconds validates uptime field exists and is reasonable
func TestREQe2efixS2EndpointReportsServiceUptimeInSeconds(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("http://%s/healthz", getHealthCheckAddr())
	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request to /healthz: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var healthz HealthzResponse
	if err := json.Unmarshal(body, &healthz); err != nil {
		t.Fatalf("Response body is not valid JSON: %v", err)
	}

	if healthz.UptimeSeconds < 0 {
		t.Errorf("uptime_seconds should be non-negative, got %d", healthz.UptimeSeconds)
	}
}

// TestREQe2efixS3UptimeIncreasesOverTime validates uptime increments correctly
func TestREQe2efixS3UptimeIncreasesOverTime(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	healthCheckAddr := getHealthCheckAddr()
	url := fmt.Sprintf("http://%s/healthz", healthCheckAddr)

	// First request
	resp1, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to make first request: %v", err)
	}
	defer resp1.Body.Close()

	body1, _ := io.ReadAll(resp1.Body)
	var healthz1 HealthzResponse
	if err := json.Unmarshal(body1, &healthz1); err != nil {
		t.Fatalf("First response is not valid JSON: %v", err)
	}

	// Wait 5 seconds
	time.Sleep(5 * time.Second)

	// Second request
	resp2, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to make second request: %v", err)
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)
	var healthz2 HealthzResponse
	if err := json.Unmarshal(body2, &healthz2); err != nil {
		t.Fatalf("Second response is not valid JSON: %v", err)
	}

	uptimeDiff := healthz2.UptimeSeconds - healthz1.UptimeSeconds
	// Allow ±1 second tolerance
	if uptimeDiff < 4 || uptimeDiff > 6 {
		t.Errorf("Expected uptime difference around 5 seconds (±1), got %d seconds", uptimeDiff)
	}
}

// TestREQe2efixS5ResponseFormatIncludesStatusField validates JSON structure
func TestREQe2efixS5ResponseFormatIncludesStatusField(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("http://%s/healthz", getHealthCheckAddr())
	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var healthz HealthzResponse
	if err := json.Unmarshal(body, &healthz); err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	if healthz.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", healthz.Status)
	}

	if healthz.UptimeSeconds < 0 {
		t.Errorf("uptime_seconds should be non-negative, got %d", healthz.UptimeSeconds)
	}

	if healthz.Timestamp == 0 {
		t.Error("timestamp field should be set to a valid Unix timestamp")
	}
}

// TestREQe2efixS6EndpointHandlesConcurrentRequests validates concurrent request handling
func TestREQe2efixS6EndpointHandlesConcurrentRequests(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	healthCheckAddr := getHealthCheckAddr()
	url := fmt.Sprintf("http://%s/healthz", healthCheckAddr)
	numRequests := 10

	done := make(chan bool, numRequests)
	errChan := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(index int) {
			defer func() { done <- true }()

			resp, err := client.Get(url)
			if err != nil {
				errChan <- fmt.Errorf("Request %d failed: %v", index, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errChan <- fmt.Errorf("Request %d: expected 200, got %d", index, resp.StatusCode)
				return
			}

			body, _ := io.ReadAll(resp.Body)
			var healthz HealthzResponse
			if err := json.Unmarshal(body, &healthz); err != nil {
				errChan <- fmt.Errorf("Request %d: invalid JSON response: %v", index, err)
				return
			}

			if healthz.UptimeSeconds < 0 {
				errChan <- fmt.Errorf("Request %d: negative uptime_seconds", index)
				return
			}
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		select {
		case <-done:
			// Request completed
		case err := <-errChan:
			t.Errorf("Concurrent request error: %v", err)
		case <-time.After(30 * time.Second):
			t.Fatal("Test timeout - concurrent requests took too long")
		}
	}
}
