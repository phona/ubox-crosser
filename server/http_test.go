package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildInfoHandler_DefaultBuildID(t *testing.T) {
	t.Parallel()
	os.Unsetenv("BUILD_ID")

	handler := buildInfoHandler("abc1234")
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp buildInfoResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "abc1234", resp.GitSHA)
	assert.Equal(t, "dev", resp.BuildID)
	assert.Equal(t, "go1.23", resp.GoVersion)
}

func TestBuildInfoHandler_WithBuildID(t *testing.T) {
	t.Setenv("BUILD_ID", "build-42")

	handler := buildInfoHandler("abc1234")
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp buildInfoResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "abc1234", resp.GitSHA)
	assert.Equal(t, "build-42", resp.BuildID)
	assert.Equal(t, "go1.23", resp.GoVersion)
}

func TestBuildInfoHandler_ContentTypeJSON(t *testing.T) {
	t.Parallel()
	os.Unsetenv("BUILD_ID")

	handler := buildInfoHandler("abc1234")
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestBuildInfoHandler_EmptyGitSHA(t *testing.T) {
	t.Parallel()
	os.Unsetenv("BUILD_ID")

	handler := buildInfoHandler("")
	req := httptest.NewRequest(http.MethodGet, "/buildinfo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp buildInfoResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "", resp.GitSHA)
	assert.Equal(t, "dev", resp.BuildID)
	assert.Equal(t, "go1.23", resp.GoVersion)
}

func TestHealthzHandler_Returns200(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	healthzHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp healthzResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "healthy", resp.Status)
	assert.GreaterOrEqual(t, resp.UptimeSeconds, int64(0))
	assert.Greater(t, resp.Timestamp, int64(0))
}

func TestHealthzHandler_ContentTypeJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	healthzHandler(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
