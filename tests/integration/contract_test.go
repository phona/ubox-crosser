//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var echoBaseURL = getEnv("ECHO_API_BASE_URL", "http://localhost:8080")

var httpClient = &http.Client{Timeout: 10 * time.Second}

type EchoResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestEcho_200_ReturnsMessage(t *testing.T) {
	u, _ := url.Parse(echoBaseURL + "/api/echo")
	q := u.Query()
	q.Set("msg", "hello")
	u.RawQuery = q.Encode()

	resp, err := httpClient.Get(u.String())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if result.Message != "hello" {
		t.Errorf("expected message %q, got %q", "hello", result.Message)
	}
}

func TestEcho_200_EmptyMessage(t *testing.T) {
	u, _ := url.Parse(echoBaseURL + "/api/echo")
	q := u.Query()
	q.Set("msg", "")
	u.RawQuery = q.Encode()

	resp, err := httpClient.Get(u.String())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if result.Message != "" {
		t.Errorf("expected empty message, got %q", result.Message)
	}
}

func TestEcho_200_SchemaValidation(t *testing.T) {
	u, _ := url.Parse(echoBaseURL + "/api/echo")
	q := u.Query()
	q.Set("msg", "test")
	u.RawQuery = q.Encode()

	resp, err := httpClient.Get(u.String())
	if err != nil {
		t.Fatalf("request failed: %v", err)
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

	if _, ok := raw["message"]; !ok {
		t.Error("required field 'message' missing from response")
	}

	if _, ok := raw["message"].(string); !ok {
		t.Error("field 'message' must be a string")
	}
}

func TestEcho_400_MissingMsgParameter(t *testing.T) {
	resp, err := httpClient.Get(echoBaseURL + "/api/echo")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result ErrorResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if result.Error == "" {
		t.Error("expected non-empty error field")
	}
}

func TestEcho_400_ErrorSchemaValidation(t *testing.T) {
	resp, err := httpClient.Get(echoBaseURL + "/api/echo")
	if err != nil {
		t.Fatalf("request failed: %v", err)
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

	if _, ok := raw["error"]; !ok {
		t.Error("required field 'error' missing from error response")
	}

	if _, ok := raw["error"].(string); !ok {
		t.Error("field 'error' must be a string")
	}
}

func TestEcho_405_PostMethodNotAllowed(t *testing.T) {
	resp, err := httpClient.Post(echoBaseURL+"/api/echo", "application/json", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result ErrorResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if result.Error == "" {
		t.Error("expected non-empty error field")
	}
}

func TestEcho_405_PutMethodNotAllowed(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, echoBaseURL+"/api/echo", nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestEcho_405_DeleteMethodNotAllowed(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, echoBaseURL+"/api/echo", nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}
