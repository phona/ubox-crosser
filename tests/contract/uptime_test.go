package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// uptimeResponse mirrors the contracted JSON schema for GET /uptime.
type uptimeResponse struct {
	UptimeSeconds int `json:"uptime_seconds"`
}

// newUptimeHandler builds a minimal handler that satisfies the /uptime contract.
// It uses the same algorithm described in design.md: time.Since(startTime) truncated to seconds.
// This stub lets the contract tests run before the real uptime package exists.
func newUptimeHandler(startTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		secs := int(time.Since(startTime).Seconds())
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(uptimeResponse{UptimeSeconds: secs})
	}
}

// REQ-736-S1: GET /uptime returns 200 with correct Content-Type and body shape.
func TestUptime_GET_Returns200WithJSON(t *testing.T) {
	h := newUptimeHandler(time.Now())
	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/uptime")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	var body uptimeResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.GreaterOrEqual(t, body.UptimeSeconds, 0)
}

// REQ-736-S2: uptime_seconds reflects real elapsed time.
func TestUptime_ReflectsElapsedTime(t *testing.T) {
	start := time.Now().Add(-2 * time.Second)
	h := newUptimeHandler(start)
	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/uptime")
	require.NoError(t, err)
	defer resp.Body.Close()

	var body uptimeResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.GreaterOrEqual(t, body.UptimeSeconds, 1)
}

// REQ-736-S3: Non-GET methods return 405.
func TestUptime_NonGET_Returns405(t *testing.T) {
	h := newUptimeHandler(time.Now())
	srv := httptest.NewServer(h)
	defer srv.Close()

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, srv.URL+"/uptime", nil)
			require.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		})
	}
}

// REQ-736-S4: Response body contains exactly one key "uptime_seconds".
func TestUptime_NoExtraFields(t *testing.T) {
	h := newUptimeHandler(time.Now())
	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/uptime")
	require.NoError(t, err)
	defer resp.Body.Close()

	var raw map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&raw))
	assert.Len(t, raw, 1, "response must have exactly one key")
	assert.Contains(t, raw, "uptime_seconds")
}
