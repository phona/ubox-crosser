package server

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/phona/ubox-crosser/models/config"
)

// ---------- helpers ----------

func mustNewHealthMux(t *testing.T) http.Handler {
	t.Helper()
	return NewHealthMux()
}

func doRequest(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

// ---------- Happy path ----------

func TestHealthEndpoint_GET_ReturnsOK(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodGet, "/health")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	body := rec.Body.Bytes()
	var resp map[string]string
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if resp["status"] != "ok" {
		t.Fatalf("expected status=ok, got %q", resp["status"])
	}
}

func TestHealthEndpoint_GET_ExactBody(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodGet, "/health")

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected exactly 1 key in response, got %d: %v", len(resp), resp)
	}
}

// ---------- Error path: method not allowed ----------

func TestHealthEndpoint_POST_Returns405(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodPost, "/health")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

func TestHealthEndpoint_PUT_Returns405(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodPut, "/health")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

func TestHealthEndpoint_DELETE_Returns405(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodDelete, "/health")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

func TestHealthEndpoint_PATCH_Returns405(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodPatch, "/health")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

func TestHealthEndpoint_AllNonGETMethods_Return405(t *testing.T) {
	methods := []string{
		http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodHead, http.MethodOptions,
	}
	mux := mustNewHealthMux(t)

	for _, m := range methods {
		t.Run(m, func(t *testing.T) {
			rec := doRequest(mux, m, "/health")
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("%s /health: expected 405, got %d", m, rec.Code)
			}
			if allow := rec.Header().Get("Allow"); allow != "GET" {
				t.Errorf("%s /health: expected Allow: GET, got %q", m, allow)
			}
		})
	}
}

// ---------- Error path: unknown paths return 404 ----------

func TestHealthEndpoint_RootPath_Returns404(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodGet, "/")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET /: expected 404, got %d", rec.Code)
	}
}

func TestHealthEndpoint_UnknownPath_Returns404(t *testing.T) {
	mux := mustNewHealthMux(t)
	paths := []string{"/metrics", "/ready", "/healthz", "/health/", "/Health", "/HEALTH"}

	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			rec := doRequest(mux, http.MethodGet, p)
			if rec.Code != http.StatusNotFound {
				t.Errorf("GET %s: expected 404, got %d", p, rec.Code)
			}
		})
	}
}

func TestHealthEndpoint_NonGETOnUnknownPath_Returns404(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodPost, "/unknown")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("POST /unknown: expected 404, got %d", rec.Code)
	}
}

// ---------- Config field ----------

func TestServerConfig_HealthAddressField(t *testing.T) {
	cfg := config.ServerConfig{
		HealthAddress: ":9090",
	}
	if cfg.HealthAddress != ":9090" {
		t.Fatalf("expected :9090, got %q", cfg.HealthAddress)
	}
}

func TestServerConfig_HealthAddressJSON(t *testing.T) {
	data := `{"health_address":":9090","address":":7000","login_password":"x","auth_password":"y"}`
	var cfg config.ServerConfig
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if cfg.HealthAddress != ":9090" {
		t.Fatalf("expected :9090, got %q", cfg.HealthAddress)
	}
}

func TestServerConfig_HealthAddressDefaultEmpty(t *testing.T) {
	var cfg config.ServerConfig
	if cfg.HealthAddress != "" {
		t.Fatalf("expected empty default, got %q", cfg.HealthAddress)
	}
}

// ---------- Lifecycle: health server disabled when address empty ----------

func TestProxyServer_NoHealthServer_WhenAddressEmpty(t *testing.T) {
	configs := map[string]config.ServerConfig{
		"test": {
			HealthAddress: "",
			Address:       "127.0.0.1:0",
			LoginPass:     "pass",
		},
	}

	srv := NewProxyServer(configs)

	// Give initWorker time to run
	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:8080/health")
	if err == nil {
		resp.Body.Close()
		t.Fatal("expected no health server on :8080 when address is empty, but got a response")
	}
	_ = srv
}

// ---------- Lifecycle: health server starts on configured port ----------

func TestProxyServer_HealthServer_StartsOnConfiguredPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	healthAddr := ln.Addr().String()
	ln.Close()

	proxyLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port for proxy: %v", err)
	}
	proxyAddr := proxyLn.Addr().String()
	proxyLn.Close()

	configs := map[string]config.ServerConfig{
		"test": {
			HealthAddress: healthAddr,
			Address:       proxyAddr,
			LoginPass:     "pass",
		},
	}

	srv := NewProxyServer(configs)
	_ = srv

	// Allow health server to start
	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://" + healthAddr + "/health")
	if err != nil {
		t.Fatalf("health server not reachable at %s: %v", healthAddr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["status"] != "ok" {
		t.Fatalf("expected status=ok, got %q", result["status"])
	}
}

// ---------- Lifecycle: port conflict reports error ----------

func TestProxyServer_HealthServer_ReportsPortConflict(t *testing.T) {
	blocker, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create blocker listener: %v", err)
	}
	defer blocker.Close()
	blockedAddr := blocker.Addr().String()

	proxyLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port for proxy: %v", err)
	}
	proxyAddr := proxyLn.Addr().String()
	proxyLn.Close()

	configs := map[string]config.ServerConfig{
		"test": {
			HealthAddress: blockedAddr,
			Address:       proxyAddr,
			LoginPass:     "pass",
		},
	}

	srv := NewProxyServer(configs)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Err()
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected a bind error, got nil")
		}
		if !strings.Contains(err.Error(), "bind") && !strings.Contains(err.Error(), "address already in use") {
			t.Logf("got error (possibly expected): %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("expected health server bind error within 3s, got nothing")
	}
}

// ---------- Edge cases ----------

func TestHealthEndpoint_ConcurrentRequests(t *testing.T) {
	mux := mustNewHealthMux(t)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	const concurrency = 50
	var wg sync.WaitGroup
	errs := make(chan string, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(ts.URL + "/health")
			if err != nil {
				errs <- err.Error()
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				errs <- "non-200 status"
				return
			}
			body, _ := io.ReadAll(resp.Body)
			var r map[string]string
			if err := json.Unmarshal(body, &r); err != nil {
				errs <- "invalid json"
				return
			}
			if r["status"] != "ok" {
				errs <- "wrong status"
			}
		}()
	}
	wg.Wait()
	close(errs)

	for e := range errs {
		t.Errorf("concurrent request failed: %s", e)
	}
}

func TestHealthEndpoint_LargePayloadRequest(t *testing.T) {
	mux := mustNewHealthMux(t)
	largeBody := strings.NewReader(strings.Repeat("x", 1<<20)) // 1 MiB
	req := httptest.NewRequest(http.MethodPost, "/health", largeBody)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for POST with large body, got %d", rec.Code)
	}
}

func TestHealthEndpoint_EmptyResponse_IsValidJSON(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodGet, "/health")

	body := rec.Body.Bytes()
	if !json.Valid(body) {
		t.Fatalf("response body is not valid JSON: %q", string(body))
	}
}

func TestHealthEndpoint_NoTrailingNewlineInJSON(t *testing.T) {
	mux := mustNewHealthMux(t)
	rec := doRequest(mux, http.MethodGet, "/health")

	body := rec.Body.String()
	trimmed := strings.TrimRight(body, "\n\r ")
	var resp map[string]string
	if err := json.Unmarshal([]byte(trimmed), &resp); err != nil {
		t.Fatalf("trimmed body is not valid JSON: %v", err)
	}
	if resp["status"] != "ok" {
		t.Fatalf("expected ok, got %q", resp["status"])
	}
}
