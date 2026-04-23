//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func managementAddr() string {
	return getEnv("MANAGEMENT_ADDR", "proxy-server:8080")
}

func buildinfoURL() string {
	return "http://" + managementAddr() + "/buildinfo"
}

// Scenario: UBOX-S1 returns 200 with all three fields on bare GET
func TestBuildinfo_UBOX_S1_Returns200WithAllFields(t *testing.T) {
	resp, err := http.Get(buildinfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	for _, field := range []string{"git_sha", "build_id", "go_version"} {
		if _, ok := data[field]; !ok {
			t.Errorf("missing required field %q", field)
		}
	}
}

// Scenario: UBOX-S2 git_sha reflects the value injected at build time
func TestBuildinfo_UBOX_S2_GitSHAPresent(t *testing.T) {
	resp, err := http.Get(buildinfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	sha, ok := data["git_sha"]
	if !ok {
		t.Fatal("git_sha field missing")
	}
	if sha == "" {
		t.Error("git_sha must be non-empty")
	}
	// spec: 7-char short SHA injected via ldflags; default "dev" when not injected at build time
	if len(sha) != 7 && sha != "dev" {
		t.Errorf("git_sha must be 7-char SHA or default \"dev\", got %q (len=%d)", sha, len(sha))
	}
}

// Scenario: UBOX-S3 build_id reflects BUILD_ID env when set
// Skipped when BUILD_ID is not present in the test environment.
func TestBuildinfo_UBOX_S3_BuildIDReflectsEnv(t *testing.T) {
	buildID := os.Getenv("BUILD_ID")
	if buildID == "" {
		t.Skip("BUILD_ID not set; skipping S3 (set BUILD_ID on both proxy-server and test-runner)")
	}

	resp, err := http.Get(buildinfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if data["build_id"] != buildID {
		t.Errorf("build_id: want %q (from BUILD_ID env), got %q", buildID, data["build_id"])
	}
}

// Scenario: UBOX-S4 build_id defaults to "dev" when BUILD_ID env is absent
func TestBuildinfo_UBOX_S4_BuildIDDefaultsDev(t *testing.T) {
	if os.Getenv("BUILD_ID") != "" {
		t.Skip("BUILD_ID is set; skipping S4 default test")
	}

	resp, err := http.Get(buildinfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if data["build_id"] != "dev" {
		t.Errorf("build_id: want \"dev\" when BUILD_ID unset, got %q", data["build_id"])
	}
}

// Scenario: UBOX-S5 go_version is always "go1.23"
func TestBuildinfo_UBOX_S5_GoVersion(t *testing.T) {
	resp, err := http.Get(buildinfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if data["go_version"] != "go1.23" {
		t.Errorf("go_version: want %q, got %q", "go1.23", data["go_version"])
	}
}

// Scenario: UBOX-S6 response Content-Type is application/json
func TestBuildinfo_UBOX_S6_ContentTypeJSON(t *testing.T) {
	resp, err := http.Get(buildinfoURL())
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) //nolint:errcheck

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type: want to contain %q, got %q", "application/json", ct)
	}
}

// Scenario: UBOX-S7 endpoint requires no authentication
func TestBuildinfo_UBOX_S7_NoAuthRequired(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, buildinfoURL(), nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	// Deliberately send no Authorization header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /buildinfo with no auth: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200 with no auth header, got %d", resp.StatusCode)
	}
}
