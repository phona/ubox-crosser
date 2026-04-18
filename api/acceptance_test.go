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

// Acceptance tests validate end-to-end business logic using Given/When/Then scenarios
// from specs/config-read/spec.md and specs/config-update/spec.md.

func newAcceptanceHandler() (http.Handler, *api.ConfigStore) {
	configs := map[string]config.ServerConfig{
		"common": {
			LoginPass: "common-pass",
			AuthPass:  "common-auth",
			Address:   "0.0.0.0:7000",
			Config: config.Config{
				Key:      "common-key",
				Method:   "chacha20",
				LogLevel: "debug",
				LogFile:  "/var/log/crosser.log",
			},
		},
		"web": {
			LoginPass: "web-pass",
			AuthPass:  "web-auth",
			Address:   "0.0.0.0:7001",
			Config: config.Config{
				Key:      "web-key",
				Method:   "aes-256-cfb",
				LogLevel: "info",
			},
		},
	}
	store := api.NewConfigStore(configs)
	return api.NewHandler(store), store
}

// --- Scenario: Successful config retrieval (happy path) ---

func TestAcceptance_GetConfig_ReturnsAllServices(t *testing.T) {
	// Given a running server with "common" and "web" services configured
	handler, _ := newAcceptanceHandler()

	// When a GET request is sent to /api/config
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response contains both services
	assert.Equal(t, http.StatusOK, rec.Code)
	var body map[string]json.RawMessage
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Contains(t, body, "common")
	assert.Contains(t, body, "web")
}

// --- Scenario: Update log level successfully ---

func TestAcceptance_PutConfig_UpdateLogLevel(t *testing.T) {
	// Given a running server with common.log_level = "debug"
	handler, store := newAcceptanceHandler()

	// When a PUT request updates common.log_level to "info"
	body := `{"common": {"log_level": "info"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response is 200 with updated fields listed
	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string][]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["updated"], "common.log_level")

	// And the config store reflects the new value
	configs := store.GetAll()
	assert.Equal(t, "info", configs["common"].LogLevel)
}

// --- Scenario: Update log level for a specific service ---

func TestAcceptance_PutConfig_UpdateServiceLogLevel(t *testing.T) {
	// Given a running server with web.log_level = "info"
	handler, store := newAcceptanceHandler()

	// When a PUT request updates web.log_level to "warn"
	body := `{"web": {"log_level": "warn"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response is 200
	assert.Equal(t, http.StatusOK, rec.Code)

	// And the web service config reflects the new log level
	configs := store.GetAll()
	assert.Equal(t, "warn", configs["web"].LogLevel)
}

// --- Scenario: Reject update of immutable field ---

func TestAcceptance_PutConfig_RejectImmutableField(t *testing.T) {
	// Given a running server with common.key = "common-key"
	handler, store := newAcceptanceHandler()

	// When a PUT request tries to update the "key" field
	body := `{"common": {"key": "new-key"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response is 400 with error identifying "key" as immutable
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	errMsg, ok := resp["error"].(string)
	require.True(t, ok)
	assert.Contains(t, errMsg, "key")

	// And no configuration is modified
	configs := store.GetAll()
	assert.Equal(t, "common-key", configs["common"].Key)
}

// --- Scenario: Reject update with invalid JSON body ---

func TestAcceptance_PutConfig_RejectInvalidJSON(t *testing.T) {
	// Given a running server
	handler, _ := newAcceptanceHandler()

	// When a PUT request is sent with malformed JSON
	body := `{"common": {invalid`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response is 400 with error about invalid JSON
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp, "error")
}

// --- Scenario: Reject update for non-existent service ---

func TestAcceptance_PutConfig_RejectUnknownService(t *testing.T) {
	// Given a running server with only "common" and "web" services
	handler, _ := newAcceptanceHandler()

	// When a PUT request targets a non-existent service
	body := `{"unknown_service": {"log_level": "info"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response is 404 indicating service not found
	assert.Equal(t, http.StatusNotFound, rec.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp, "error")
}

// --- Scenario: Mixed mutable and immutable fields rejected entirely (atomicity) ---

func TestAcceptance_PutConfig_AtomicRejectMixedFields(t *testing.T) {
	// Given a running server with common.log_level = "debug"
	handler, store := newAcceptanceHandler()

	// When a PUT request contains both mutable (log_level) and immutable (key) fields
	body := `{"common": {"log_level": "info", "key": "new-key"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response is 400
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// And log_level is NOT updated (atomicity — entire request rejected)
	configs := store.GetAll()
	assert.Equal(t, "debug", configs["common"].LogLevel, "log_level SHALL NOT be updated when request contains immutable fields")
	assert.Equal(t, "common-key", configs["common"].Key, "key SHALL NOT be updated")
}

// --- Scenario: Sensitive fields masked in GET after PUT ---

func TestAcceptance_GetConfigAfterPut_SensitiveStillMasked(t *testing.T) {
	// Given a running server where log_level was updated
	handler, _ := newAcceptanceHandler()

	putBody := `{"common": {"log_level": "warn"}}`
	putReq := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(putBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	handler.ServeHTTP(putRec, putReq)
	require.Equal(t, http.StatusOK, putRec.Code)

	// When GET is called after the update
	getReq := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)

	// Then sensitive fields are still masked
	var body map[string]map[string]interface{}
	err := json.Unmarshal(getRec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "***", body["common"]["key"])
	assert.Equal(t, "***", body["common"]["login_password"])
	assert.Equal(t, "***", body["common"]["auth_password"])

	// And the updated field reflects the new value
	assert.Equal(t, "warn", body["common"]["log_level"])
}

// --- Scenario: Method not allowed ---

func TestAcceptance_PostConfig_MethodNotAllowed(t *testing.T) {
	// Given a running server
	handler, _ := newAcceptanceHandler()

	// When a POST request is sent to /api/config
	req := httptest.NewRequest(http.MethodPost, "/api/config", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response status code is 405
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

// --- Scenario: Multiple services updated in single request ---

func TestAcceptance_PutConfig_MultipleServicesUpdated(t *testing.T) {
	// Given a running server with "common" and "web" services
	handler, store := newAcceptanceHandler()

	// When a PUT request updates log_level for both services
	body := `{"common": {"log_level": "error"}, "web": {"log_level": "error"}}`
	req := httptest.NewRequest(http.MethodPut, "/api/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Then the response is 200
	assert.Equal(t, http.StatusOK, rec.Code)

	// And both services reflect the updated value
	configs := store.GetAll()
	assert.Equal(t, "error", configs["common"].LogLevel)
	assert.Equal(t, "error", configs["web"].LogLevel)
}
