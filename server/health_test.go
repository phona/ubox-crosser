package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHealthMux_GET_Healthz_200(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestNewHealthMux_POST_Healthz_405(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
	if allow := resp.Header.Get("Allow"); allow != "GET" {
		t.Errorf("expected Allow: GET, got %q", allow)
	}
}

func TestNewHealthMux_PUT_Healthz_405(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodPut, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
	if allow := resp.Header.Get("Allow"); allow != "GET" {
		t.Errorf("expected Allow: GET, got %q", allow)
	}
}

func TestNewHealthMux_DELETE_Healthz_405(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodDelete, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
	if allow := resp.Header.Get("Allow"); allow != "GET" {
		t.Errorf("expected Allow: GET, got %q", allow)
	}
}

func TestNewHealthMux_HEAD_Healthz_405(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodHead, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
	if allow := resp.Header.Get("Allow"); allow != "GET" {
		t.Errorf("expected Allow: GET, got %q", allow)
	}
}

func TestNewHealthMux_PATCH_Healthz_405(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodPatch, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
	if allow := resp.Header.Get("Allow"); allow != "GET" {
		t.Errorf("expected Allow: GET, got %q", allow)
	}
}

func TestNewHealthMux_UnknownPath_404(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Result().StatusCode)
	}
}

func TestNewHealthMux_RootPath_404(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Result().StatusCode)
	}
}

func TestNewHealthMux_HealthzTrailingSlash_404(t *testing.T) {
	mux := newHealthMux()
	req := httptest.NewRequest(http.MethodGet, "/healthz/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Result().StatusCode)
	}
}
