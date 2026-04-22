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

func TestManagementServer_VersionEndpoint(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/version")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var versionInfo VersionInfo
	err = json.NewDecoder(resp.Body).Decode(&versionInfo)
	require.NoError(t, err)

	assert.Equal(t, "github.com/phona/ubox-crosser", versionInfo.Module)
	assert.NotEmpty(t, versionInfo.Version)
}

func TestManagementServer_VersionEndpoint_MethodNotAllowed(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/version", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestManagementServer_HealthEndpoint(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var healthResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&healthResp)
	require.NoError(t, err)

	assert.Equal(t, "ok", healthResp["status"])
}

func TestManagementServer_Create(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	assert.NotNil(t, mgmt)
	assert.Equal(t, "127.0.0.1:8888", mgmt.address)
	assert.NotNil(t, mgmt.mux)
	assert.NotNil(t, mgmt.errs)
}

func TestManagementServer_HealthzEndpoint_GET_Returns200(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/healthz")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func TestManagementServer_HealthzEndpoint_ReturnsValidJSON(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	var healthz HealthzResponse
	err = json.NewDecoder(resp.Body).Decode(&healthz)
	require.NoError(t, err)

	assert.Equal(t, "ok", healthz.Status)
	assert.GreaterOrEqual(t, healthz.UptimeSeconds, int64(0))
}

func TestManagementServer_HealthzEndpoint_UptimeIncreases(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	resp1, err := http.Get(server.URL + "/healthz")
	require.NoError(t, err)
	defer resp1.Body.Close()

	var healthz1 HealthzResponse
	err = json.NewDecoder(resp1.Body).Decode(&healthz1)
	require.NoError(t, err)
	uptime1 := healthz1.UptimeSeconds

	time.Sleep(100 * time.Millisecond)

	resp2, err := http.Get(server.URL + "/healthz")
	require.NoError(t, err)
	defer resp2.Body.Close()

	var healthz2 HealthzResponse
	err = json.NewDecoder(resp2.Body).Decode(&healthz2)
	require.NoError(t, err)
	uptime2 := healthz2.UptimeSeconds

	assert.GreaterOrEqual(t, uptime2, uptime1)
}

func TestManagementServer_HealthzEndpoint_POST_Returns405(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/healthz", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestManagementServer_HealthzEndpoint_PUT_Returns405(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPut, server.URL+"/healthz", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestManagementServer_HealthzEndpoint_DELETE_Returns405(t *testing.T) {
	mgmt := NewManagementServer("127.0.0.1:8888")
	server := httptest.NewServer(mgmt.mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/healthz", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
