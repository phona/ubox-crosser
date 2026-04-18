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

func TestNewRoutesHandler_NilProvider_ReturnsEmptyArray(t *testing.T) {
	handler := NewRoutesHandler(nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
}

func TestNewRoutesHandler_ProviderReturnsNil_ReturnsEmptyArray(t *testing.T) {
	handler := NewRoutesHandler(&mockProvider{routes: nil})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
}

func TestNewRoutesHandler_SerializesRoutes(t *testing.T) {
	provider := &mockProvider{
		routes: []RouteInfo{
			{Name: "r1", ListenAddress: "0.0.0.0:8000", Method: "aes-256-cfb", Active: true},
			{Name: "r2", ListenAddress: "0.0.0.0:8001", Method: "", Active: false},
		},
	}
	handler := NewRoutesHandler(provider)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	var routes []RouteInfo
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 2)

	assert.Equal(t, "r1", routes[0].Name)
	assert.Equal(t, "0.0.0.0:8000", routes[0].ListenAddress)
	assert.Equal(t, "aes-256-cfb", routes[0].Method)
	assert.True(t, routes[0].Active)

	assert.Equal(t, "r2", routes[1].Name)
	assert.Equal(t, "", routes[1].Method)
	assert.False(t, routes[1].Active)
}

func TestNewRoutesHandler_RejectsNonGetMethods(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead}
	handler := NewRoutesHandler(nil)

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(method, "/api/routes", nil))
			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}

func TestNewRoutesHandler_ContentTypeIsJSON(t *testing.T) {
	handler := NewRoutesHandler(&mockProvider{routes: []RouteInfo{{Name: "x"}}})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestNewRoutesHandler_MethodNotAllowed_NoBody(t *testing.T) {
	handler := NewRoutesHandler(nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/routes", nil))

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	assert.Empty(t, rec.Body.String())
}
