package contract

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type VersionResponse struct {
	Commit string `json:"commit"`
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	commit := "abc123def456"
	_ = json.NewEncoder(w).Encode(VersionResponse{Commit: commit})
}

// REQ-969-S1: GET /api/version returns 200 with JSON {"commit":"<string>"}
func TestVersionReturnsOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(versionHandler))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/version")
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

	var ver VersionResponse
	if err := json.Unmarshal(body, &ver); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if ver.Commit == "" {
		t.Error("expected non-empty commit field")
	}
}

// REQ-969-S2: GET /api/version response has required JSON schema fields
func TestVersionResponseSchema(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(versionHandler))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/version")
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

	commitVal, ok := raw["commit"]
	if !ok {
		t.Fatal("missing required field \"commit\"")
	}

	if _, ok := commitVal.(string); !ok {
		t.Errorf("\"commit\" field should be a string, got %T", commitVal)
	}

	for key := range raw {
		if key != "commit" {
			t.Errorf("unexpected field %q in response", key)
		}
	}
}

// REQ-969-S3: Non-GET methods return 405 Method Not Allowed
func TestVersionRejectsNonGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(versionHandler))
	defer srv.Close()

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, srv.URL+"/api/version", nil)
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
