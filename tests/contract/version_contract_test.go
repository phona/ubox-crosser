//go:build contract

package contract

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func baseURL(t *testing.T) string {
	t.Helper()
	u := os.Getenv("CONTRACT_BASE_URL")
	if u == "" {
		t.Fatal("CONTRACT_BASE_URL not set")
	}
	return u
}

func TestVersion_S1_Returns200WithJSON(t *testing.T) {
	resp, err := http.Get(baseURL(t) + "/version")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}

	required := []string{"version", "commit", "build_time"}
	for _, key := range required {
		val, ok := body[key]
		if !ok {
			t.Errorf("missing required field %q", key)
			continue
		}
		if _, isStr := val.(string); !isStr {
			t.Errorf("field %q is %T, want string", key, val)
		}
	}

	if len(body) != 3 {
		t.Errorf("response has %d fields, want exactly 3 (additionalProperties: false)", len(body))
	}
}

func TestVersion_S5_NoAuthRequired(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, baseURL(t)+"/version", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unauthenticated request: status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestVersion_S6_PostReturns405(t *testing.T) {
	resp, err := http.Post(baseURL(t)+"/version", "application/json", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("POST status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Errorf("Allow header = %q, want %q", allow, "GET")
	}
}

func TestVersion_S6b_PutReturns405(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, baseURL(t)+"/version", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("PUT status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	allow := resp.Header.Get("Allow")
	if allow != "GET" {
		t.Errorf("Allow header = %q, want %q", allow, "GET")
	}
}

func TestVersion_S8_RootPathReturns404(t *testing.T) {
	resp, err := http.Get(baseURL(t) + "/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET / status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestVersion_S9_UnknownPathReturns404(t *testing.T) {
	resp, err := http.Get(baseURL(t) + "/metrics")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /metrics status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestVersion_S10_TrailingSlashReturns404(t *testing.T) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(baseURL(t) + "/version/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /version/ status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestVersion_S11_DefaultFieldsPresent(t *testing.T) {
	resp, err := http.Get(baseURL(t) + "/version")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}

	for _, key := range []string{"version", "commit", "build_time"} {
		val, ok := body[key]
		if !ok {
			t.Errorf("missing field %q in default build", key)
			continue
		}
		s, isStr := val.(string)
		if !isStr {
			t.Errorf("field %q is %T, want string", key, val)
			continue
		}
		if s == "" {
			t.Errorf("field %q is empty in default build, want non-empty default", key)
		}
	}
}
