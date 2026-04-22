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

type VersionResponse struct {
	Version string `json:"version"`
	Module  string `json:"module"`
	GoOS    string `json:"go_os"`
	GoArch  string `json:"go_arch"`
	Commit  string `json:"commit"`
}

func TestVersionEndpointReturnsGitSHA(t *testing.T) {
	adminAddr := getEnv("ADMIN_SERVER_ADDR", "localhost:8080")
	versionURL := fmt.Sprintf("http://%s/version", adminAddr)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Test 1: GET /version returns 200 OK
	t.Run("returns 200 OK", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 200, got %d. body: %s", resp.StatusCode, string(body))
		}
	})

	// Test 2: Response has correct Content-Type
	t.Run("returns application/json content type", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", contentType)
		}
	})

	// Test 3: Response contains commit field with valid git SHA
	t.Run("response contains valid commit SHA", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		var versionResp VersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		// Verify commit field exists and is not empty
		if versionResp.Commit == "" {
			t.Error("commit field is empty or missing")
		}

		// Verify commit is a valid 40-character hexadecimal SHA
		hexRegex := regexp.MustCompile(`^[0-9a-f]{40}$`)
		if !hexRegex.MatchString(versionResp.Commit) {
			t.Errorf("commit value %q is not a valid 40-character hex SHA", versionResp.Commit)
		}
	})

	// Test 4: Commit SHA is consistent across multiple requests
	t.Run("commit SHA is immutable across requests", func(t *testing.T) {
		resp1, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp1.Body.Close()

		var versionResp1 VersionResponse
		if err := json.NewDecoder(resp1.Body).Decode(&versionResp1); err != nil {
			t.Fatalf("failed to decode first response: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		resp2, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp2.Body.Close()

		var versionResp2 VersionResponse
		if err := json.NewDecoder(resp2.Body).Decode(&versionResp2); err != nil {
			t.Fatalf("failed to decode second response: %v", err)
		}

		if versionResp1.Commit != versionResp2.Commit {
			t.Errorf("commit SHA changed between requests: %q vs %q",
				versionResp1.Commit, versionResp2.Commit)
		}
	})

	// Test 5: Other fields remain unchanged
	t.Run("other version fields are preserved", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		var versionResp VersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		// Module should be set
		if versionResp.Module == "" {
			t.Error("module field is empty")
		}

		// GoOS and GoArch should be set
		if versionResp.GoOS == "" {
			t.Error("go_os field is empty")
		}
		if versionResp.GoArch == "" {
			t.Error("go_arch field is empty")
		}
	})
}

func TestVersionEndpointWithCustomAdminPort(t *testing.T) {
	// This test validates that the /version endpoint works correctly
	// when the admin server is on a custom port
	adminAddr := getEnv("CUSTOM_ADMIN_ADDR", "")
	if adminAddr == "" {
		t.Skip("CUSTOM_ADMIN_ADDR not set, skipping custom port test")
	}

	versionURL := fmt.Sprintf("http://%s/version", adminAddr)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(versionURL)
	if err != nil {
		t.Fatalf("failed to GET /version on custom port: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var versionResp VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	// Verify commit is a valid 40-character hexadecimal SHA
	hexRegex := regexp.MustCompile(`^[0-9a-f]{40}$`)
	if !hexRegex.MatchString(versionResp.Commit) {
		t.Errorf("commit value %q is not a valid 40-character hex SHA", versionResp.Commit)
	}
}

func TestVersionEndpointSecurityNoLeakage(t *testing.T) {
	// This test verifies that the /version endpoint doesn't leak sensitive information
	adminAddr := getEnv("ADMIN_SERVER_ADDR", "localhost:8080")
	versionURL := fmt.Sprintf("http://%s/version", adminAddr)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(versionURL)
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	// Verify that only expected fields are present in the response
	// This is a basic check to ensure no unexpected data is leaked
	var versionResp map[string]interface{}
	if err := json.Unmarshal(body, &versionResp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Expected fields: version, module, go_os, go_arch, commit
	expectedFields := map[string]bool{
		"version": true,
		"module":  true,
		"go_os":   true,
		"go_arch": true,
		"commit":  true,
	}

	for field := range versionResp {
		if !expectedFields[field] {
			t.Warnf("unexpected field in response: %s", field)
		}
	}

	// Verify no sensitive keywords are leaked
	sensitivePatterns := []string{
		"PASSWORD",
		"TOKEN",
		"SECRET",
		"LDFLAGS",
		"buildinfo",
		"/root",
		"/home",
	}

	for _, pattern := range sensitivePatterns {
		if regexp.MustCompile(`(?i)` + pattern).MatchString(bodyStr) {
			t.Errorf("sensitive pattern found in response: %s", pattern)
		}
	}
}
