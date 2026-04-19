//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

var httpBaseURL = getEnv("PROXY_HTTP_URL", "http://proxy-server:8080")

func versionURL() string {
	return httpBaseURL + "/version"
}

func TestVersionEndpoint_SuccessfulResponse(t *testing.T) {
	resp, err := httpClient().Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v\nbody: %s", err, body)
	}

	for _, field := range []string{"version", "commit", "buildTime"} {
		val, ok := result[field]
		if !ok {
			t.Errorf("missing required field %q", field)
			continue
		}
		str, ok := val.(string)
		if !ok {
			t.Errorf("field %q should be string, got %T", field, val)
			continue
		}
		if str == "" {
			t.Errorf("field %q should not be empty", field)
		}
	}
}

func TestVersionEndpoint_ContentType(t *testing.T) {
	resp, err := httpClient().Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestVersionEndpoint_NoExtraFields(t *testing.T) {
	resp, err := httpClient().Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	allowed := map[string]bool{"version": true, "commit": true, "buildTime": true}
	for key := range result {
		if !allowed[key] {
			t.Errorf("unexpected field %q in response", key)
		}
	}
}

func TestVersionEndpoint_DefaultValues(t *testing.T) {
	resp, err := httpClient().Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result struct {
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		BuildTime string `json:"buildTime"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result.Version == "" {
		t.Error("version must not be empty; expected at least 'dev'")
	}
	if result.Commit == "" {
		t.Error("commit must not be empty; expected at least 'unknown'")
	}
	if result.BuildTime == "" {
		t.Error("buildTime must not be empty; expected at least 'unknown'")
	}
}

func TestVersionEndpoint_MethodNotAllowed(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	client := httpClient()

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, versionURL(), nil)
			if err != nil {
				t.Fatalf("failed to create %s request: %v", method, err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("%s /version failed: %v", method, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				t.Errorf("%s /version should not return 200", method)
			}
		})
	}
}

func TestVersionEndpoint_NotFound(t *testing.T) {
	resp, err := httpClient().Get(httpBaseURL + "/nonexistent")
	if err != nil {
		t.Fatalf("GET /nonexistent failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for unknown path, got %d", resp.StatusCode)
	}
}

func TestVersionEndpoint_ResponseIsStrictJSON(t *testing.T) {
	resp, err := httpClient().Get(versionURL())
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	dec := json.NewDecoder(nil)
	_ = dec

	var first json.RawMessage
	if err := json.Unmarshal(body, &first); err != nil {
		t.Fatalf("body is not valid JSON: %v", err)
	}

	var obj map[string]string
	if err := json.Unmarshal(body, &obj); err != nil {
		t.Errorf("all fields should be string-typed: %v", err)
	}
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}
