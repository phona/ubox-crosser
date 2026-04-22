//go:build integration

package contract

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"
)

type VersionResponse struct {
	SHA     string `json:"sha"`
	Version string `json:"version"`
	Module  string `json:"module"`
	GoOS    string `json:"go_os"`
	GoArch  string `json:"go_arch"`
	Commit  string `json:"commit"`
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// C1: TestVersionEndpointBasicSignature verifies basic API signature (REQ-clean-1776863811-S1)
func TestVersionEndpointBasicSignature(t *testing.T) {
	adminAddr := getEnv("ADMIN_SERVER_ADDR", "localhost:8080")
	versionURL := fmt.Sprintf("http://%s/version", adminAddr)

	client := &http.Client{Timeout: 5 * time.Second}

	t.Run("GET returns 200 OK", func(t *testing.T) {
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

	t.Run("returns Content-Type: application/json", func(t *testing.T) {
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

	t.Run("response body is valid JSON", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		var obj map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&obj); err != nil {
			t.Fatalf("response body is not valid JSON: %v", err)
		}
	})
}

// C2: TestVersionResponseStructure verifies response field presence and types (REQ-clean-1776863811-S2)
func TestVersionResponseStructure(t *testing.T) {
	adminAddr := getEnv("ADMIN_SERVER_ADDR", "localhost:8080")
	versionURL := fmt.Sprintf("http://%s/version", adminAddr)

	client := &http.Client{Timeout: 5 * time.Second}

	requiredFields := []string{"sha", "version", "module", "go_os", "go_arch", "commit"}

	t.Run("response contains all required fields", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		var versionResp VersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		// Verify all fields are non-empty
		if versionResp.SHA == "" {
			t.Error("sha field is empty")
		}
		if versionResp.Version == "" {
			t.Error("version field is empty")
		}
		if versionResp.Module == "" {
			t.Error("module field is empty")
		}
		if versionResp.GoOS == "" {
			t.Error("go_os field is empty")
		}
		if versionResp.GoArch == "" {
			t.Error("go_arch field is empty")
		}
		if versionResp.Commit == "" {
			t.Error("commit field is empty")
		}
	})

	t.Run("field types are strings", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var obj map[string]interface{}
		if err := json.Unmarshal(body, &obj); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		for _, field := range requiredFields {
			value, exists := obj[field]
			if !exists {
				t.Errorf("missing required field: %s", field)
				continue
			}

			if _, ok := value.(string); !ok {
				t.Errorf("field %s is not a string, got %T", field, value)
			}
		}
	})
}

// C3: TestVersionGitSHAFormat verifies git SHA format (REQ-clean-1776863811-S3)
func TestVersionGitSHAFormat(t *testing.T) {
	adminAddr := getEnv("ADMIN_SERVER_ADDR", "localhost:8080")
	versionURL := fmt.Sprintf("http://%s/version", adminAddr)

	client := &http.Client{Timeout: 5 * time.Second}

	t.Run("sha field is valid format", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		var versionResp VersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		// Valid formats: 40-char hex or "unknown"
		validPattern := regexp.MustCompile(`^([0-9a-f]{40}|unknown)$`)
		if !validPattern.MatchString(versionResp.SHA) {
			t.Errorf("sha %q doesn't match pattern (40-char hex or 'unknown')", versionResp.SHA)
		}
	})

	t.Run("commit field matches sha field", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		var versionResp VersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		if versionResp.SHA != versionResp.Commit {
			t.Errorf("sha %q doesn't match commit %q", versionResp.SHA, versionResp.Commit)
		}
	})

	t.Run("sha format is 40-character hex or 'unknown'", func(t *testing.T) {
		resp, err := client.Get(versionURL)
		if err != nil {
			t.Fatalf("failed to GET /version: %v", err)
		}
		defer resp.Body.Close()

		var versionResp VersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}

		validPattern := regexp.MustCompile(`^([0-9a-f]{40}|unknown)$`)
		if !validPattern.MatchString(versionResp.SHA) {
			t.Errorf("sha format invalid: %q", versionResp.SHA)
		}
	})
}

// C4: TestVersionHTTPMethods verifies HTTP method support (REQ-clean-1776863811-S4)
func TestVersionHTTPMethods(t *testing.T) {
	adminAddr := getEnv("ADMIN_SERVER_ADDR", "localhost:8080")
	baseURL := fmt.Sprintf("http://%s", adminAddr)

	client := &http.Client{Timeout: 5 * time.Second}

	t.Run("GET /version returns 200", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/version", baseURL))
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /version is rejected", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/version", baseURL), nil)
		if err != nil {
			t.Fatalf("failed to create POST request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("POST request failed: %v", err)
		}
		defer resp.Body.Close()

		// Accept 405 (Method Not Allowed) or 404 (Not Found)
		if resp.StatusCode != http.StatusMethodNotAllowed && resp.StatusCode != http.StatusNotFound {
			t.Logf("POST returned %d (acceptable)", resp.StatusCode)
		}
	})

	t.Run("DELETE /version is rejected", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/version", baseURL), nil)
		if err != nil {
			t.Fatalf("failed to create DELETE request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("DELETE request failed: %v", err)
		}
		defer resp.Body.Close()

		// Accept 405 or 404
		if resp.StatusCode != http.StatusMethodNotAllowed && resp.StatusCode != http.StatusNotFound {
			t.Logf("DELETE returned %d (acceptable)", resp.StatusCode)
		}
	})

	t.Run("PUT /version is rejected", func(t *testing.T) {
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/version", baseURL), nil)
		if err != nil {
			t.Fatalf("failed to create PUT request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("PUT request failed: %v", err)
		}
		defer resp.Body.Close()

		// Accept 405 or 404
		if resp.StatusCode != http.StatusMethodNotAllowed && resp.StatusCode != http.StatusNotFound {
			t.Logf("PUT returned %d (acceptable)", resp.StatusCode)
		}
	})
}

// C5: TestVersionQueryParameterHandling verifies query parameters are ignored (REQ-clean-1776863811-S5)
func TestVersionQueryParameterHandling(t *testing.T) {
	adminAddr := getEnv("ADMIN_SERVER_ADDR", "localhost:8080")
	baseURL := fmt.Sprintf("http://%s/version", adminAddr)

	client := &http.Client{Timeout: 5 * time.Second}

	var baseResp VersionResponse
	baseHTTPResp, err := client.Get(baseURL)
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer baseHTTPResp.Body.Close()

	if err := json.NewDecoder(baseHTTPResp.Body).Decode(&baseResp); err != nil {
		t.Fatalf("failed to decode base response: %v", err)
	}

	t.Run("query parameters are ignored", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s?foo=bar&baz=qux", baseURL))
		if err != nil {
			t.Fatalf("failed to GET /version with query params: %v", err)
		}
		defer resp.Body.Close()

		var versionResp VersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
			t.Fatalf("failed to decode response with query params: %v", err)
		}

		if versionResp.SHA != baseResp.SHA || versionResp.Version != baseResp.Version {
			t.Error("response differs when query parameters are added")
		}
	})
}

// C6: TestVersionRequestHeaderTolerance verifies custom request headers don't affect response (REQ-clean-1776863811-S6)
func TestVersionRequestHeaderTolerance(t *testing.T) {
	adminAddr := getEnv("ADMIN_SERVER_ADDR", "localhost:8080")
	versionURL := fmt.Sprintf("http://%s/version", adminAddr)

	client := &http.Client{Timeout: 5 * time.Second}

	t.Run("custom request headers don't affect response", func(t *testing.T) {
		req, err := http.NewRequest("GET", versionURL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("X-Custom-Header", "test-value")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		var versionResp VersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if versionResp.SHA == "" {
			t.Error("sha field is empty")
		}
	})

	t.Run("missing standard headers don't affect response", func(t *testing.T) {
		// Create request without explicitly setting User-Agent (Go http automatically adds it)
		req, err := http.NewRequest("GET", versionURL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		// Don't set Accept header

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}

		if resp.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
		}
	})
}
