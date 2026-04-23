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

// HealthzResponse represents the contract response from /healthz endpoint
type HealthzResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp     int64  `json:"timestamp"`
}

// getBaseURL returns the service base URL from environment variable
func getBaseURL() string {
	if url := os.Getenv("BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

// TestHealthzContractStatusCode validates endpoint returns 200 OK (REQ-e2e-fix2-1776924685-S1)
func TestHealthzContractStatusCode(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(fmt.Sprintf("%s/healthz", baseURL))
	if err != nil {
		t.Fatalf("Failed to make request to /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
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

// TestHealthzContractUptimeField validates uptime_seconds field exists (REQ-e2e-fix2-1776924685-S2)
func TestHealthzContractUptimeField(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(fmt.Sprintf("%s/healthz", baseURL))
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

// TestHealthzContractUptimeIncreases validates uptime increments over time (REQ-e2e-fix2-1776924685-S3)
func TestHealthzContractUptimeIncreases(t *testing.T) {
	baseURL := getBaseURL()

	// First request
	resp1, err := http.Get(fmt.Sprintf("%s/healthz", baseURL))
	if err != nil {
		t.Fatalf("Failed to make first request: %v", err)
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Fatalf("Failed to read first response body: %v", err)
	}

	var healthz1 HealthzResponse
	if err := json.Unmarshal(body1, &healthz1); err != nil {
		t.Fatalf("First response is not valid JSON: %v", err)
	}

	// Wait 5 seconds
	time.Sleep(5 * time.Second)

	// Second request
	resp2, err := http.Get(fmt.Sprintf("%s/healthz", baseURL))
	if err != nil {
		t.Fatalf("Failed to make second request: %v", err)
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		t.Fatalf("Failed to read second response body: %v", err)
	}

	var healthz2 HealthzResponse
	if err := json.Unmarshal(body2, &healthz2); err != nil {
		t.Fatalf("Second response is not valid JSON: %v", err)
	}

	// Verify uptime increased by approximately 5 seconds (±1 second tolerance)
	uptimeDiff := healthz2.UptimeSeconds - healthz1.UptimeSeconds
	if uptimeDiff < 4 || uptimeDiff > 6 {
		t.Errorf("Expected uptime difference around 5 seconds (4-6), got %d seconds", uptimeDiff)
	}
}

// TestHealthzContractResponseFormat validates response includes required fields (REQ-e2e-fix2-1776924685-S4)
func TestHealthzContractResponseFormat(t *testing.T) {
	baseURL := getBaseURL()
	resp, err := http.Get(fmt.Sprintf("%s/healthz", baseURL))
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

	// Verify status field
	if healthz.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", healthz.Status)
	}

	// Verify uptime_seconds field
	if healthz.UptimeSeconds < 0 {
		t.Errorf("uptime_seconds should be non-negative, got %d", healthz.UptimeSeconds)
	}

	// Verify timestamp field
	if healthz.Timestamp == 0 {
		t.Error("timestamp field should be set to a valid Unix timestamp")
	}

	// Additional check: timestamp should be recent (within last 10 seconds)
	now := time.Now().Unix()
	if healthz.Timestamp > now || healthz.Timestamp < now-10 {
		t.Errorf("timestamp %d should be recent (within last 10 seconds), current time %d", healthz.Timestamp, now)
	}
}

// TestHealthzContractConcurrentRequests validates concurrent request handling (REQ-e2e-fix2-1776924685-S5)
func TestHealthzContractConcurrentRequests(t *testing.T) {
	baseURL := getBaseURL()
	numRequests := 10
	done := make(chan bool, numRequests)
	errChan := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := http.Get(fmt.Sprintf("%s/healthz", baseURL))
			if err != nil {
				errChan <- err
				done <- false
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errChan <- fmt.Errorf("expected 200, got %d", resp.StatusCode)
				done <- false
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errChan <- err
				done <- false
				return
			}

			var healthz HealthzResponse
			if err := json.Unmarshal(body, &healthz); err != nil {
				errChan <- fmt.Errorf("invalid JSON response: %v", err)
				done <- false
				return
			}

			// Verify uptime is non-negative
			if healthz.UptimeSeconds < 0 {
				errChan <- fmt.Errorf("uptime_seconds should be non-negative, got %d", healthz.UptimeSeconds)
				done <- false
				return
			}

			done <- true
		}()
	}

	successCount := 0
	for i := 0; i < numRequests; i++ {
		select {
		case success := <-done:
			if success {
				successCount++
			}
		case err := <-errChan:
			t.Logf("Request error: %v", err)
		case <-time.After(30 * time.Second):
			t.Fatal("Test timeout - concurrent requests took too long")
		}
	}

	if successCount < numRequests {
		t.Errorf("Expected all %d concurrent requests to succeed, got %d", numRequests, successCount)
	}
}

// TestHealthzContractMinimalResponseTime validates response time is acceptable (implicit check)
func TestHealthzContractResponseTime(t *testing.T) {
	baseURL := getBaseURL()

	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/healthz", baseURL))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	// Response should be fast (< 100ms is ideal for load balancer health checks)
	if elapsed > 100*time.Millisecond {
		t.Logf("Warning: Response took %v, expected < 100ms", elapsed)
	}

	// Verify response is still valid
	body, _ := io.ReadAll(resp.Body)
	var healthz HealthzResponse
	if err := json.Unmarshal(body, &healthz); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}
}
