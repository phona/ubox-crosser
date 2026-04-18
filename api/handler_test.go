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

func testHandler() (http.Handler, *ConfigStore) {
	configs := map[string]config.ServerConfig{
		"svc1": {
			LoginPass: "pass1",
			AuthPass:  "auth1",
			Address:   "0.0.0.0:7000",
			Config: config.Config{
				Key:      "key1",
				Method:   "chacha20",
				LogLevel: "debug",
				LogFile:  "/var/log/test.log",
			},
		},
	}
	store := NewConfigStore(configs)
	return NewHandler(store), store
}

func TestHandler_GET_MasksSensitiveFields(t *testing.T) {
	handler, _ := testHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var body map[string]map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "***", body["svc1"]["key"])
	assert.Equal(t, "***", body["svc1"]["login_password"])
	assert.Equal(t, "***", body["svc1"]["auth_password"])
}

func TestHandler_GET_ShowsNonSensitiveFields(t *testing.T) {
	handler, _ := testHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var body map[string]map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "chacha20", body["svc1"]["method"])
	assert.Equal(t, "debug", body["svc1"]["log_level"])
	assert.Equal(t, "0.0.0.0:7000", body["svc1"]["address"])
}

func TestHandler_PUT_ValidUpdate(t *testing.T) {
	handler, store := testHandler()
	body := `{"svc1": {"log_level": "error"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	configs := store.GetAll()
	assert.Equal(t, "error", configs["svc1"].LogLevel)
}

func TestHandler_PUT_ImmutableFieldRejected(t *testing.T) {
	handler, store := testHandler()
	body := `{"svc1": {"method": "aes"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	configs := store.GetAll()
	assert.Equal(t, "chacha20", configs["svc1"].Method)
}

func TestHandler_PUT_InvalidJSON(t *testing.T) {
	handler, _ := testHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_PUT_UnknownService(t *testing.T) {
	handler, _ := testHandler()
	body := `{"missing": {"log_level": "info"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	handler, _ := testHandler()
	for _, method := range []string{http.MethodPatch, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/config", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "method %s should return 405", method)
	}
}
