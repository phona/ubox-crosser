package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_GET_ReturnsOK(t *testing.T) {
	handler := newHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
	body, _ := io.ReadAll(rec.Body)
	if string(body) != `{"status":"ok"}` {
		t.Fatalf("expected {\"status\":\"ok\"}, got %s", body)
	}
}

func TestHealthHandler_POST_Returns405(t *testing.T) {
	handler := newHealthHandler()
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHealthHandler_PUT_Returns405(t *testing.T) {
	handler := newHealthHandler()
	req := httptest.NewRequest(http.MethodPut, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHealthHandler_DELETE_Returns405(t *testing.T) {
	handler := newHealthHandler()
	req := httptest.NewRequest(http.MethodDelete, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHealthHandler_UnknownPath_Returns404(t *testing.T) {
	handler := newHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHealthHandler_RootPath_Returns404(t *testing.T) {
	handler := newHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHealthHandler_HEAD_Returns405(t *testing.T) {
	handler := newHealthHandler()
	req := httptest.NewRequest(http.MethodHead, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
