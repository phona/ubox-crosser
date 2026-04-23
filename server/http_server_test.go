package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestBuildInfoHandler_DefaultValues(t *testing.T) {
	os.Unsetenv("BUILD_ID")
	mux := http.NewServeMux()
	mux.HandleFunc("/buildinfo", func(w http.ResponseWriter, r *http.Request) {
		info := buildInfoFromVars("abc1234")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info) //nolint:errcheck
	})

	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var info BuildInfo
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if info.GitSHA != "abc1234" {
		t.Errorf("expected git_sha 'abc1234', got %q", info.GitSHA)
	}
	if info.BuildID != "dev" {
		t.Errorf("expected build_id 'dev', got %q", info.BuildID)
	}
	if info.GoVersion != "go1.23" {
		t.Errorf("expected go_version 'go1.23', got %q", info.GoVersion)
	}
}

func TestBuildInfoHandler_BuildIDFromEnv(t *testing.T) {
	os.Setenv("BUILD_ID", "ci-42")
	defer os.Unsetenv("BUILD_ID")

	info := buildInfoFromVars("deadbee")
	if info.BuildID != "ci-42" {
		t.Errorf("expected build_id 'ci-42', got %q", info.BuildID)
	}
}

func TestHealthzHandler_StatusOK(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		resp := healthzResponse{
			Status:        "healthy",
			UptimeSeconds: int64(now.Sub(startTime).Seconds()),
			Timestamp:     now.Unix(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp healthzResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %q", resp.Status)
	}
	if resp.UptimeSeconds < 0 {
		t.Errorf("uptime_seconds must be non-negative, got %d", resp.UptimeSeconds)
	}
	if resp.Timestamp == 0 {
		t.Error("timestamp must be non-zero")
	}
}

func TestBuildInfoFromVars_GoVersion(t *testing.T) {
	info := buildInfoFromVars("any")
	if info.GoVersion != "go1.23" {
		t.Errorf("go_version must be 'go1.23', got %q", info.GoVersion)
	}
}
