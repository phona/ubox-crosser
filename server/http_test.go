package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	w := httptest.NewRecorder()

	handlePing(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["message"] != "pong" {
		t.Errorf("expected message 'pong', got '%s'", body["message"])
	}
}

func TestHandlePingMethodPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/ping", nil)
	w := httptest.NewRecorder()

	handlePing(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
