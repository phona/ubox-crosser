package echo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_WithMsg(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo?msg=hello", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("expected Content-Type text/plain; charset=utf-8, got %q", ct)
	}
	if rec.Body.String() != "hello" {
		t.Fatalf("expected body %q, got %q", "hello", rec.Body.String())
	}
}

func TestHandler_SpecialChars(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo?msg=hello%20world%21", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "hello world!" {
		t.Fatalf("expected body %q, got %q", "hello world!", rec.Body.String())
	}
}

func TestHandler_EmptyMsg(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo?msg=", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "" {
		t.Fatalf("expected empty body, got %q", rec.Body.String())
	}
}

func TestHandler_NoMsgParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	rec := httptest.NewRecorder()

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("expected Content-Type text/plain; charset=utf-8, got %q", ct)
	}
	if rec.Body.String() != "" {
		t.Fatalf("expected empty body, got %q", rec.Body.String())
	}
}
