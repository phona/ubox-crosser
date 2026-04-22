// +build integration

package contract

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type VersionResponse struct {
	Version string `json:"version"`
	Module  string `json:"module"`
	GoOS    string `json:"go_os"`
	GoArch  string `json:"go_arch"`
	GitSha  string `json:"git_sha"`
}

func getManagementServerAddrForVersion() string {
	// In contract test environment, management server should be accessible
	// Default to localhost:8888 if not set
	if addr := os.Getenv("MANAGEMENT_SERVER_ADDR"); addr != "" {
		return addr
	}
	return "localhost:8888"
}

func TestVersionEndpoint_GET_Returns200(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK")
}

func TestVersionEndpoint_ResponseIsValidJSON(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	var version VersionResponse
	err = json.NewDecoder(resp.Body).Decode(&version)
	require.NoError(t, err, "response should be valid JSON")
	assert.NotEmpty(t, version.Module, "module field should not be empty")
}

func TestVersionEndpoint_IncludesAllRequiredFields(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	var version VersionResponse
	err = json.NewDecoder(resp.Body).Decode(&version)
	require.NoError(t, err, "response should be valid JSON")

	assert.NotEmpty(t, version.Version, "version field should not be empty")
	assert.NotEmpty(t, version.Module, "module field should not be empty")
	assert.NotEmpty(t, version.GoOS, "go_os field should not be empty")
	assert.NotEmpty(t, version.GoArch, "go_arch field should not be empty")
	assert.NotEmpty(t, version.GitSha, "git_sha field should not be empty")
}

func TestVersionEndpoint_GitShaIsValidFormat(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	var version VersionResponse
	err = json.NewDecoder(resp.Body).Decode(&version)
	require.NoError(t, err, "response should be valid JSON")

	// Git SHA should be hexadecimal and at least 7 characters (short form)
	sha := version.GitSha
	assert.Greater(t, len(sha), 6, "git_sha should be at least 7 characters")
	// Verify it's valid hexadecimal
	matched, err := regexp.MatchString("^[0-9a-f]{7,40}$", sha)
	require.NoError(t, err, "regex validation failed")
	assert.True(t, matched, "git_sha should be valid hexadecimal (7-40 characters)")
}

func TestVersionEndpoint_GitShaConsistentAcrossRequests(t *testing.T) {
	addr := getManagementServerAddrForVersion()

	// First call
	resp1, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "first request failed")
	defer resp1.Body.Close()

	var version1 VersionResponse
	err = json.NewDecoder(resp1.Body).Decode(&version1)
	require.NoError(t, err, "first response should be valid JSON")
	sha1 := version1.GitSha

	// Second call
	resp2, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "second request failed")
	defer resp2.Body.Close()

	var version2 VersionResponse
	err = json.NewDecoder(resp2.Body).Decode(&version2)
	require.NoError(t, err, "second response should be valid JSON")
	sha2 := version2.GitSha

	assert.Equal(t, sha1, sha2, "git_sha should be identical across requests")
}

func TestVersionEndpoint_ModuleIsCorrect(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	var version VersionResponse
	err = json.NewDecoder(resp.Body).Decode(&version)
	require.NoError(t, err, "response should be valid JSON")

	assert.Equal(t, "github.com/phona/ubox-crosser", version.Module, "module should be github.com/phona/ubox-crosser")
}

func TestVersionEndpoint_GoOSAndGoArchPresent(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	var version VersionResponse
	err = json.NewDecoder(resp.Body).Decode(&version)
	require.NoError(t, err, "response should be valid JSON")

	// Go OS should be one of common values
	validOS := map[string]bool{
		"linux": true, "darwin": true, "windows": true,
		"freebsd": true, "openbsd": true, "netbsd": true,
	}
	assert.True(t, validOS[version.GoOS], fmt.Sprintf("go_os should be a valid OS, got %s", version.GoOS))

	// Go Arch should be one of common values
	validArch := map[string]bool{
		"amd64": true, "arm64": true, "386": true,
		"arm": true, "ppc64": true, "ppc64le": true,
	}
	assert.True(t, validArch[version.GoArch], fmt.Sprintf("go_arch should be a valid architecture, got %s", version.GoArch))
}

func TestVersionEndpoint_ContentTypeJSON(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "Content-Type should be application/json")
}

func TestVersionEndpoint_POST_Returns405(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	resp, err := http.Post(fmt.Sprintf("http://%s/version", addr), "application/json", nil)
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "POST should return 405 Method Not Allowed")
}

func TestVersionEndpoint_PUT_Returns405(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/version", addr), nil)
	require.NoError(t, err, "failed to create request")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "PUT should return 405 Method Not Allowed")
}

func TestVersionEndpoint_DELETE_Returns405(t *testing.T) {
	addr := getManagementServerAddrForVersion()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/version", addr), nil)
	require.NoError(t, err, "failed to create request")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "DELETE should return 405 Method Not Allowed")
}

func TestVersionEndpoint_ResponseTimeAcceptable(t *testing.T) {
	addr := getManagementServerAddrForVersion()

	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
	elapsed := time.Since(start)
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	// Read and discard body
	io.ReadAll(resp.Body)

	// Should complete in under 100ms
	assert.Less(t, elapsed, 100*time.Millisecond, "response time should be under 100ms")
}

func TestVersionEndpoint_AllFieldsPresistentAcrossRequests(t *testing.T) {
	addr := getManagementServerAddrForVersion()

	for i := 0; i < 3; i++ {
		resp, err := http.Get(fmt.Sprintf("http://%s/version", addr))
		require.NoError(t, err, "request %d failed", i)
		defer resp.Body.Close()

		var version VersionResponse
		err = json.NewDecoder(resp.Body).Decode(&version)
		require.NoError(t, err, "response %d should be valid JSON", i)

		assert.NotEmpty(t, version.Version, "version field should not be empty (request %d)", i)
		assert.NotEmpty(t, version.Module, "module field should not be empty (request %d)", i)
		assert.NotEmpty(t, version.GoOS, "go_os field should not be empty (request %d)", i)
		assert.NotEmpty(t, version.GoArch, "go_arch field should not be empty (request %d)", i)
		assert.NotEmpty(t, version.GitSha, "git_sha field should not be empty (request %d)", i)
	}
}
