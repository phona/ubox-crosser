//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var httpAddr = getEnv("HTTP_SERVER_ADDR", "proxy-server:8080")

func versionURL() string {
	return "http://" + httpAddr + "/version"
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}

func TestVersionEndpoint_HappyPath(t *testing.T) {
	resp, err := httpClient().Get(versionURL())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "v3", result["version"])
}

func TestVersionEndpoint_ResponseBodyExact(t *testing.T) {
	resp, err := httpClient().Get(versionURL())
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &parsed))

	assert.Len(t, parsed, 1, "response body should contain exactly one field")
	assert.Contains(t, parsed, "version")
}

func TestVersionEndpoint_MethodNotAllowed_POST(t *testing.T) {
	resp, err := httpClient().Post(versionURL(), "application/json", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestVersionEndpoint_MethodNotAllowed_PUT(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, versionURL(), nil)
	require.NoError(t, err)

	resp, err := httpClient().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestVersionEndpoint_MethodNotAllowed_DELETE(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, versionURL(), nil)
	require.NoError(t, err)

	resp, err := httpClient().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestVersionEndpoint_MethodNotAllowed_PATCH(t *testing.T) {
	req, err := http.NewRequest(http.MethodPatch, versionURL(), nil)
	require.NoError(t, err)

	resp, err := httpClient().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestVersionEndpoint_ContentTypeJSON(t *testing.T) {
	resp, err := httpClient().Get(versionURL())
	require.NoError(t, err)
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	assert.Contains(t, ct, "application/json")
}

func TestVersionEndpoint_Idempotent(t *testing.T) {
	var bodies [3]string
	for i := range bodies {
		resp, err := httpClient().Get(versionURL())
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(t, err)
		bodies[i] = string(body)
	}
	assert.Equal(t, bodies[0], bodies[1])
	assert.Equal(t, bodies[1], bodies[2])
}
