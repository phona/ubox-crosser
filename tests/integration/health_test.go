//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

var healthAddr = getEnv("HEALTH_ADDR", "localhost:8080")

func healthURL(path string) string {
	return fmt.Sprintf("http://%s%s", healthAddr, path)
}

// --- Happy Path ---

func TestHealthEndpoint_GET_ReturnsOK(t *testing.T) {
	// GIVEN a running proxy server with health endpoint
	// WHEN a client sends GET /health
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	// THEN the server responds with HTTP 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// AND Content-Type is application/json
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	// AND body is {"status":"ok"}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to parse JSON body: %v (body: %s)", err, string(body))
	}

	if result["status"] != "ok" {
		t.Fatalf("expected status \"ok\", got %q", result["status"])
	}
}

// --- Error Path: Method Not Allowed ---

func TestHealthEndpoint_POST_Returns405(t *testing.T) {
	// GIVEN a running proxy server with health endpoint
	// WHEN a client sends POST /health
	resp, err := http.Post(healthURL("/health"), "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("POST /health failed: %v", err)
	}
	defer resp.Body.Close()

	// THEN the server responds with 405 Method Not Allowed
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_PUT_Returns405(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, healthURL("/health"), strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("failed to create PUT request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_DELETE_Returns405(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, healthURL("/health"), nil)
	if err != nil {
		t.Fatalf("failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

// --- Error Path: Unknown Path ---

func TestHealthEndpoint_UnknownPath_Returns404(t *testing.T) {
	// WHEN a client sends GET /unknown
	resp, err := http.Get(healthURL("/unknown"))
	if err != nil {
		t.Fatalf("GET /unknown failed: %v", err)
	}
	defer resp.Body.Close()

	// THEN the server responds with 404 Not Found
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_RootPath_Returns404(t *testing.T) {
	resp, err := http.Get(healthURL("/"))
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

// --- Edge Cases ---

func TestHealthEndpoint_ResponseBodyExact(t *testing.T) {
	// Verify the exact JSON body matches the contract spec
	resp, err := http.Get(healthURL("/health"))
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("expected exactly 1 field in response, got %d: %v", len(parsed), parsed)
	}

	status, ok := parsed["status"]
	if !ok {
		t.Fatal("response missing 'status' field")
	}
	if status != "ok" {
		t.Fatalf("expected status 'ok', got %v", status)
	}
}

func TestHealthEndpoint_ConcurrentRequests(t *testing.T) {
	const concurrency = 20
	errs := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			resp, err := http.Get(healthURL("/health"))
			if err != nil {
				errs <- fmt.Errorf("request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errs <- fmt.Errorf("expected 200, got %d", resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errs <- fmt.Errorf("read body failed: %v", err)
				return
			}

			var result map[string]string
			if err := json.Unmarshal(body, &result); err != nil {
				errs <- fmt.Errorf("unmarshal failed: %v", err)
				return
			}

			if result["status"] != "ok" {
				errs <- fmt.Errorf("expected status 'ok', got %q", result["status"])
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

func TestHealthEndpoint_IsReachable(t *testing.T) {
	// Verify health endpoint is reachable via TCP (health server started with proxy)
	conn, err := net.DialTimeout("tcp", healthAddr, 5*time.Second)
	if err != nil {
		t.Fatalf("health endpoint not reachable at %s: %v", healthAddr, err)
	}
	conn.Close()
}

func TestHealthEndpoint_HEAD_Returns405(t *testing.T) {
	req, err := http.NewRequest(http.MethodHead, healthURL("/health"), nil)
	if err != nil {
		t.Fatalf("failed to create HEAD request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HEAD /health failed: %v", err)
	}
	defer resp.Body.Close()

	// Spec says only GET is allowed; HEAD should also be 405
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint_PATCH_Returns405(t *testing.T) {
	req, err := http.NewRequest(http.MethodPatch, healthURL("/health"), strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("failed to create PATCH request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}
