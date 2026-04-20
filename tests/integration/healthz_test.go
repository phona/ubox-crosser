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

var healthAddr = getEnv("HEALTH_SERVER_ADDR", "proxy-server:8080")

type healthResponse struct {
	Status        string   `json:"status"`
	UptimeSeconds *float64 `json:"uptime_seconds"`
	Services      []string `json:"services"`
}

func healthURL(path string) string {
	return fmt.Sprintf("http://%s%s", healthAddr, path)
}

func doHealthRequest(t *testing.T, method, path string) (*http.Response, []byte) {
	t.Helper()
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(method, healthURL(path), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send %s %s: %v", method, path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	return resp, body
}

func TestHealthzReturnsOK(t *testing.T) {
	resp, _ := doHealthRequest(t, http.MethodGet, "/healthz")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestHealthzContentTypeJSON(t *testing.T) {
	resp, _ := doHealthRequest(t, http.MethodGet, "/healthz")
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}

func TestHealthzBodyStructure(t *testing.T) {
	_, body := doHealthRequest(t, http.MethodGet, "/healthz")

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	for _, field := range []string{"status", "uptime_seconds", "services"} {
		if _, ok := raw[field]; !ok {
			t.Errorf("missing required field %q in response", field)
		}
	}
}

func TestHealthzStatusFieldIsOK(t *testing.T) {
	_, body := doHealthRequest(t, http.MethodGet, "/healthz")

	var hr healthResponse
	if err := json.Unmarshal(body, &hr); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if hr.Status != "ok" {
		t.Fatalf("expected status %q, got %q", "ok", hr.Status)
	}
}

func TestHealthzUptimeNonNegative(t *testing.T) {
	_, body := doHealthRequest(t, http.MethodGet, "/healthz")

	var hr healthResponse
	if err := json.Unmarshal(body, &hr); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if hr.UptimeSeconds == nil {
		t.Fatal("uptime_seconds is null")
	}
	if *hr.UptimeSeconds < 0 {
		t.Fatalf("uptime_seconds is negative: %f", *hr.UptimeSeconds)
	}
}

func TestHealthzServicesNotNull(t *testing.T) {
	_, body := doHealthRequest(t, http.MethodGet, "/healthz")

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	servicesRaw, ok := raw["services"]
	if !ok {
		t.Fatal("services field missing")
	}
	if string(servicesRaw) == "null" {
		t.Fatal("services field is null, expected an array")
	}

	var services []string
	if err := json.Unmarshal(servicesRaw, &services); err != nil {
		t.Fatalf("services is not a string array: %v", err)
	}
}

func TestHealthzServicesContainsConfigured(t *testing.T) {
	_, body := doHealthRequest(t, http.MethodGet, "/healthz")

	var hr healthResponse
	if err := json.Unmarshal(body, &hr); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expected := serveName
	found := false
	for _, s := range hr.Services {
		if s == expected {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("services %v does not contain configured service %q", hr.Services, expected)
	}
}

func TestHealthzMethodNotAllowed(t *testing.T) {
	resp, _ := doHealthRequest(t, http.MethodPost, "/healthz")
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestHealthzNotFound(t *testing.T) {
	resp, _ := doHealthRequest(t, http.MethodGet, "/nonexistent")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestHealthzUptimeIncreases(t *testing.T) {
	_, body1 := doHealthRequest(t, http.MethodGet, "/healthz")

	var hr1 healthResponse
	if err := json.Unmarshal(body1, &hr1); err != nil {
		t.Fatalf("failed to unmarshal first response: %v", err)
	}
	if hr1.UptimeSeconds == nil {
		t.Fatal("first uptime_seconds is null")
	}

	time.Sleep(1 * time.Second)

	_, body2 := doHealthRequest(t, http.MethodGet, "/healthz")
	var hr2 healthResponse
	if err := json.Unmarshal(body2, &hr2); err != nil {
		t.Fatalf("failed to unmarshal second response: %v", err)
	}
	if hr2.UptimeSeconds == nil {
		t.Fatal("second uptime_seconds is null")
	}

	if *hr2.UptimeSeconds < *hr1.UptimeSeconds {
		t.Fatalf("uptime decreased: first=%f, second=%f", *hr1.UptimeSeconds, *hr2.UptimeSeconds)
	}
}
