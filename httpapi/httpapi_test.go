package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUptimeHandler_Returns200WithJSON(t *testing.T) {
	Init()
	handler := UptimeHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/uptime", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp UptimeResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.StartedAt)
	assert.GreaterOrEqual(t, resp.UptimeSeconds, 0)
}

func TestUptimeHandler_ResponseSchema(t *testing.T) {
	Init()
	handler := UptimeHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/uptime", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var m map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &m)
	require.NoError(t, err)

	assert.Contains(t, m, "started_at")
	assert.Contains(t, m, "uptime_seconds")
	assert.Len(t, m, 2)
}

func TestUptimeHandler_StartedAtIsRFC3339(t *testing.T) {
	Init()
	handler := UptimeHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/uptime", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var resp UptimeResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	_, err = time.Parse(time.RFC3339, resp.StartedAt)
	assert.NoError(t, err, "started_at should be a valid RFC 3339 timestamp")
}

func TestUptimeHandler_StartedAtStableAndUptimeIncreases(t *testing.T) {
	Init()
	handler := UptimeHandler()

	req1 := httptest.NewRequest(http.MethodGet, "/api/uptime", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	var resp1 UptimeResponse
	require.NoError(t, json.Unmarshal(w1.Body.Bytes(), &resp1))

	time.Sleep(1100 * time.Millisecond)

	req2 := httptest.NewRequest(http.MethodGet, "/api/uptime", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	var resp2 UptimeResponse
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp2))

	assert.Equal(t, resp1.StartedAt, resp2.StartedAt, "started_at should be stable across calls")
	assert.Greater(t, resp2.UptimeSeconds, resp1.UptimeSeconds, "uptime_seconds should increase over time")
}

func TestUnknownRoute_Returns404(t *testing.T) {
	Init()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/uptime", UptimeHandler())

	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
