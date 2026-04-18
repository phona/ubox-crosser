package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phona/ubox-crosser/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Contract tests validate the GET /api/routes API against the OpenAPI contract
// defined in contract.spec.yaml.
//
// These tests are expected to FAIL until the implementation is complete.

// --- GET /api/routes: status code ---

func TestGetRoutes_StatusCode(t *testing.T) {
	handler := api.NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "GET /api/routes SHALL return 200")
}

// --- GET /api/routes: content type ---

func TestGetRoutes_ContentType(t *testing.T) {
	handler := api.NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	assert.Equal(t, "application/json", ct, "Content-Type SHALL be application/json")
}

// --- GET /api/routes: response is JSON array ---

func TestGetRoutes_ResponseIsJSONArray(t *testing.T) {
	handler := api.NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var body []json.RawMessage
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err, "response body SHALL be a valid JSON array")
}

// --- GET /api/routes: empty config returns empty array ---

func TestGetRoutes_EmptyConfig_ReturnsEmptyArray(t *testing.T) {
	handler := api.NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var body []json.RawMessage
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Empty(t, body, "empty config SHALL return empty JSON array []")
}

// --- GET /api/routes: required fields present ---

func TestGetRoutes_RequiredFieldsPresent(t *testing.T) {
	provider := &stubRouteProvider{
		routes: []api.RouteInfo{
			{Name: "svc1", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: true},
		},
	}
	handler := api.NewRoutesHandler(provider)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var routes []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 1)

	route := routes[0]
	assert.Contains(t, route, "name", "route SHALL contain 'name' field")
	assert.Contains(t, route, "listen_address", "route SHALL contain 'listen_address' field")
	assert.Contains(t, route, "method", "route SHALL contain 'method' field")
	assert.Contains(t, route, "active", "route SHALL contain 'active' field")
}

// --- GET /api/routes: field types ---

func TestGetRoutes_FieldTypes(t *testing.T) {
	provider := &stubRouteProvider{
		routes: []api.RouteInfo{
			{Name: "svc1", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: true},
		},
	}
	handler := api.NewRoutesHandler(provider)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var routes []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 1)

	route := routes[0]
	_, nameOk := route["name"].(string)
	assert.True(t, nameOk, "name SHALL be a string")

	_, addrOk := route["listen_address"].(string)
	assert.True(t, addrOk, "listen_address SHALL be a string")

	_, methodOk := route["method"].(string)
	assert.True(t, methodOk, "method SHALL be a string")

	_, activeOk := route["active"].(bool)
	assert.True(t, activeOk, "active SHALL be a boolean")
}

// --- GET /api/routes: field values ---

func TestGetRoutes_FieldValues(t *testing.T) {
	provider := &stubRouteProvider{
		routes: []api.RouteInfo{
			{Name: "UBox_cytm", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: true},
			{Name: "UBox_ylxy", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: false},
		},
	}
	handler := api.NewRoutesHandler(provider)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var routes []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 2)

	assert.Equal(t, "UBox_cytm", routes[0]["name"])
	assert.Equal(t, "127.0.0.1:7000", routes[0]["listen_address"])
	assert.Equal(t, "chacha20", routes[0]["method"])
	assert.Equal(t, true, routes[0]["active"])

	assert.Equal(t, "UBox_ylxy", routes[1]["name"])
	assert.Equal(t, false, routes[1]["active"])
}

// --- GET /api/routes: no sensitive fields ---

func TestGetRoutes_NoSensitiveFields(t *testing.T) {
	provider := &stubRouteProvider{
		routes: []api.RouteInfo{
			{Name: "svc1", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: true},
		},
	}
	handler := api.NewRoutesHandler(provider)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var routes []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 1)

	route := routes[0]
	assert.NotContains(t, route, "key", "response SHALL NOT contain encryption key")
	assert.NotContains(t, route, "login_password", "response SHALL NOT contain login_password")
	assert.NotContains(t, route, "auth_password", "response SHALL NOT contain auth_password")
}

// --- Method not allowed ---

func TestPostRoutes_MethodNotAllowed(t *testing.T) {
	handler := api.NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "POST SHALL return 405")
}

func TestPutRoutes_MethodNotAllowed(t *testing.T) {
	handler := api.NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodPut, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "PUT SHALL return 405")
}

func TestDeleteRoutes_MethodNotAllowed(t *testing.T) {
	handler := api.NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodDelete, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "DELETE SHALL return 405")
}

// --- Active vs inactive routes ---

func TestGetRoutes_ActiveInactiveStatus(t *testing.T) {
	provider := &stubRouteProvider{
		routes: []api.RouteInfo{
			{Name: "active-svc", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: true},
			{Name: "inactive-svc", ListenAddress: "127.0.0.1:7001", Method: "aes-256-cfb", Active: false},
		},
	}
	handler := api.NewRoutesHandler(provider)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var routes []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 2)

	assert.Equal(t, true, routes[0]["active"], "connected service SHALL have active=true")
	assert.Equal(t, false, routes[1]["active"], "disconnected service SHALL have active=false")
}

// --- Only four fields in response (no extra fields) ---

func TestGetRoutes_OnlyContractFields(t *testing.T) {
	provider := &stubRouteProvider{
		routes: []api.RouteInfo{
			{Name: "svc1", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: true},
		},
	}
	handler := api.NewRoutesHandler(provider)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var routes []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 1)

	allowed := map[string]bool{"name": true, "listen_address": true, "method": true, "active": true}
	for key := range routes[0] {
		assert.True(t, allowed[key], "unexpected field %q in response — only contract fields allowed", key)
	}
}

// --- Stub for tests ---

type stubRouteProvider struct {
	routes []api.RouteInfo
}

func (s *stubRouteProvider) ListRoutes() []api.RouteInfo {
	return s.routes
}
