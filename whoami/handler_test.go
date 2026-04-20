package whoami_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/phona/ubox-crosser/whoami"
)

func TestHandler_StatusAndContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	rec := httptest.NewRecorder()

	whoami.Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Fatalf("expected Content-Type text/plain; charset=utf-8, got %q", ct)
	}
}

func TestHandler_BodyMatchesHostname(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	rec := httptest.NewRecorder()

	whoami.Handler(rec, req)

	want, err := os.Hostname()
	if err != nil {
		want = "unknown"
	}
	got := rec.Body.String()
	if got != want {
		t.Fatalf("expected body %q, got %q", want, got)
	}
}

func TestHandler_NonEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	rec := httptest.NewRecorder()

	whoami.Handler(rec, req)

	body := strings.TrimSpace(rec.Body.String())
	if len(body) == 0 {
		t.Fatal("expected non-empty body")
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /whoami", whoami.Handler)
	return mux
}

func TestMux_MethodNotAllowed(t *testing.T) {
	mux := newMux()
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/whoami", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s /whoami, got %d", method, rec.Code)
			}
		})
	}
}
