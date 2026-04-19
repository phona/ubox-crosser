//go:build contract

package contract

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var (
	baseURL = getEnv("VERSION_HTTP_ADDR", "http://localhost:8080")
)

type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
}

func newClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}

func TestGetVersion_StatusOK(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetVersion_ContentType(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestGetVersion_ResponseSchema(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var info VersionInfo
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("failed to parse raw JSON: %v", err)
	}

	requiredFields := []string{"version", "git_commit", "build_time"}
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("required field %q missing from response", field)
		}
	}
}

func TestGetVersion_FieldTypes(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	for _, field := range []string{"version", "git_commit", "build_time"} {
		val, ok := raw[field]
		if !ok {
			t.Errorf("field %q missing", field)
			continue
		}
		if _, ok := val.(string); !ok {
			t.Errorf("field %q should be string, got %T", field, val)
		}
	}
}

func TestGetVersion_NoExtraFields(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	allowed := map[string]bool{"version": true, "git_commit": true, "build_time": true}
	for key := range raw {
		if !allowed[key] {
			t.Errorf("unexpected field %q in response", key)
		}
	}
}

func TestGetVersion_DefaultVersionIsDev(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var info VersionInfo
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if info.Version == "" {
		t.Error("version field should not be empty; expected at least \"dev\"")
	}
}

func TestGetVersion_BuildTimeFormat(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var info VersionInfo
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if info.BuildTime != "" {
		if _, err := time.Parse(time.RFC3339, info.BuildTime); err != nil {
			t.Errorf("build_time %q is not valid RFC 3339: %v", info.BuildTime, err)
		}
	}
}

func TestGetVersion_MethodNotAllowed(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	client := newClient()

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, baseURL+"/version", nil)
			if err != nil {
				t.Fatalf("failed to create %s request: %v", method, err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("%s /version failed: %v", method, err)
			}
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				t.Errorf("%s /version should not return 200", method)
			}
		})
	}
}

func TestGetVersion_UnknownPathReturns404(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/nonexistent")
	if err != nil {
		t.Fatalf("GET /nonexistent failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for unknown path, got %d", resp.StatusCode)
	}
}

func TestGetVersion_ResponseBodyNotEmpty(t *testing.T) {
	resp, err := newClient().Get(baseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	if len(body) == 0 {
		t.Error("response body should not be empty")
	}
}
