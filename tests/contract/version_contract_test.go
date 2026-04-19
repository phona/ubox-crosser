//go:build contract

package contract

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

var httpBaseURL = getEnv("HTTP_BASE_URL", "http://localhost:8080")

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestGetVersion_Status200(t *testing.T) {
	resp, err := http.Get(httpBaseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetVersion_ContentTypeJSON(t *testing.T) {
	resp, err := http.Get(httpBaseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}

func TestGetVersion_ResponseBodySchema(t *testing.T) {
	resp, err := http.Get(httpBaseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v (body: %s)", err, body)
	}

	version, ok := result["version"]
	if !ok {
		t.Fatal("response JSON missing required field \"version\"")
	}

	versionStr, ok := version.(string)
	if !ok {
		t.Fatalf("field \"version\" must be a string, got %T", version)
	}

	if versionStr != "v3" {
		t.Fatalf("expected version \"v3\", got %q", versionStr)
	}
}

func TestGetVersion_NoExtraRequiredFields(t *testing.T) {
	resp, err := http.Get(httpBaseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if _, ok := result["version"]; !ok {
		t.Fatal("response JSON missing required field \"version\"")
	}
}

func TestGetVersion_ResponseBodyNotEmpty(t *testing.T) {
	resp, err := http.Get(httpBaseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("response body is empty")
	}
}

func TestPostVersion_MethodNotAllowed(t *testing.T) {
	resp, err := http.Post(httpBaseURL+"/version", "application/json", nil)
	if err != nil {
		t.Fatalf("POST /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 for POST, got %d", resp.StatusCode)
	}
}

func TestPutVersion_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, httpBaseURL+"/version", nil)
	if err != nil {
		t.Fatalf("failed to create PUT request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 for PUT, got %d", resp.StatusCode)
	}
}

func TestDeleteVersion_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, httpBaseURL+"/version", nil)
	if err != nil {
		t.Fatalf("failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 for DELETE, got %d", resp.StatusCode)
	}
}

func TestPatchVersion_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodPatch, httpBaseURL+"/version", nil)
	if err != nil {
		t.Fatalf("failed to create PATCH request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH /version failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 for PATCH, got %d", resp.StatusCode)
	}
}

func TestGetVersion_VersionFieldIsString(t *testing.T) {
	resp, err := http.Get(httpBaseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}

	versionRaw, ok := raw["version"]
	if !ok {
		t.Fatal("response JSON missing required field \"version\"")
	}

	var versionStr string
	if err := json.Unmarshal(versionRaw, &versionStr); err != nil {
		t.Fatalf("field \"version\" is not a JSON string: %s", versionRaw)
	}
}

func TestGetVersion_ResponseIsJSONObject(t *testing.T) {
	resp, err := http.Get(httpBaseURL + "/version")
	if err != nil {
		t.Fatalf("GET /version failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("response body is empty")
	}

	if body[0] != '{' {
		t.Fatalf("response is not a JSON object, starts with %q", string(body[0]))
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		t.Fatalf("response is not a valid JSON object: %v", err)
	}
}
