//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

var healthAddr = getHealthAddr()

func getHealthAddr() string {
	if v := os.Getenv("HEALTH_ADDR"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

type healthResponse struct {
	Status string `json:"status"`
}

func healthzURL() string {
	return healthAddr + "/healthz"
}

func TestHealthzContract_GET_200(t *testing.T) {
	resp, err := http.Get(healthzURL())
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var hr healthResponse
	if err := json.Unmarshal(body, &hr); err != nil {
		t.Fatalf("failed to unmarshal response: %v (body: %s)", err, body)
	}

	if hr.Status != "ok" {
		t.Errorf("expected status \"ok\", got %q", hr.Status)
	}
}

func TestHealthzContract_GET_NoExtraFields(t *testing.T) {
	resp, err := http.Get(healthzURL())
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("failed to unmarshal response as map: %v", err)
	}

	if len(raw) != 1 {
		t.Errorf("expected exactly 1 field, got %d: %v", len(raw), raw)
	}
	if _, ok := raw["status"]; !ok {
		t.Error("missing required field \"status\"")
	}
}

func TestHealthzContract_POST_405(t *testing.T) {
	resp, err := http.Post(healthzURL(), "application/json", nil)
	if err != nil {
		t.Fatalf("POST /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Errorf("expected Allow: GET header, got %q", allow)
	}
}

func TestHealthzContract_PUT_405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, healthzURL(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Errorf("expected Allow: GET header, got %q", allow)
	}
}

func TestHealthzContract_DELETE_405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, healthzURL(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Errorf("expected Allow: GET header, got %q", allow)
	}
}

func TestHealthzContract_UnknownPath_404(t *testing.T) {
	resp, err := http.Get(healthAddr + "/nonexistent")
	if err != nil {
		t.Fatalf("GET /nonexistent failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestHealthzContract_RootPath_404(t *testing.T) {
	resp, err := http.Get(healthAddr + "/")
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404 for root path, got %d", resp.StatusCode)
	}
}

func TestHealthzContract_HEAD_405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodHead, healthzURL(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HEAD /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Errorf("expected Allow: GET header, got %q", allow)
	}
}

func TestHealthzContract_PATCH_405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, healthzURL(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Errorf("expected Allow: GET header, got %q", allow)
	}
}
