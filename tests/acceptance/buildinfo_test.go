package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// BuildInfoResponse represents the expected /buildinfo JSON response.
type BuildInfoResponse struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

func getBuildInfoAddr() string {
	return "localhost:8080"
}

// TestBuildinfoReturns200 [UBOX-S1]
func TestBuildinfoReturns200(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("request to /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

// TestBuildinfoResponseIsJSON [UBOX-S2]
func TestBuildinfoResponseIsJSON(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("request to /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var info BuildInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Errorf("response is not valid JSON: %v\nbody: %s", err, string(body))
	}
}

// TestBuildinfoHasGitSHA [UBOX-S3]
func TestBuildinfoHasGitSHA(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("request to /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if info.GitSHA == "" {
		t.Error("git_sha field must not be empty")
	}
}

// TestBuildinfoHasBuildID [UBOX-S4]
func TestBuildinfoHasBuildID(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("request to /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if info.BuildID == "" {
		t.Error("build_id field must not be empty")
	}
}

// TestBuildinfoHasGoVersion [UBOX-S5]
func TestBuildinfoHasGoVersion(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("request to /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if info.GoVersion != "go1.23" {
		t.Errorf("expected go_version 'go1.23', got %q", info.GoVersion)
	}
}

// TestBuildinfoBuildIDDefaultsToDevWhenEnvUnset [UBOX-S6]
// When BUILD_ID env is absent the server returns build_id="dev".
// In the docker-compose acceptance stack BUILD_ID is not set, so this confirms
// the fallback behaviour.
func TestBuildinfoBuildIDDefaultsToDevWhenEnvUnset(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()))
	if err != nil {
		t.Fatalf("request to /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildInfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// In the acceptance docker-compose the container is started without BUILD_ID.
	if info.BuildID != "dev" {
		t.Errorf("expected build_id 'dev' when BUILD_ID env is absent, got %q", info.BuildID)
	}
}

// TestBuildinfoNoAuthentication [UBOX-S7]
func TestBuildinfoNoAuthentication(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/buildinfo", getBuildInfoAddr()), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		t.Errorf("/buildinfo must not require authentication, got %d", resp.StatusCode)
	}
}
