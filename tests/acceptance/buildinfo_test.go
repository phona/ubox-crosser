package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
)

// BuildInfoResponse represents the expected response from /buildinfo endpoint
type BuildInfoResponse struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

// TestBuildinfoEndpointReturns200OK verifies /buildinfo returns 200 OK [ubox-crosser-S1]
func TestBuildinfoEndpointReturns200OK(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("Failed to make request to /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestBuildinfoResponseIsValidJSON verifies response is valid JSON [ubox-crosser-S6]
func TestBuildinfoResponseIsValidJSON(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var buildInfo BuildInfoResponse
	if err := json.Unmarshal(body, &buildInfo); err != nil {
		t.Errorf("Response body is not valid JSON: %v", err)
	}
}

// TestBuildinfoContainsAllRequiredFields verifies all fields are present
func TestBuildinfoContainsAllRequiredFields(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var buildInfo BuildInfoResponse
	json.Unmarshal(body, &buildInfo)

	if buildInfo.GitSHA == "" {
		t.Error("git_sha field is missing or empty")
	}

	if buildInfo.BuildID == "" {
		t.Error("build_id field is missing or empty")
	}

	if buildInfo.GoVersion == "" {
		t.Error("go_version field is missing or empty")
	}
}

// TestBuildinfoGitSHAFormat verifies git_sha is 7 hex characters [ubox-crosser-S2]
func TestBuildinfoGitSHAFormat(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var buildInfo BuildInfoResponse
	json.Unmarshal(body, &buildInfo)

	// Should match pattern: 7 hex characters
	matched, _ := regexp.MatchString(`^[0-9a-f]{7}$`, buildInfo.GitSHA)
	if !matched {
		t.Errorf("git_sha '%s' does not match pattern ^[0-9a-f]{7}$", buildInfo.GitSHA)
	}
}

// TestBuildinfoBuildIDDefault verifies build_id defaults to 'dev' [ubox-crosser-S3]
func TestBuildinfoBuildIDDefault(t *testing.T) {
	// Make sure BUILD_ID env is not set
	oldBuildID := os.Getenv("BUILD_ID")
	os.Unsetenv("BUILD_ID")
	defer os.Setenv("BUILD_ID", oldBuildID)

	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var buildInfo BuildInfoResponse
	json.Unmarshal(body, &buildInfo)

	if buildInfo.BuildID != "dev" {
		t.Errorf("Expected build_id 'dev' when env var not set, got '%s'", buildInfo.BuildID)
	}
}

// TestBuildinfoGoVersion verifies go_version is always 'go1.23' [ubox-crosser-S5]
func TestBuildinfoGoVersion(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var buildInfo BuildInfoResponse
	json.Unmarshal(body, &buildInfo)

	if buildInfo.GoVersion != "go1.23" {
		t.Errorf("Expected go_version 'go1.23', got '%s'", buildInfo.GoVersion)
	}
}

// TestBuildinfoContentType verifies Content-Type header [ubox-crosser-S8]
func TestBuildinfoContentType(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

// TestBuildinfoGETMethodOnly verifies only GET is allowed
func TestBuildinfoGETMethodOnly(t *testing.T) {
	// Test POST
	resp, err := http.Post(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()), "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 for POST, got %d", resp.StatusCode)
	}
}

// TestBuildinfoConcurrentRequests verifies endpoint handles concurrent requests
func TestBuildinfoConcurrentRequests(t *testing.T) {
	addr := getHealthCheckAddr()
	done := make(chan bool, 10)
	errChan := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", addr))
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
			var buildInfo BuildInfoResponse
			if err := json.Unmarshal(body, &buildInfo); err != nil {
				errChan <- fmt.Errorf("invalid JSON response: %v", err)
				done <- false
				return
			}

			if buildInfo.GitSHA == "" || buildInfo.BuildID == "" || buildInfo.GoVersion == "" {
				errChan <- fmt.Errorf("missing required fields")
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
		}
	}

	if successCount < 10 {
		t.Errorf("Expected all 10 concurrent requests to succeed, got %d", successCount)
	}
}

