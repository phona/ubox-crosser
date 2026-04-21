package contract

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

// REQ-953-S1: GET /api/healthz returns 200 with JSON {"status":"ok"}
func TestHealthzReturnsOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(healthzHandler))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var health HealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if health.Status != "ok" {
		t.Errorf("expected status \"ok\", got %q", health.Status)
	}
}

// REQ-953-S2: GET /api/healthz response has required JSON schema fields
func TestHealthzResponseSchema(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(healthzHandler))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	statusVal, ok := raw["status"]
	if !ok {
		t.Fatal("missing required field \"status\"")
	}

	if _, ok := statusVal.(string); !ok {
		t.Errorf("\"status\" field should be a string, got %T", statusVal)
	}
}

// REQ-953-S3: Non-GET methods return 405 Method Not Allowed
func TestHealthzRejectsNonGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(healthzHandler))
	defer srv.Close()

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, srv.URL+"/api/healthz", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			_ = resp.Body.Close()

			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("expected status 405 for %s, got %d", method, resp.StatusCode)
			}
		})
	}
}
