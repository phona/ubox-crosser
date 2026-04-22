package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
