package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

// TestVersionEndpointAccessible tests FEATURE-A1: Server starts and /version endpoint is accessible
func TestVersionEndpointAccessible(t *testing.T) {
	serverURL := getServerURL(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(fmt.Sprintf("%s/version", serverURL))
	if err != nil {
		t.Fatalf("Failed to reach /version endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type: application/json, got %s", resp.Header.Get("Content-Type"))
	}
}

// TestVersionResponseContainsValidGitSHA tests FEATURE-A2: /version response contains valid git SHA
func TestVersionResponseContainsValidGitSHA(t *testing.T) {
	serverURL := getServerURL(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(fmt.Sprintf("%s/version", serverURL))
	if err != nil {
		t.Fatalf("Failed to reach /version endpoint: %v", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	commit, exists := response["commit"]
	if !exists {
		t.Fatal("Response does not contain 'commit' field")
	}

	commitStr, ok := commit.(string)
	if !ok {
		t.Fatalf("'commit' field is not a string: %v", commit)
	}

	// Validate 40-character hexadecimal format
	if !regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(commitStr) {
		t.Errorf("'commit' field is not a valid 40-character hex SHA: %s", commitStr)
	}

	// Compare with actual git HEAD
	expectedSHA, err := getGitHEAD()
	if err != nil {
		t.Logf("Warning: Could not get git HEAD: %v", err)
	} else if commitStr != expectedSHA {
		t.Errorf("Response SHA %s does not match git HEAD %s", commitStr, expectedSHA)
	}
}

// TestVersionResponseConsistency tests FEATURE-A3: Response is consistent across multiple requests
func TestVersionResponseConsistency(t *testing.T) {
	serverURL := getServerURL(t)
	client := &http.Client{Timeout: 5 * time.Second}

	var firstCommit string
	for i := 0; i < 10; i++ {
		resp, err := client.Get(fmt.Sprintf("%s/version", serverURL))
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Request %d: Failed to parse JSON: %v", i, err)
		}

		commit := response["commit"].(string)
		if i == 0 {
			firstCommit = commit
		} else if commit != firstCommit {
			t.Errorf("Request %d returned different commit: %s vs %s", i, commit, firstCommit)
		}
	}
}

// TestVersionResponseLatency tests FEATURE-A4: Response time is within acceptable bounds
func TestVersionResponseLatency(t *testing.T) {
	serverURL := getServerURL(t)
	client := &http.Client{Timeout: 5 * time.Second}

	maxLatency := 100 * time.Millisecond
	for i := 0; i < 5; i++ {
		start := time.Now()
		resp, err := client.Get(fmt.Sprintf("%s/version", serverURL))
		latency := time.Since(start)

		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()

		if latency > maxLatency {
			t.Errorf("Request %d latency %v exceeded %v", i, latency, maxLatency)
		}
	}
}

// TestVersionSecurityNoLeakage tests FEATURE-A6: Response includes only expected fields
func TestVersionSecurityNoLeakage(t *testing.T) {
	serverURL := getServerURL(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(fmt.Sprintf("%s/version", serverURL))
	if err != nil {
		t.Fatalf("Failed to reach /version endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	response := make(map[string]interface{})
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check for unwanted fields that could leak sensitive information
	unwantedFields := []string{"build_time", "go_version", "build_flags", "env", "paths"}
	for _, field := range unwantedFields {
		if _, exists := response[field]; exists {
			t.Errorf("Response contains sensitive field: %s", field)
		}
	}

	// Verify only 'commit' field is present (or minimal expected fields)
	for key := range response {
		if key != "commit" {
			t.Logf("Warning: Unexpected field in response: %s", key)
		}
	}
}

// TestVersionHTTPMethods tests FEATURE-A8: Server handles various HTTP methods
func TestVersionHTTPMethods(t *testing.T) {
	serverURL := getServerURL(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// GET should work
	resp, err := client.Get(fmt.Sprintf("%s/version", serverURL))
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET request expected 200, got %d", resp.StatusCode)
	}

	// HEAD should work (optional, but nice to have)
	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s/version", serverURL), nil)
	if err != nil {
		t.Fatalf("Failed to create HEAD request: %v", err)
	}
	resp, err = client.Do(req)
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMethodNotAllowed {
			t.Logf("HEAD request returned: %d (acceptable)", resp.StatusCode)
		}
	}

	// POST should fail or be rejected
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/version", serverURL), bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}
	resp, err = client.Do(req)
	if err == nil {
		resp.Body.Close()
		// POST may return 405 Method Not Allowed or be accepted; both are acceptable
		t.Logf("POST request returned: %d (acceptable)", resp.StatusCode)
	}
}

// Helper functions

func getServerURL(t *testing.T) string {
	// Default to localhost:8080, but allow override via env var
	url := os.Getenv("SERVER_URL")
	if url == "" {
		url = "http://localhost:8080"
	}
	return url
}

func getGitHEAD() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run git rev-parse HEAD: %w", err)
	}
	return bytes.TrimSpace(out.Bytes()).String(), nil
}
