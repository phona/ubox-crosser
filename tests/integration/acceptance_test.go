//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// --- Happy Path ---

func TestAcceptance_EchoHappyPath_SimpleMessage(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made with msg=hello
	u := echoBaseURL + "/api/echo?msg=hello"
	resp, err := httpClient.Get(u)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// And Content-Type is application/json
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	// And the body contains {"message": "hello"}
	body, _ := io.ReadAll(resp.Body)
	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Message != "hello" {
		t.Errorf("expected message %q, got %q", "hello", result.Message)
	}
}

func TestAcceptance_EchoHappyPath_EmptyMessage(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made with msg= (empty value)
	u := echoBaseURL + "/api/echo?msg="
	resp, err := httpClient.Get(u)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// And the message field is an empty string
	body, _ := io.ReadAll(resp.Body)
	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Message != "" {
		t.Errorf("expected empty message, got %q", result.Message)
	}
}

// --- Error Cases ---

func TestAcceptance_EchoError_MissingMsgParam(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made without the msg parameter
	resp, err := httpClient.Get(echoBaseURL + "/api/echo")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 400
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	// And the body contains {"error": "msg parameter is required"}
	body, _ := io.ReadAll(resp.Body)
	var result ErrorResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Error != "msg parameter is required" {
		t.Errorf("expected error %q, got %q", "msg parameter is required", result.Error)
	}
}

func TestAcceptance_EchoError_PostMethodNotAllowed(t *testing.T) {
	// Given a running echo API server
	// When a POST request is made to /api/echo
	resp, err := httpClient.Post(echoBaseURL+"/api/echo?msg=hello", "application/json", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 405
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}

	// And the response is JSON with a non-empty error field
	body, _ := io.ReadAll(resp.Body)
	var result ErrorResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Error == "" {
		t.Error("expected non-empty error field")
	}
}

func TestAcceptance_EchoError_PutMethodNotAllowed(t *testing.T) {
	// Given a running echo API server
	// When a PUT request is made to /api/echo
	req, _ := http.NewRequest(http.MethodPut, echoBaseURL+"/api/echo?msg=hello", nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 405
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestAcceptance_EchoError_DeleteMethodNotAllowed(t *testing.T) {
	// Given a running echo API server
	// When a DELETE request is made to /api/echo
	req, _ := http.NewRequest(http.MethodDelete, echoBaseURL+"/api/echo?msg=hello", nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 405
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestAcceptance_EchoError_PatchMethodNotAllowed(t *testing.T) {
	// Given a running echo API server
	// When a PATCH request is made to /api/echo
	req, _ := http.NewRequest(http.MethodPatch, echoBaseURL+"/api/echo?msg=hello", nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 405
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

// --- Edge Cases ---

func TestAcceptance_EchoEdge_UnicodeMessage(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made with a Unicode message
	u, _ := url.Parse(echoBaseURL + "/api/echo")
	q := u.Query()
	q.Set("msg", "你好世界🌍")
	u.RawQuery = q.Encode()

	resp, err := httpClient.Get(u.String())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// And the message is returned correctly with Unicode preserved
	body, _ := io.ReadAll(resp.Body)
	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Message != "你好世界🌍" {
		t.Errorf("expected Unicode message preserved, got %q", result.Message)
	}
}

func TestAcceptance_EchoEdge_SpecialCharacters(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made with special characters including spaces and symbols
	u, _ := url.Parse(echoBaseURL + "/api/echo")
	q := u.Query()
	q.Set("msg", "hello world & foo=bar <script>")
	u.RawQuery = q.Encode()

	resp, err := httpClient.Get(u.String())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// And the message is returned with special characters preserved (no HTML escaping)
	body, _ := io.ReadAll(resp.Body)
	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Message != "hello world & foo=bar <script>" {
		t.Errorf("expected special chars preserved, got %q", result.Message)
	}
}

func TestAcceptance_EchoEdge_DuplicateMsgParams(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made with duplicate msg parameters
	resp, err := httpClient.Get(echoBaseURL + "/api/echo?msg=first&msg=second")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// And the first value is used
	body, _ := io.ReadAll(resp.Body)
	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Message != "first" {
		t.Errorf("expected first value %q, got %q", "first", result.Message)
	}
}

func TestAcceptance_EchoEdge_LongMessage(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made with a long message (1000 chars)
	longMsg := strings.Repeat("a", 1000)
	u, _ := url.Parse(echoBaseURL + "/api/echo")
	q := u.Query()
	q.Set("msg", longMsg)
	u.RawQuery = q.Encode()

	resp, err := httpClient.Get(u.String())
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// And the full message is echoed back
	body, _ := io.ReadAll(resp.Body)
	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Message != longMsg {
		t.Errorf("expected %d-char message echoed back, got %d chars", len(longMsg), len(result.Message))
	}
}

func TestAcceptance_EchoEdge_ExtraQueryParams(t *testing.T) {
	// Given a running echo API server
	// When a GET request includes msg plus extra unrelated query parameters
	resp, err := httpClient.Get(echoBaseURL + "/api/echo?msg=hello&foo=bar&baz=qux")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response status is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// And only the msg value is echoed back
	body, _ := io.ReadAll(resp.Body)
	var result EchoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Message != "hello" {
		t.Errorf("expected %q, got %q", "hello", result.Message)
	}
}

func TestAcceptance_EchoEdge_ResponseBodyIsExactJSON(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made with msg=test
	resp, err := httpClient.Get(echoBaseURL + "/api/echo?msg=test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the response contains only a "message" field (no extra fields)
	body, _ := io.ReadAll(resp.Body)
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(raw) != 1 {
		t.Errorf("expected exactly 1 field in response, got %d: %v", len(raw), raw)
	}
	if _, ok := raw["message"]; !ok {
		t.Error("expected 'message' field in response")
	}
}

func TestAcceptance_EchoEdge_ErrorResponseBodyIsExactJSON(t *testing.T) {
	// Given a running echo API server
	// When a GET request is made without msg parameter
	resp, err := httpClient.Get(echoBaseURL + "/api/echo")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the error response contains only an "error" field (no extra fields)
	body, _ := io.ReadAll(resp.Body)
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(raw) != 1 {
		t.Errorf("expected exactly 1 field in error response, got %d: %v", len(raw), raw)
	}
	if _, ok := raw["error"]; !ok {
		t.Error("expected 'error' field in error response")
	}
}

func TestAcceptance_EchoEdge_ContentTypeOnError(t *testing.T) {
	// Given a running echo API server
	// When a GET request triggers a 400 error
	resp, err := httpClient.Get(echoBaseURL + "/api/echo")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the Content-Type is still application/json
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json on error, got %q", ct)
	}
}

func TestAcceptance_EchoEdge_ContentTypeOn405(t *testing.T) {
	// Given a running echo API server
	// When a POST request triggers a 405 error
	resp, err := httpClient.Post(echoBaseURL+"/api/echo", "", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Then the Content-Type is still application/json
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json on 405, got %q", ct)
	}
}
