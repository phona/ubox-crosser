package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoHandler_WithMessage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/echo?msg=hello", nil)
	rec := httptest.NewRecorder()

	EchoHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp echoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.Message != "hello" {
		t.Errorf("expected %q, got %q", "hello", resp.Message)
	}
}

func TestEchoHandler_EmptyMessage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/echo?msg=", nil)
	rec := httptest.NewRecorder()

	EchoHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp echoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.Message != "" {
		t.Errorf("expected empty message, got %q", resp.Message)
	}
}

func TestEchoHandler_MissingMsgParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/echo", nil)
	rec := httptest.NewRecorder()

	EchoHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.Error != "msg parameter is required" {
		t.Errorf("unexpected error: %q", resp.Error)
	}
}

func TestEchoHandler_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/echo?msg=hello", nil)
	rec := httptest.NewRecorder()

	EchoHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.Error != "method not allowed" {
		t.Errorf("unexpected error: %q", resp.Error)
	}
}
