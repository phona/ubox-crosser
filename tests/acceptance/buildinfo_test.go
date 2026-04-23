//go:build acceptance

package acceptance

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
)

func buildInfoURL() string {
	if v := os.Getenv("MANAGEMENT_ADDR"); v != "" {
		return "http://" + v + "/buildinfo"
	}
	return "http://localhost:8080/buildinfo"
}

type buildInfoResponse struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

func getBuildInfo(t *testing.T) *buildInfoResponse {
	t.Helper()
	resp, err := http.Get(buildInfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()
	var d buildInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	return &d
}

// Scenario: UBOX-S1 returns 200 with all three fields on GET
func TestBuildinfo_S1_Returns200WithAllFields(t *testing.T) {
	resp, err := http.Get(buildInfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}

	var d buildInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if d.GitSHA == "" {
		t.Error("git_sha field is missing or empty")
	}
	if d.BuildID == "" {
		t.Error("build_id field is missing or empty")
	}
	if d.GoVersion == "" {
		t.Error("go_version field is missing or empty")
	}
}

// Scenario: UBOX-S2 git_sha reflects the value injected at build time
func TestBuildinfo_S2_GitSHAMatchesBuild(t *testing.T) {
	expected := os.Getenv("EXPECTED_GIT_SHA")
	if expected == "" {
		t.Skip("EXPECTED_GIT_SHA not set; set it to the 7-char SHA used at build time")
	}

	d := getBuildInfo(t)
	if d.GitSHA != expected {
		t.Errorf("git_sha: want %q, got %q", expected, d.GitSHA)
	}
}

// Scenario: UBOX-S3 build_id reflects BUILD_ID env when set
func TestBuildinfo_S3_BuildIDFromEnv(t *testing.T) {
	expected := os.Getenv("EXPECTED_BUILD_ID")
	if expected == "" {
		t.Skip("EXPECTED_BUILD_ID not set; start server with BUILD_ID=<value> and set this env")
	}

	d := getBuildInfo(t)
	if d.BuildID != expected {
		t.Errorf("build_id: want %q, got %q", expected, d.BuildID)
	}
}

// Scenario: UBOX-S4 build_id defaults to "dev" when BUILD_ID env is absent
func TestBuildinfo_S4_BuildIDDefaultsDev(t *testing.T) {
	if os.Getenv("EXPECTED_BUILD_ID") != "" {
		t.Skip("server started with BUILD_ID set; S4 (default=dev) does not apply")
	}

	d := getBuildInfo(t)
	if d.BuildID != "dev" {
		t.Errorf("build_id: want \"dev\", got %q", d.BuildID)
	}
}

// Scenario: UBOX-S5 go_version is always "go1.23"
func TestBuildinfo_S5_GoVersionIsGo123(t *testing.T) {
	d := getBuildInfo(t)
	if d.GoVersion != "go1.23" {
		t.Errorf("go_version: want \"go1.23\", got %q", d.GoVersion)
	}
}

// Scenario: UBOX-S6 response Content-Type is application/json
func TestBuildinfo_S6_ContentTypeIsJSON(t *testing.T) {
	resp, err := http.Get(buildInfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type: want to contain \"application/json\", got %q", ct)
	}
}

// Scenario: UBOX-S7 endpoint requires no authentication
func TestBuildinfo_S7_NoAuthRequired(t *testing.T) {
	resp, err := http.Get(buildInfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200 without auth header, got %d", resp.StatusCode)
	}
}
