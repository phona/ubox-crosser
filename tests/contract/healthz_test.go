// +build integration

package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HealthzResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

func getManagementServerAddr() string {
	// In contract test environment, management server should be accessible
	// Default to localhost:8888 if not set
	if addr := os.Getenv("MANAGEMENT_SERVER_ADDR"); addr != "" {
		return addr
	}
	return "localhost:8888"
}

func TestHealthzEndpoint_GET_Returns200(t *testing.T) {
	addr := getManagementServerAddr()
	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK")
}

func TestHealthzEndpoint_ResponseIsValidJSON(t *testing.T) {
	addr := getManagementServerAddr()
	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	var healthz HealthzResponse
	err = json.NewDecoder(resp.Body).Decode(&healthz)
	require.NoError(t, err, "response should be valid JSON")
	assert.Equal(t, "ok", healthz.Status, "status should be 'ok'")
}

func TestHealthzEndpoint_IncludesUptimeSeconds(t *testing.T) {
	addr := getManagementServerAddr()
	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	var healthz HealthzResponse
	err = json.NewDecoder(resp.Body).Decode(&healthz)
	require.NoError(t, err, "response should be valid JSON")

	assert.Greater(t, healthz.UptimeSeconds, int64(-1), "uptime_seconds should be non-negative")
}

func TestHealthzEndpoint_UptimeIncreases(t *testing.T) {
	addr := getManagementServerAddr()

	// First call
	resp1, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	require.NoError(t, err, "first request failed")
	defer resp1.Body.Close()

	var healthz1 HealthzResponse
	err = json.NewDecoder(resp1.Body).Decode(&healthz1)
	require.NoError(t, err, "first response should be valid JSON")
	uptime1 := healthz1.UptimeSeconds

	// Wait a second
	time.Sleep(1 * time.Second)

	// Second call
	resp2, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	require.NoError(t, err, "second request failed")
	defer resp2.Body.Close()

	var healthz2 HealthzResponse
	err = json.NewDecoder(resp2.Body).Decode(&healthz2)
	require.NoError(t, err, "second response should be valid JSON")
	uptime2 := healthz2.UptimeSeconds

	assert.GreaterOrEqual(t, uptime2, uptime1, "uptime should be monotonically increasing")
}

func TestHealthzEndpoint_ContentTypeJSON(t *testing.T) {
	addr := getManagementServerAddr()
	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "Content-Type should be application/json")
}

func TestHealthzEndpoint_POST_Returns405(t *testing.T) {
	addr := getManagementServerAddr()
	resp, err := http.Post(fmt.Sprintf("http://%s/healthz", addr), "application/json", bytes.NewBuffer([]byte{}))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "POST should return 405 Method Not Allowed")
}

func TestHealthzEndpoint_PUT_Returns405(t *testing.T) {
	addr := getManagementServerAddr()
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/healthz", addr), nil)
	require.NoError(t, err, "failed to create request")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "PUT should return 405 Method Not Allowed")
}

func TestHealthzEndpoint_DELETE_Returns405(t *testing.T) {
	addr := getManagementServerAddr()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/healthz", addr), nil)
	require.NoError(t, err, "failed to create request")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "DELETE should return 405 Method Not Allowed")
}

func TestHealthzEndpoint_ResponseTimeAcceptable(t *testing.T) {
	addr := getManagementServerAddr()

	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	elapsed := time.Since(start)
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	// Read and discard body
	io.ReadAll(resp.Body)

	// Should complete in under 100ms
	assert.Less(t, elapsed, 100*time.Millisecond, "response time should be under 100ms")
}

func TestHealthzEndpoint_StatusAlwaysOk(t *testing.T) {
	addr := getManagementServerAddr()

	for i := 0; i < 5; i++ {
		resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
		require.NoError(t, err, "request %d failed", i)
		defer resp.Body.Close()

		var healthz HealthzResponse
		err = json.NewDecoder(resp.Body).Decode(&healthz)
		require.NoError(t, err, "response %d should be valid JSON", i)

		assert.Equal(t, "ok", healthz.Status, "status should always be 'ok'")
	}
}
