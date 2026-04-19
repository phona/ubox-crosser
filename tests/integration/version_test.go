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

var managementAddr = getEnv("MANAGEMENT_ADDR", "proxy-server:8080")

// ---------------------------------------------------------------------------
// P0 — Happy Path
// ---------------------------------------------------------------------------

func TestVersionEndpoint_GET_ReturnsVersionJSON(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get("http://" + managementAddr + "/version")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))

	assert.Equal(t, "v2", result["version"])
}

func TestVersionEndpoint_GET_ResponseHasVersionFieldOnly(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get("http://" + managementAddr + "/version")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &result))

	assert.Len(t, result, 1, "response should contain exactly one field")
	assert.Contains(t, result, "version")
}

// ---------------------------------------------------------------------------
// P1 — Error Paths
// ---------------------------------------------------------------------------

func TestVersionEndpoint_POST_MethodNotAllowed(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Post("http://"+managementAddr+"/version", "application/json", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestVersionEndpoint_PUT_MethodNotAllowed(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodPut, "http://"+managementAddr+"/version", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestVersionEndpoint_DELETE_MethodNotAllowed(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodDelete, "http://"+managementAddr+"/version", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestVersionEndpoint_UnknownPath_NotFound(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get("http://" + managementAddr + "/nonexistent")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// ---------------------------------------------------------------------------
// P2 — Edge Cases
// ---------------------------------------------------------------------------

func TestVersionEndpoint_ConcurrentRequests(t *testing.T) {
	const concurrency = 10
	client := &http.Client{Timeout: 5 * time.Second}

	type result struct {
		status  int
		version string
		err     error
	}

	results := make(chan result, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			resp, err := client.Get("http://" + managementAddr + "/version")
			if err != nil {
				results <- result{err: err}
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				results <- result{status: resp.StatusCode, err: err}
				return
			}

			var parsed map[string]string
			if err := json.Unmarshal(body, &parsed); err != nil {
				results <- result{status: resp.StatusCode, err: err}
				return
			}

			results <- result{status: resp.StatusCode, version: parsed["version"]}
		}()
	}

	for i := 0; i < concurrency; i++ {
		r := <-results
		require.NoError(t, r.err, "concurrent request %d failed", i)
		assert.Equal(t, http.StatusOK, r.status)
		assert.Equal(t, "v2", r.version)
	}
}

func TestVersionEndpoint_ResponseIsValidJSON(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get("http://" + managementAddr + "/version")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.True(t, json.Valid(body), "response body should be valid JSON: %s", string(body))
}

func TestVersionEndpoint_NoTrailingSlashRedirect(t *testing.T) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get("http://" + managementAddr + "/version/")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.NotEqual(t, http.StatusMovedPermanently, resp.StatusCode,
		"should not redirect /version/ with 301")
}
