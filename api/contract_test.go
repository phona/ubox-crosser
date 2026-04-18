package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phona/ubox-crosser/api"
	"github.com/phona/ubox-crosser/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Contract tests validate the API schema, request/response formats, and status codes
// as defined in specs/config-read/spec.md and specs/config-update/spec.md.

func newTestHandler() http.Handler {
	configs := map[string]config.ServerConfig{
		"service1": {
			LoginPass: "secret-pass",
			AuthPass:  "auth-secret",
			Address:   "0.0.0.0:7000",
			Config: config.Config{
				Key:      "encryption-key",
				Method:   "chacha20",
				LogLevel: "debug",
			},
		},
	}
	store := api.NewConfigStore(configs)
	return api.NewHandler(store)
}

// --- GET /api/config contract ---

func TestGetConfig_StatusCode(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "GET /api/config SHALL return 200")
}

func TestGetConfig_ContentType(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	assert.Equal(t, "application/json", ct, "Content-Type SHALL be application/json")
}

func TestGetConfig_ResponseIsJSONObject(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var body map[string]json.RawMessage
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err, "response body SHALL be a valid JSON object")
	assert.Contains(t, body, "service1", "response SHALL contain service names as keys")
}

func TestGetConfig_SensitiveFieldsMasked(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var body map[string]map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)

	svc := body["service1"]
	assert.Equal(t, "***", svc["key"], "key SHALL be masked")
	assert.Equal(t, "***", svc["login_password"], "login_password SHALL be masked")
	assert.Equal(t, "***", svc["auth_password"], "auth_password SHALL be masked")
}

func TestGetConfig_NonSensitiveFieldsVisible(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var body map[string]map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)

	svc := body["service1"]
	assert.Equal(t, "chacha20", svc["method"], "method SHALL be visible")
	assert.Equal(t, "debug", svc["log_level"], "log_level SHALL be visible")
}

// --- Method not allowed contract ---

func TestPostConfig_MethodNotAllowed(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "POST SHALL return 405")
}

func TestDeleteConfig_MethodNotAllowed(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "DELETE SHALL return 405")
}

// --- PUT /api/config contract ---

func TestPutConfig_SuccessStatusCode(t *testing.T) {
	handler := newTestHandler()
	body := `{"service1": {"log_level": "info"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "PUT with valid mutable field SHALL return 200")
}

func TestPutConfig_SuccessResponseFormat(t *testing.T) {
	handler := newTestHandler()
	body := `{"service1": {"log_level": "info"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err, "response SHALL be valid JSON")
	assert.Contains(t, resp, "updated", "response SHALL contain 'updated' key")
}

func TestPutConfig_ImmutableFieldStatusCode(t *testing.T) {
	handler := newTestHandler()
	body := `{"service1": {"key": "new-key"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code, "PUT with immutable field SHALL return 400")
}

func TestPutConfig_InvalidJSONStatusCode(t *testing.T) {
	handler := newTestHandler()
	body := `{invalid json`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code, "PUT with invalid JSON SHALL return 400")
}

func TestPutConfig_UnknownServiceStatusCode(t *testing.T) {
	handler := newTestHandler()
	body := `{"unknown_service": {"log_level": "info"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code, "PUT with unknown service SHALL return 404")
}
