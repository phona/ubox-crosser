package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phona/ubox-crosser/echo"
)

func newEchoMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /echo", echo.Handler)
	return mux
}

// REQ-924-S1: Echo endpoint returns msg parameter as plain text
func TestEcho_WithMsg(t *testing.T) {
	mux := newEchoMux()
	req := httptest.NewRequest(http.MethodGet, "/echo?msg=hello", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("[REQ-924-S1] expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Fatalf("[REQ-924-S1] expected Content-Type text/plain; charset=utf-8, got %q", ct)
	}
	if rec.Body.String() != "hello" {
		t.Fatalf("[REQ-924-S1] expected body %q, got %q", "hello", rec.Body.String())
	}
}

// REQ-924-S2: Echo endpoint returns empty body when msg is empty string
func TestEcho_EmptyMsg(t *testing.T) {
	mux := newEchoMux()
	req := httptest.NewRequest(http.MethodGet, "/echo?msg=", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("[REQ-924-S2] expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Fatalf("[REQ-924-S2] expected Content-Type text/plain; charset=utf-8, got %q", ct)
	}
	if rec.Body.String() != "" {
		t.Fatalf("[REQ-924-S2] expected empty body, got %q", rec.Body.String())
	}
}

// REQ-924-S3: Echo endpoint returns empty body when msg parameter is absent
func TestEcho_NoMsgParam(t *testing.T) {
	mux := newEchoMux()
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("[REQ-924-S3] expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Fatalf("[REQ-924-S3] expected Content-Type text/plain; charset=utf-8, got %q", ct)
	}
	if rec.Body.String() != "" {
		t.Fatalf("[REQ-924-S3] expected empty body, got %q", rec.Body.String())
	}
}

// REQ-924-S4: Echo endpoint rejects non-GET methods with 405
func TestEcho_MethodNotAllowed(t *testing.T) {
	mux := newEchoMux()
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/echo", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("[REQ-924-S4] expected 405 for %s /echo, got %d", method, rec.Code)
			}
		})
	}
}
