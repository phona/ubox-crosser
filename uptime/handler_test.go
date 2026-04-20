package uptime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GET_Returns200WithJSON(t *testing.T) {
	startTime = time.Now()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/uptime", nil)

	Handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var body response
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.GreaterOrEqual(t, body.UptimeSeconds, 0)
}

func TestHandler_ReflectsElapsedTime(t *testing.T) {
	startTime = time.Now().Add(-3 * time.Second)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/uptime", nil)

	Handler(rec, req)

	var body response
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.GreaterOrEqual(t, body.UptimeSeconds, 2)
}

func TestHandler_NonGET_Returns405(t *testing.T) {
	startTime = time.Now()
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/uptime", nil)

			Handler(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}

func TestHandler_NoExtraFields(t *testing.T) {
	startTime = time.Now()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/uptime", nil)

	Handler(rec, req)

	var raw map[string]interface{}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&raw))
	assert.Len(t, raw, 1, "response must have exactly one key")
	assert.Contains(t, raw, "uptime_seconds")
}
