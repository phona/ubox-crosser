package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestBuildInfoHandler_StatusAndContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()

	BuildInfoHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want status 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("want Content-Type application/json, got %q", got)
	}
}

func TestBuildInfoHandler_JSONShape(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()

	BuildInfoHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	var info BuildInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if info.GitSHA == "" {
		t.Errorf("git_sha must be non-empty")
	}
	if info.BuildID == "" {
		t.Errorf("build_id must be non-empty")
	}
	if info.GoVersion != goVersionLiteral {
		t.Errorf("go_version want %q, got %q", goVersionLiteral, info.GoVersion)
	}
}

func TestBuildInfoHandler_GitSHAInjection(t *testing.T) {
	orig := GitSHA
	t.Cleanup(func() { GitSHA = orig })

	GitSHA = "abc1234"
	info := doBuildInfoRequest(t)
	if info.GitSHA != "abc1234" {
		t.Errorf("want git_sha=abc1234, got %q", info.GitSHA)
	}
}

func TestBuildInfoHandler_GitSHADefaultUnknown(t *testing.T) {
	orig := GitSHA
	t.Cleanup(func() { GitSHA = orig })

	GitSHA = "unknown"
	info := doBuildInfoRequest(t)
	if info.GitSHA != "unknown" {
		t.Errorf("want git_sha=unknown, got %q", info.GitSHA)
	}
}

func TestBuildInfoHandler_BuildIDFromEnv(t *testing.T) {
	t.Setenv("BUILD_ID", "ci-99999")
	info := doBuildInfoRequest(t)
	if info.BuildID != "ci-99999" {
		t.Errorf("want build_id=ci-99999, got %q", info.BuildID)
	}
}

func TestBuildInfoHandler_BuildIDDefaultDev(t *testing.T) {
	t.Setenv("BUILD_ID", "")
	info := doBuildInfoRequest(t)
	if info.BuildID != "dev" {
		t.Errorf("want build_id=dev (env empty), got %q", info.BuildID)
	}
}

func TestBuildInfoHandler_GoVersionLiteral(t *testing.T) {
	info := doBuildInfoRequest(t)
	if info.GoVersion != "go1.23" {
		t.Errorf("want go_version=go1.23, got %q", info.GoVersion)
	}
}

// TestGoVersionLiteralMatchesGoMod guards against toolchain drift: if
// someone bumps the `go` directive in go.mod without updating the hard-coded
// literal, this test fails CI and forces a conscious decision about the
// public contract.
func TestGoVersionLiteralMatchesGoMod(t *testing.T) {
	f, err := os.Open("../../go.mod")
	if err != nil {
		t.Fatalf("open go.mod: %v", err)
	}
	defer f.Close()

	var gomodGo string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "go ") {
			gomodGo = "go" + strings.TrimSpace(strings.TrimPrefix(line, "go "))
			break
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan go.mod: %v", err)
	}
	if gomodGo == "" {
		t.Fatal("no `go` directive found in go.mod")
	}
	if gomodGo != goVersionLiteral {
		t.Fatalf("go.mod `go` directive %q does not match /buildinfo literal %q — update goVersionLiteral and the contract spec together",
			gomodGo, goVersionLiteral)
	}
}

func doBuildInfoRequest(t *testing.T) BuildInfo {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	BuildInfoHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	var info BuildInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return info
}
