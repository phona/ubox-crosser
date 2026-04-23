package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/phona/ubox-crosser/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuildInfoHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		gitSHA  string
		buildID string
	}{
		{"with sha and build id", "abc1234", "build-42"},
		{"empty sha dev build", "", "dev"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			info := server.BuildInfo{
				GitSHA:    tt.gitSHA,
				BuildID:   tt.buildID,
				GoVersion: "go1.23",
			}
			handler := server.NewBuildInfoHandler(info)
			req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
			w := httptest.NewRecorder()
			handler(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var got server.BuildInfo
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
			assert.Equal(t, tt.gitSHA, got.GitSHA)
			assert.Equal(t, tt.buildID, got.BuildID)
			assert.Equal(t, "go1.23", got.GoVersion)
		})
	}
}

func TestNewHealthzHandler(t *testing.T) {
	t.Parallel()
	start := time.Now().Add(-5 * time.Second)
	handler := server.NewHealthzHandler(start)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got server.HealthzResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "healthy", got.Status)
	assert.GreaterOrEqual(t, got.UptimeSeconds, int64(4))
	assert.NotZero(t, got.Timestamp)
}

func TestNewBuildInfoHandlerGoVersionFixed(t *testing.T) {
	t.Parallel()
	info := server.BuildInfo{GitSHA: "deadbee", BuildID: "ci-99", GoVersion: "go1.23"}
	handler := server.NewBuildInfoHandler(info)
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	var got map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "go1.23", got["go_version"])
}
