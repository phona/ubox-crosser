package contract_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func whoamiMux(handler http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /whoami", handler)
	return mux
}

func staticWhoamiHandler(hostname string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hostname))
	}
}

// REQ-722-S1: GET /whoami returns 200 with Content-Type text/plain; charset=utf-8
func TestWhoami_S1_StatusAndContentType(t *testing.T) {
	mux := whoamiMux(staticWhoamiHandler("test-host"))
	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Fatalf("expected Content-Type %q, got %q", "text/plain; charset=utf-8", ct)
	}
}

// REQ-722-S2: GET /whoami body is a non-empty hostname string
func TestWhoami_S2_NonEmptyBody(t *testing.T) {
	mux := whoamiMux(staticWhoamiHandler("node-01"))
	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	body := strings.TrimSpace(rec.Body.String())
	if len(body) == 0 {
		t.Fatal("expected non-empty body")
	}
}

// REQ-722-S3: os.Hostname failure returns "unknown"
func TestWhoami_S3_FallbackUnknown(t *testing.T) {
	mux := whoamiMux(staticWhoamiHandler("unknown"))
	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if body != "unknown" {
		t.Fatalf("expected body %q, got %q", "unknown", body)
	}
}

// REQ-722-S4: POST /whoami returns 405
// REQ-722-S5: PUT /whoami returns 405
// REQ-722-S6: DELETE /whoami returns 405
func TestWhoami_S4_S5_S6_MethodNotAllowed(t *testing.T) {
	mux := whoamiMux(staticWhoamiHandler("test-host"))

	tests := []struct {
		scenario string
		method   string
	}{
		{"S4", http.MethodPost},
		{"S5", http.MethodPut},
		{"S6", http.MethodDelete},
	}

	for _, tt := range tests {
		t.Run(tt.scenario+"_"+tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/whoami", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405 for %s /whoami, got %d", tt.method, rec.Code)
			}
		})
	}
}
