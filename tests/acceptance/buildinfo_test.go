// +build integration

package acceptance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
)

type BuildInfoResponse struct {
	GitSha    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

func getBuildInfoAddr() string {
	addr := os.Getenv("BUILD_INFO_ADDR")
	if addr == "" {
		addr = "proxy-server:8080"
	}
	return addr
}

// TestBuildInfoEndpointReturns200OK verifies the endpoint returns 200 OK (FEATURE-A1)
func TestBuildInfoEndpointReturns200OK(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("Failed to make request to /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestBuildInfoResponseContainsGitSha verifies git_sha is included (FEATURE-A2)
func TestBuildInfoResponseContainsGitSha(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("Failed to make request to /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	var info BuildInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if info.GitSha == "" {
		t.Errorf("git_sha should not be empty")
	}
	if len(info.GitSha) != 7 {
		t.Errorf("git_sha should be 7 characters, got %d", len(info.GitSha))
	}
}

// TestBuildInfoResponseContainsBuildId verifies build_id is included (FEATURE-A3)
func TestBuildInfoResponseContainsBuildId(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("Failed to make request to /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	var info BuildInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if info.BuildID == "" {
		t.Errorf("build_id should not be empty")
	}
}

// TestBuildInfoResponseContainsGoVersion verifies go_version is included (FEATURE-A4)
func TestBuildInfoResponseContainsGoVersion(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("Failed to make request to /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	var info BuildInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if info.GoVersion != "go1.23" {
		t.Errorf("Expected go_version 'go1.23', got '%s'", info.GoVersion)
	}
}

// TestBuildInfoResponseFormatIsValidJSON verifies all fields are present (FEATURE-A5)
func TestBuildInfoResponseFormatIsValidJSON(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("Failed to make request to /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	var info BuildInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify all three required fields are present and non-empty
	if info.GitSha == "" {
		t.Errorf("git_sha field is missing or empty")
	}
	if info.BuildID == "" {
		t.Errorf("build_id field is missing or empty")
	}
	if info.GoVersion == "" {
		t.Errorf("go_version field is missing or empty")
	}
}

// TestBuildInfoEndpointDoesNotRequireAuth verifies endpoint is unauthenticated (FEATURE-A6)
func TestBuildInfoEndpointDoesNotRequireAuth(t *testing.T) {
	// Making an unauthenticated request should still work
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("Failed to make request to /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for unauthenticated request, got %d", resp.StatusCode)
	}

	var info BuildInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify complete build information is returned
	if info.GitSha == "" || info.BuildID == "" || info.GoVersion == "" {
		t.Errorf("Incomplete build information returned for unauthenticated request")
	}
}

// TestBuildInfoEndpointHandlesConcurrentRequests verifies concurrent requests work (FEATURE-A7)
func TestBuildInfoEndpointHandlesConcurrentRequests(t *testing.T) {
	// Make 5 concurrent requests
	done := make(chan error, 5)
	for i := 0; i < 5; i++ {
		go func() {
			resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
			if err != nil {
				done <- fmt.Errorf("Failed to make request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				done <- fmt.Errorf("Expected status 200, got %d", resp.StatusCode)
				return
			}

			var info BuildInfoResponse
			if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
				done <- fmt.Errorf("Failed to decode response: %v", err)
				return
			}

			done <- nil
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 5; i++ {
		if err := <-done; err != nil {
			t.Errorf("Concurrent request %d failed: %v", i+1, err)
		}
	}
}
