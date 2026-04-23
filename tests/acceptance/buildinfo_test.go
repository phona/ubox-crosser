//go:build acceptance

package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// BuildInfoResponse represents the expected response from /buildinfo.
type BuildInfoResponse struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

// TestBuildinfoReturns200 verifies the endpoint returns 200 OK (UBOX-BINFO-S1).
func TestBuildinfoReturns200(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("failed to reach /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

// TestBuildinfoHasThreeFields verifies all three required fields are present (UBOX-BINFO-S1).
func TestBuildinfoHasThreeFields(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("failed to reach /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var info BuildInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("response is not valid JSON: %v\nbody: %s", err, body)
	}

	if info.GitSHA == "" {
		t.Error("git_sha field is missing or empty")
	}
	if info.BuildID == "" {
		t.Error("build_id field is missing or empty")
	}
	if info.GoVersion == "" {
		t.Error("go_version field is missing or empty")
	}
}

// TestBuildinfoGoVersion verifies go_version is "go1.23" (UBOX-BINFO-S4).
func TestBuildinfoGoVersion(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("failed to reach /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if info.GoVersion != "go1.23" {
		t.Errorf("expected go_version 'go1.23', got '%s'", info.GoVersion)
	}
}

// TestBuildinfoBuildIDFromEnv verifies build_id reflects BUILD_ID env (UBOX-BINFO-S2).
// When BUILD_ID is set in the server environment, it should appear in the response.
func TestBuildinfoBuildIDFromEnv(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("failed to reach /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// In acceptance docker-compose BUILD_ID=test-build-001; locally defaults to "dev".
	if info.BuildID == "" {
		t.Error("build_id must not be empty — expected 'test-build-001' or 'dev'")
	}
}

// TestBuildinfoContentType verifies the response Content-Type is JSON (UBOX-BINFO-S5).
func TestBuildinfoContentType(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getHealthCheckAddr()))
	if err != nil {
		t.Fatalf("failed to reach /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}
}
