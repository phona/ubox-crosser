package ping_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phona/ubox-crosser/ping"
)

func TestHandler_StatusAndContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	ping.Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Fatalf("expected Content-Type text/plain; charset=utf-8, got %q", ct)
	}
}

func TestHandler_Body(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	ping.Handler(rec, req)

	if rec.Body.String() != "pong" {
		t.Fatalf("expected body %q, got %q", "pong", rec.Body.String())
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", ping.Handler)
	return mux
}

func TestMux_MethodNotAllowed(t *testing.T) {
	mux := newMux()
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/ping", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s /ping, got %d", method, rec.Code)
			}
		})
	}
}
