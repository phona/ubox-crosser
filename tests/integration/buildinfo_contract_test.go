//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
)

func managementURL(path string) string {
	addr := getEnv("MANAGEMENT_ADDR", "localhost:8080")
	return "http://" + addr + path
}

// Scenario UBOX-S1: returns 200 with all three fields on bare GET
func TestBuildinfo_S1_Returns200WithAllFields(t *testing.T) {
	resp, err := http.Get(managementURL("/buildinfo"))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}

	var d map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		t.Fatalf("response not valid JSON: %v", err)
	}

	for _, k := range []string{"git_sha", "build_id", "go_version"} {
		if _, ok := d[k]; !ok {
			t.Errorf("missing field %q in response", k)
		}
	}
}

// Scenario UBOX-S2: build_id defaults to "dev" when BUILD_ID env var is unset
func TestBuildinfo_S2_BuildIDDefaultsDev(t *testing.T) {
	if os.Getenv("EXPECTED_BUILD_ID") != "" {
		t.Skip("EXPECTED_BUILD_ID set — server has non-default BUILD_ID; S2 requires server without BUILD_ID")
	}

	resp, err := http.Get(managementURL("/buildinfo"))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	var d map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		t.Fatal(err)
	}

	got, ok := d["build_id"].(string)
	if !ok {
		t.Fatalf("build_id is not a string: %T", d["build_id"])
	}
	if got != "dev" {
		t.Errorf("want build_id=%q, got %q", "dev", got)
	}
}

// Scenario UBOX-S3: build_id reflects BUILD_ID env var when set
func TestBuildinfo_S3_BuildIDReflectsEnv(t *testing.T) {
	expected := os.Getenv("EXPECTED_BUILD_ID")
	if expected == "" {
		t.Skip("EXPECTED_BUILD_ID not set — start server with BUILD_ID=<val> and set EXPECTED_BUILD_ID=<val> to run S3")
	}

	resp, err := http.Get(managementURL("/buildinfo"))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	var d map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		t.Fatal(err)
	}

	got, ok := d["build_id"].(string)
	if !ok {
		t.Fatalf("build_id is not a string: %T", d["build_id"])
	}
	if got != expected {
		t.Errorf("want build_id=%q, got %q", expected, got)
	}
}

// Scenario UBOX-S4: go_version is always "go1.23"
func TestBuildinfo_S4_GoVersionIsGo123(t *testing.T) {
	resp, err := http.Get(managementURL("/buildinfo"))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	var d map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		t.Fatal(err)
	}

	got, ok := d["go_version"].(string)
	if !ok {
		t.Fatalf("go_version is not a string: %T", d["go_version"])
	}
	if got != "go1.23" {
		t.Errorf("want go_version=%q, got %q", "go1.23", got)
	}
}

// Scenario UBOX-S5: git_sha matches the 7-char short hash injected via ldflags
func TestBuildinfo_S5_GitShaIs7CharHex(t *testing.T) {
	resp, err := http.Get(managementURL("/buildinfo"))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	var d map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		t.Fatal(err)
	}

	sha, ok := d["git_sha"].(string)
	if !ok || sha == "" {
		t.Fatalf("git_sha missing or empty, got: %v", d["git_sha"])
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{7}$`, sha)
	if !matched {
		t.Errorf("git_sha %q is not a 7-char lowercase hex string", sha)
	}
}

// Scenario UBOX-S6: response Content-Type is application/json
func TestBuildinfo_S6_ContentTypeIsJSON(t *testing.T) {
	resp, err := http.Get(managementURL("/buildinfo"))
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("want Content-Type application/json, got %q", ct)
	}
}

// Scenario UBOX-S7: endpoint is unauthenticated (no credentials required, must return 200)
func TestBuildinfo_S7_Unauthenticated(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, managementURL("/buildinfo"), nil)
	if err != nil {
		t.Fatal(err)
	}
	// Explicitly send no Authorization header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /buildinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		t.Errorf("endpoint requires auth (got %d); spec requires no authentication", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}
