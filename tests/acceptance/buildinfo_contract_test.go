//go:build acceptance

package acceptance

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
)

func managementAddr() string {
	if v := os.Getenv("MANAGEMENT_ADDR"); v != "" {
		return v
	}
	return "http://ubox-crosser:8080"
}

type buildinfoResponse struct {
	GitSHA     string `json:"git_sha"`
	BuildID    string `json:"build_id"`
	GoVersion  string `json:"go_version"`
}

func getBuildinfoContract(t *testing.T) (*http.Response, buildinfoResponse) {
	t.Helper()
	url := managementAddr() + "/buildinfo"
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		t.Fatalf("GET /buildinfo failed: %v", err)
	}
	t.Cleanup(func() { resp.Body.Close() })

	var body buildinfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	return resp, body
}

// Scenario: UBOX-S1 returns 200 with all three fields on bare GET
func TestBuildinfoContract_S1_Returns200WithAllFields(t *testing.T) {
	resp, body := getBuildinfoContract(t)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("UBOX-S1: want status 200, got %d", resp.StatusCode)
	}
	if body.GitSHA == "" {
		t.Error("UBOX-S1: missing or empty field git_sha")
	}
	if body.BuildID == "" {
		t.Error("UBOX-S1: missing or empty field build_id")
	}
	if body.GoVersion == "" {
		t.Error("UBOX-S1: missing or empty field go_version")
	}
}

// Scenario: UBOX-S2 git_sha reflects the value injected at build time.
// In black-box mode, assert the value is a 7-character alphanumeric string.
// Set EXPECTED_GIT_SHA env to assert an exact value when the expected SHA is known.
func TestBuildinfoContract_S2_GitSHAReflectsBuildInjection(t *testing.T) {
	_, body := getBuildinfoContract(t)

	if expected := os.Getenv("EXPECTED_GIT_SHA"); expected != "" {
		if body.GitSHA != expected {
			t.Errorf("UBOX-S2: want git_sha=%q, got %q", expected, body.GitSHA)
		}
		return
	}

	// Without a known expected value, assert the format: exactly 7 hex chars.
	if len(body.GitSHA) != 7 {
		t.Errorf("UBOX-S2: git_sha must be 7 characters, got %q (len=%d)", body.GitSHA, len(body.GitSHA))
	}
	for _, c := range body.GitSHA {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Errorf("UBOX-S2: git_sha contains non-hex character %q in %q", c, body.GitSHA)
		}
	}
}

// Scenario: UBOX-S3 build_id reflects BUILD_ID env when set.
// Set EXPECTED_BUILD_ID env to a non-empty, non-"dev" value to exercise this scenario.
// If EXPECTED_BUILD_ID is absent, this test is skipped (see UBOX-S4).
func TestBuildinfoContract_S3_BuildIDReflectsEnv(t *testing.T) {
	expected := os.Getenv("EXPECTED_BUILD_ID")
	if expected == "" || expected == "dev" {
		t.Skip("UBOX-S3: set EXPECTED_BUILD_ID to a custom value to exercise this scenario")
	}

	_, body := getBuildinfoContract(t)

	if body.BuildID != expected {
		t.Errorf("UBOX-S3: want build_id=%q (from BUILD_ID env), got %q", expected, body.BuildID)
	}
}

// Scenario: UBOX-S4 build_id defaults to "dev" when BUILD_ID env is absent.
// This test runs when EXPECTED_BUILD_ID is not set or is "dev".
func TestBuildinfoContract_S4_BuildIDDefaultsDev(t *testing.T) {
	if v := os.Getenv("EXPECTED_BUILD_ID"); v != "" && v != "dev" {
		t.Skip("UBOX-S4: EXPECTED_BUILD_ID is set to a custom value; run S3 instead")
	}

	_, body := getBuildinfoContract(t)

	if body.BuildID != "dev" {
		t.Errorf("UBOX-S4: want build_id=%q when BUILD_ID env absent, got %q", "dev", body.BuildID)
	}
}

// Scenario: UBOX-S5 go_version is always "go1.23"
func TestBuildinfoContract_S5_GoVersionIsGo123(t *testing.T) {
	_, body := getBuildinfoContract(t)

	if body.GoVersion != "go1.23" {
		t.Errorf("UBOX-S5: want go_version=%q, got %q", "go1.23", body.GoVersion)
	}
}

// Scenario: UBOX-S6 response Content-Type is application/json
func TestBuildinfoContract_S6_ContentTypeJSON(t *testing.T) {
	resp, _ := getBuildinfoContract(t)

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("UBOX-S6: want Content-Type containing application/json, got %q", ct)
	}
}

// Scenario: UBOX-S7 endpoint requires no authentication
func TestBuildinfoContract_S7_NoAuthRequired(t *testing.T) {
	url := managementAddr() + "/buildinfo"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	// Explicitly send no Authorization header.

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("UBOX-S7: GET /buildinfo failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("UBOX-S7: want 200 without auth, got %d", resp.StatusCode)
	}
}
