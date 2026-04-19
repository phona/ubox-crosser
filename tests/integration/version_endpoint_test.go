//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

var httpAddr = getEnv("PROXY_HTTP_ADDR", "localhost:8080")

type versionResponse struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
}

func versionURL() string {
	return fmt.Sprintf("http://%s/version", httpAddr)
}

func TestVersionEndpoint_HappyPath(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var v versionResponse
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatalf("failed to parse JSON response: %v (body: %s)", err, string(body))
	}

	if v.Version == "" {
		t.Fatal("version field is empty, expected at least \"dev\"")
	}

	t.Logf("Version: %s, GitCommit: %s, BuildTime: %s", v.Version, v.GitCommit, v.BuildTime)
}

func TestVersionEndpoint_ResponseContainsAllFields(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	for _, field := range []string{"version", "git_commit", "build_time"} {
		if _, ok := raw[field]; !ok {
			t.Errorf("missing required field %q in response", field)
		}
	}
}

func TestVersionEndpoint_DefaultVersionIsDev(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var v versionResponse
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// In the test Docker environment the binary is built without -ldflags version injection,
	// so version should default to "dev" and git_commit/build_time should be empty.
	if v.Version != "dev" {
		t.Errorf("expected version \"dev\" (no ldflags injection), got %q", v.Version)
	}
	if v.GitCommit != "" {
		t.Errorf("expected empty git_commit (no ldflags injection), got %q", v.GitCommit)
	}
	if v.BuildTime != "" {
		t.Errorf("expected empty build_time (no ldflags injection), got %q", v.BuildTime)
	}
}

func TestVersionEndpoint_WrongMethod(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Post(versionURL(), "application/json", nil)
	if err != nil {
		t.Fatalf("POST /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Fatal("expected non-200 status for POST /version, but got 200")
	}
	t.Logf("POST /version returned status %d", resp.StatusCode)
}

func TestVersionEndpoint_UnknownRoute(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(fmt.Sprintf("http://%s/nonexistent", httpAddr))
	if err != nil {
		t.Fatalf("GET /nonexistent failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown route, got %d", resp.StatusCode)
	}
}

func TestVersionEndpoint_NoExtraFields(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	allowed := map[string]bool{"version": true, "git_commit": true, "build_time": true}
	for key := range raw {
		if !allowed[key] {
			t.Errorf("unexpected field %q in response", key)
		}
	}
}

func TestVersionEndpoint_FieldTypesAreStrings(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	for _, field := range []string{"version", "git_commit", "build_time"} {
		val, ok := raw[field]
		if !ok {
			t.Errorf("missing field %q", field)
			continue
		}
		if _, isStr := val.(string); !isStr {
			t.Errorf("field %q should be string, got %T", field, val)
		}
	}
}

func TestVersionEndpoint_ConcurrentRequests(t *testing.T) {
	const concurrency = 10
	client := &http.Client{Timeout: 5 * time.Second}
	errs := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			resp, err := client.Get(versionURL())
			if err != nil {
				errs <- fmt.Errorf("request failed: %w", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errs <- fmt.Errorf("expected 200, got %d", resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errs <- fmt.Errorf("failed to read body: %w", err)
				return
			}

			var v versionResponse
			if err := json.Unmarshal(body, &v); err != nil {
				errs <- fmt.Errorf("invalid JSON: %w", err)
				return
			}

			errs <- nil
		}()
	}

	for i := 0; i < concurrency; i++ {
		if err := <-errs; err != nil {
			t.Errorf("concurrent request %d: %v", i, err)
		}
	}
}
