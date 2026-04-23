package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// HealthzResponse represents the expected response from /healthz endpoint
type HealthzResponse struct {
	Status       string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp    int64  `json:"timestamp"`
}

// TestHealthzEndpointReturns200OK verifies the endpoint returns 200 OK (FEATURE-A1)
func TestHealthzEndpointReturns200OK(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", getHealthCheckAddr()))
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

// TestHealthzEndpointReportsUptime verifies uptime is included in response (FEATURE-A2)
func TestHealthzEndpointReportsUptime(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", getHealthCheckAddr()))
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

// TestHealthzUptimeIncreases verifies uptime increases over time (FEATURE-A3)
func TestHealthzUptimeIncreases(t *testing.T) {
	healthCheckAddr := getHealthCheckAddr()

	resp1, err := http.Get(fmt.Sprintf("http://%s/healthz", healthCheckAddr))
	if err != nil {
		t.Fatalf("Failed to make first request: %v", err)
	}
	defer resp1.Body.Close()

	body1, _ := io.ReadAll(resp1.Body)
	var healthz1 HealthzResponse
	json.Unmarshal(body1, &healthz1)

	// Wait 5 seconds
	time.Sleep(5 * time.Second)

	resp2, err := http.Get(fmt.Sprintf("http://%s/healthz", healthCheckAddr))
	if err != nil {
		t.Fatalf("Failed to make second request: %v", err)
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)
	var healthz2 HealthzResponse
	json.Unmarshal(body2, &healthz2)

	uptimeDiff := healthz2.UptimeSeconds - healthz1.UptimeSeconds
	if uptimeDiff < 4 || uptimeDiff > 6 {
		t.Errorf("Expected uptime difference around 5 seconds, got %d seconds", uptimeDiff)
	}
}

// TestHealthzResponseFormat verifies response includes required fields (FEATURE-A5)
func TestHealthzResponseFormat(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", getHealthCheckAddr()))
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

// TestHealthzEndpointConcurrent verifies the endpoint handles concurrent requests (FEATURE-A6)
func TestHealthzEndpointConcurrent(t *testing.T) {
	healthCheckAddr := getHealthCheckAddr()
	done := make(chan bool, 10)
	errChan := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			resp, err := http.Get(fmt.Sprintf("http://%s/healthz", healthCheckAddr))
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

			body, _ := io.ReadAll(resp.Body)
			var healthz HealthzResponse
			if err := json.Unmarshal(body, &healthz); err != nil {
				errChan <- fmt.Errorf("invalid JSON response: %v", err)
				done <- false
				return
			}

			done <- true
		}()
	}

	successCount := 0
	for i := 0; i < 10; i++ {
		select {
		case success := <-done:
			if success {
				successCount++
			}
		case err := <-errChan:
			t.Logf("Request error: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatal("Test timeout - concurrent requests took too long")
		}
	}

	if successCount < 10 {
		t.Errorf("Expected all 10 concurrent requests to succeed, got %d", successCount)
	}
}

// getHealthCheckAddr returns the server address for HTTP endpoint tests.
// Set SERVER_ADDR env var to override (e.g. in docker-compose).
func getHealthCheckAddr() string {
	if addr := os.Getenv("SERVER_ADDR"); addr != "" {
		return addr
	}
	return "localhost:8080"
}
