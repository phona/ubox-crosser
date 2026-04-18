package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProvider struct {
	routes []RouteInfo
}

func (m *mockProvider) ListRoutes() []RouteInfo {
	return m.routes
}

func TestNewRoutesHandler_NilProvider_Returns200(t *testing.T) {
	handler := NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestNewRoutesHandler_NilProvider_ReturnsEmptyArray(t *testing.T) {
	handler := NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.JSONEq(t, "[]", rec.Body.String())
}

func TestNewRoutesHandler_ProviderReturnsNil_ReturnsEmptyArray(t *testing.T) {
	handler := NewRoutesHandler(&mockProvider{routes: nil})
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.JSONEq(t, "[]", rec.Body.String())
}

func TestNewRoutesHandler_ProviderReturnsRoutes(t *testing.T) {
	provider := &mockProvider{
		routes: []RouteInfo{
			{Name: "svc1", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: true},
			{Name: "svc2", ListenAddress: "127.0.0.1:7001", Method: "", Active: false},
		},
	}
	handler := NewRoutesHandler(provider)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var routes []RouteInfo
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	assert.Len(t, routes, 2)
	assert.Equal(t, "svc1", routes[0].Name)
	assert.Equal(t, true, routes[0].Active)
	assert.Equal(t, "svc2", routes[1].Name)
	assert.Equal(t, "", routes[1].Method)
	assert.Equal(t, false, routes[1].Active)
}

func TestNewRoutesHandler_MethodNotAllowed(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead}
	handler := NewRoutesHandler(nil)

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/routes", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
			assert.Empty(t, rec.Body.String())
		})
	}
}

func TestNewRoutesHandler_ContentTypeJSON(t *testing.T) {
	handler := NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestNewRoutesHandler_MethodNotAllowed_NoContentType(t *testing.T) {
	handler := NewRoutesHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Content-Type"))
}

func TestNewRoutesHandler_EmptyMethodField(t *testing.T) {
	provider := &mockProvider{
		routes: []RouteInfo{
			{Name: "svc", ListenAddress: "127.0.0.1:7000", Method: "", Active: false},
		},
	}
	handler := NewRoutesHandler(provider)
	req := httptest.NewRequest(http.MethodGet, "/api/routes", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var routes []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, "", routes[0]["method"])
}

func TestRouteInfo_JSONSerialization(t *testing.T) {
	route := RouteInfo{
		Name:          "test",
		ListenAddress: "127.0.0.1:8000",
		Method:        "aes-256-cfb",
		Active:        true,
	}

	data, err := json.Marshal(route)
	require.NoError(t, err)

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	require.NoError(t, err)

	assert.Equal(t, "test", m["name"])
	assert.Equal(t, "127.0.0.1:8000", m["listen_address"])
	assert.Equal(t, "aes-256-cfb", m["method"])
	assert.Equal(t, true, m["active"])
	assert.Len(t, m, 4, "RouteInfo should serialize to exactly 4 fields")
}
