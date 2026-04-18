package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phona/ubox-crosser/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testHandler() http.Handler {
	configs := map[string]config.ServerConfig{
		"test": {
			LoginPass: "secret",
			AuthPass:  "auth-secret",
			Address:   "0.0.0.0:7000",
			Config: config.Config{
				Key:      "enc-key",
				Method:   "chacha20",
				LogLevel: "debug",
				LogFile:  "/var/log/test.log",
			},
		},
	}
	store := NewConfigStore(configs)
	return NewHandler(store)
}

func TestHandler_GetConfig_MasksAllSensitiveFields(t *testing.T) {
	t.Parallel()
	handler := testHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var body map[string]map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	svc := body["test"]
	assert.Equal(t, "***", svc["key"])
	assert.Equal(t, "***", svc["login_password"])
	assert.Equal(t, "***", svc["auth_password"])
	assert.Equal(t, "chacha20", svc["method"])
	assert.Equal(t, "debug", svc["log_level"])
}

func TestHandler_PutConfig_UpdatesAndReturnsFieldList(t *testing.T) {
	t.Parallel()
	handler := testHandler()
	body := `{"test": {"log_level": "warn"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string][]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Contains(t, resp["updated"], "test.log_level")
}

func TestHandler_PutConfig_ImmutableFieldReturns400(t *testing.T) {
	t.Parallel()
	handler := testHandler()
	body := `{"test": {"method": "aes-256-cfb"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_PutConfig_InvalidJSONReturns400(t *testing.T) {
	t.Parallel()
	handler := testHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_PutConfig_UnknownServiceReturns404(t *testing.T) {
	t.Parallel()
	handler := testHandler()
	body := `{"nope": {"log_level": "info"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	t.Parallel()
	methods := []string{http.MethodPost, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			handler := testHandler()
			req := httptest.NewRequest(method, "/api/config", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}
