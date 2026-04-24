//go:build acceptance

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

// BuildinfoResponse is the expected JSON shape from GET /buildinfo per contract.spec.yaml.
type BuildinfoResponse struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

// getMgmtAddr returns the management HTTP server address.
// Override with MGMT_ADDR env; default matches acceptance docker-compose port mapping.
func getMgmtAddr() string {
	if addr := os.Getenv("MGMT_ADDR"); addr != "" {
		return addr
	}
	return "localhost:8080"
}

// Scenario UBOX-S1: returns 200 with all three fields on bare GET
func TestBuildinfoEndpoint_S1_Returns200WithAllFields(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getMgmtAddr()))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("S1: want status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("S1: read response body: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("S1: response is not valid JSON: %v\nbody: %s", err, body)
	}
	for _, field := range []string{"git_sha", "build_id", "go_version"} {
		if _, ok := raw[field]; !ok {
			t.Errorf("S1: missing required field %q", field)
		}
	}
}

// Scenario UBOX-S2: git_sha reflects the value injected at build time via -ldflags.
func TestBuildinfoEndpoint_S2_GitSHAIsSevenCharHex(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getMgmtAddr()))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildinfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("S2: decode JSON: %v", err)
	}

	matched, _ := regexp.MatchString(`^[0-9a-f]{7}$`, info.GitSHA)
	if !matched {
		t.Errorf("S2: git_sha %q must be a 7-char lowercase hex string (injected via -ldflags)", info.GitSHA)
	}
}

// Scenario UBOX-S3/S4: build_id reflects BUILD_ID env when set; defaults to "dev" otherwise.
func TestBuildinfoEndpoint_S3S4_BuildIDReflectsEnv(t *testing.T) {
	expected := os.Getenv("EXPECTED_BUILD_ID")
	if expected == "" {
		t.Log("EXPECTED_BUILD_ID not set — exercising S4: build_id should default to \"dev\"")
		expected = "dev"
	} else {
		t.Logf("EXPECTED_BUILD_ID=%s — exercising S3: build_id should reflect BUILD_ID env", expected)
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getMgmtAddr()))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildinfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("S3/S4: decode JSON: %v", err)
	}

	if info.BuildID != expected {
		t.Errorf("S3/S4: want build_id=%q, got %q", expected, info.BuildID)
	}
}

// Scenario UBOX-S5: go_version is always the literal string "go1.23".
func TestBuildinfoEndpoint_S5_GoVersionIsGo123(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getMgmtAddr()))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var info BuildinfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("S5: decode JSON: %v", err)
	}

	if info.GoVersion != "go1.23" {
		t.Errorf("S5: want go_version=%q, got %q", "go1.23", info.GoVersion)
	}
}

// Scenario UBOX-S6: response Content-Type header contains "application/json".
func TestBuildinfoEndpoint_S6_ContentTypeIsJSON(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/buildinfo", getMgmtAddr()))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	matched, _ := regexp.MatchString(`application/json`, ct)
	if !matched {
		t.Errorf("S6: want Content-Type containing \"application/json\", got %q", ct)
	}
}

// Scenario UBOX-S7: endpoint must not require authentication.
func TestBuildinfoEndpoint_S7_NoAuthRequired(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/buildinfo", getMgmtAddr()), nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /buildinfo (no auth): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("S7: want 200 without Authorization header, got %d", resp.StatusCode)
	}
}
