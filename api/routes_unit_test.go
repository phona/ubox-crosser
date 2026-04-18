package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Unit tests for internal handler logic (package-internal access).

type mockProvider struct {
	routes []RouteInfo
}

func (m *mockProvider) ListRoutes() []RouteInfo {
	return m.routes
}

func TestHandler_NilProvider_Returns200EmptyArray(t *testing.T) {
	handler := NewRoutesHandler(nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	assert.Equal(t, http.StatusOK, rec.Code)

	var routes []RouteInfo
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	assert.Equal(t, []RouteInfo{}, routes)
}

func TestHandler_ProviderReturnsNilSlice_Returns200EmptyArray(t *testing.T) {
	handler := NewRoutesHandler(&mockProvider{routes: nil})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
}

func TestHandler_ProviderReturnsRoutes_SerializesCorrectly(t *testing.T) {
	provider := &mockProvider{
		routes: []RouteInfo{
			{Name: "svc1", ListenAddress: "127.0.0.1:7000", Method: "chacha20", Active: true},
			{Name: "svc2", ListenAddress: "127.0.0.1:7001", Method: "", Active: false},
		},
	}
	handler := NewRoutesHandler(provider)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	assert.Equal(t, http.StatusOK, rec.Code)

	var routes []RouteInfo
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 2)

	assert.Equal(t, "svc1", routes[0].Name)
	assert.Equal(t, "127.0.0.1:7000", routes[0].ListenAddress)
	assert.Equal(t, "chacha20", routes[0].Method)
	assert.True(t, routes[0].Active)

	assert.Equal(t, "svc2", routes[1].Name)
	assert.Equal(t, "", routes[1].Method)
	assert.False(t, routes[1].Active)
}

func TestHandler_MethodNotAllowed_NoBody(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	handler := NewRoutesHandler(nil)

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(method, "/api/routes", nil))

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
			assert.Empty(t, rec.Body.String())
		})
	}
}

func TestHandler_ContentTypeSetOnlyFor200(t *testing.T) {
	handler := NewRoutesHandler(nil)

	// GET should have Content-Type
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	// POST (405) should not have Content-Type set
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/routes", nil))
	assert.Empty(t, rec.Header().Get("Content-Type"))
}

func TestHandler_EmptyProviderRoutes_ReturnsEmptyArray(t *testing.T) {
	handler := NewRoutesHandler(&mockProvider{routes: []RouteInfo{}})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
}

func TestHandler_EmptyMethodField_SerializesAsEmptyString(t *testing.T) {
	provider := &mockProvider{
		routes: []RouteInfo{
			{Name: "plain", ListenAddress: "127.0.0.1:8000", Method: "", Active: true},
		},
	}
	handler := NewRoutesHandler(provider)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/routes", nil))

	var routes []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &routes)
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, "", routes[0]["method"])
}
