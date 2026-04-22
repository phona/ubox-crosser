//go:build integration

package contract

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var (
	proxyAddr = getEnv("PROXY_SERVER_ADDR", "localhost:7000")
)

// TestVersionEndpointReturns200 verifies that GET /version returns HTTP 200
func TestVersionEndpointReturns200(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("http://" + proxyAddr + "/version")
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 200, got %d; body: %s", resp.StatusCode, string(body))
	}
}

// VersionResponse represents the expected response structure
type VersionResponse struct {
	SHA string `json:"sha"`
}

// TestVersionResponseContainsSHA verifies that response contains "sha" field
func TestVersionResponseContainsSHA(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("http://" + proxyAddr + "/version")
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var vr VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if vr.SHA == "" {
		t.Fatal("sha field is empty")
	}
}

// TestVersionSHAIsValidHex verifies that the SHA value is a valid hex string
func TestVersionSHAIsValidHex(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("http://" + proxyAddr + "/version")
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var vr VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify SHA is a valid 40-character hex string (full git SHA)
	hexRegex := regexp.MustCompile("^[a-f0-9]{40}$")
	if !hexRegex.MatchString(vr.SHA) {
		t.Fatalf("expected valid hex string (40 chars), got: %s", vr.SHA)
	}
}

// TestVersionContentType verifies response has correct content type
func TestVersionContentType(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("http://" + proxyAddr + "/version")
	if err != nil {
		t.Fatalf("failed to GET /version: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type: application/json, got: %s", contentType)
	}
}

// TestVersionConnectionTimeout tests connection behavior
func TestVersionConnectionTimeout(t *testing.T) {
	// Attempt TCP connection with timeout to verify server is accessible
	conn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		t.Fatalf("failed to establish TCP connection to %s: %v", proxyAddr, err)
	}
	defer conn.Close()
}
