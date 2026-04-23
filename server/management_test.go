package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildInfoHandler(t *testing.T) {
	tests := []struct {
		name          string
		gitSHA        string
		buildIDEnv    string
		wantGitSHA    string
		wantBuildID   string
		wantGoVersion string
	}{
		{
			name:          "default build_id when BUILD_ID unset",
			gitSHA:        "abc1234",
			buildIDEnv:    "",
			wantGitSHA:    "abc1234",
			wantBuildID:   "dev",
			wantGoVersion: "go1.23",
		},
		{
			name:          "uses BUILD_ID env when set",
			gitSHA:        "def5678",
			buildIDEnv:    "ci-42",
			wantGitSHA:    "def5678",
			wantBuildID:   "ci-42",
			wantGoVersion: "go1.23",
		},
		{
			name:          "empty git_sha is preserved",
			gitSHA:        "",
			buildIDEnv:    "",
			wantGitSHA:    "",
			wantBuildID:   "dev",
			wantGoVersion: "go1.23",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("BUILD_ID", tt.buildIDEnv)

			ms := NewManagementServer(tt.gitSHA)
			req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
			w := httptest.NewRecorder()
			ms.handleBuildinfo(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

			var resp buildInfoResponse
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, tt.wantGitSHA, resp.GitSHA)
			assert.Equal(t, tt.wantBuildID, resp.BuildID)
			assert.Equal(t, tt.wantGoVersion, resp.GoVersion)
		})
	}
}

func TestHealthzHandler(t *testing.T) {
	t.Parallel()
	ms := NewManagementServer("abc1234")
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	ms.handleHealthz(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var resp healthzResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "healthy", resp.Status)
	assert.GreaterOrEqual(t, resp.UptimeSeconds, int64(0))
	assert.Greater(t, resp.Timestamp, int64(0))
}

func TestHealthzUptimeReflectsStartTime(t *testing.T) {
	t.Parallel()
	ms := &ManagementServer{
		startTime: time.Now().Add(-10 * time.Second),
		gitSHA:    "abc1234",
		buildID:   "dev",
		mux:       http.NewServeMux(),
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	ms.handleHealthz(w, req)

	var resp healthzResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.GreaterOrEqual(t, resp.UptimeSeconds, int64(10),
		"uptime should be at least 10s when startTime was set 10s ago")
}

func TestBuildInfoRouteRegistered(t *testing.T) {
	t.Parallel()
	ms := NewManagementServer("sha123")
	srv := httptest.NewServer(ms.mux)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/buildinfo")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body buildInfoResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "sha123", body.GitSHA)
	assert.Equal(t, "go1.23", body.GoVersion)
}
