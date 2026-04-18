package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ---------- mock registry for contract tests ----------

type mockConnection struct {
	ID            string     `json:"id"`
	ServeName     string     `json:"serve_name"`
	RemoteAddr    string     `json:"remote_addr"`
	Status        string     `json:"status"`
	Type          string     `json:"type"`
	ConnectedAt   time.Time  `json:"connected_at"`
	LastHeartbeat *time.Time `json:"last_heartbeat,omitempty"`
}

// contractListResponse mirrors the expected JSON envelope for GET /api/connections
type contractListResponse struct {
	Connections []mockConnection `json:"connections"`
	Pagination  *struct {
		Page       *int `json:"page"`
		PageSize   *int `json:"page_size"`
		Total      *int `json:"total"`
		TotalPages *int `json:"total_pages"`
	} `json:"pagination"`
}

// contractSingleResponse mirrors the expected JSON envelope for GET /api/connections/{id}
type contractSingleResponse struct {
	Connection *mockConnection `json:"connection"`
}

// contractErrorResponse mirrors the expected error envelope
type contractErrorResponse struct {
	Error   *string `json:"error"`
	Message *string `json:"message"`
}

// ---------- helpers ----------

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	handler := NewHandler(nil) // registry injected per implementation
	mux := NewServeMux(handler)
	return httptest.NewServer(mux)
}

func mustGet(t *testing.T, url string) *http.Response {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s failed: %v", url, err)
	}
	return resp
}

// ---------- Contract Tests: GET /api/connections ----------

func TestContract_ListConnections_ContentType(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections")
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestContract_ListConnections_200ResponseSchema(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body contractListResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// connections must be a non-nil array (may be empty)
	if body.Connections == nil {
		t.Error("connections field must be a JSON array, got null")
	}

	// pagination must be present with all required fields
	if body.Pagination == nil {
		t.Fatal("pagination field is required")
	}
	if body.Pagination.Page == nil {
		t.Error("pagination.page is required")
	}
	if body.Pagination.PageSize == nil {
		t.Error("pagination.page_size is required")
	}
	if body.Pagination.Total == nil {
		t.Error("pagination.total is required")
	}
	if body.Pagination.TotalPages == nil {
		t.Error("pagination.total_pages is required")
	}
}

func TestContract_ListConnections_PaginationDefaults(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections")
	defer resp.Body.Close()

	var body contractListResponse
	json.NewDecoder(resp.Body).Decode(&body)

	if body.Pagination == nil {
		t.Fatal("pagination is required")
	}
	if body.Pagination.Page != nil && *body.Pagination.Page != 1 {
		t.Errorf("default page should be 1, got %d", *body.Pagination.Page)
	}
	if body.Pagination.PageSize != nil && *body.Pagination.PageSize != 20 {
		t.Errorf("default page_size should be 20, got %d", *body.Pagination.PageSize)
	}
}

func TestContract_ListConnections_ConnectionObjectSchema(t *testing.T) {
	// When connections exist, each object must have required fields
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections")
	defer resp.Body.Close()

	var raw map[string]json.RawMessage
	json.NewDecoder(resp.Body).Decode(&raw)

	connsRaw, ok := raw["connections"]
	if !ok {
		t.Fatal("missing connections field")
	}

	var conns []map[string]interface{}
	json.Unmarshal(connsRaw, &conns)

	requiredFields := []string{"id", "serve_name", "remote_addr", "status", "type", "connected_at"}

	for i, conn := range conns {
		for _, field := range requiredFields {
			if _, exists := conn[field]; !exists {
				t.Errorf("connection[%d] missing required field %q", i, field)
			}
		}

		// Validate enum values
		if status, ok := conn["status"].(string); ok {
			validStatuses := map[string]bool{"active": true, "idle": true, "terminated": true}
			if !validStatuses[status] {
				t.Errorf("connection[%d].status has invalid value %q", i, status)
			}
		}

		if connType, ok := conn["type"].(string); ok {
			validTypes := map[string]bool{"control": true, "worker": true}
			if !validTypes[connType] {
				t.Errorf("connection[%d].type has invalid value %q", i, connType)
			}
		}
	}
}

// ---------- Contract Tests: GET /api/connections/{id} ----------

func TestContract_GetConnection_ContentType(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections/any-id")
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestContract_GetConnection_200ResponseSchema(t *testing.T) {
	// This test will pass once there is an implementation that can return a real connection
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections/test-id")
	defer resp.Body.Close()

	// Accept either 200 (found) or 404 (not found) - both are valid contract responses
	if resp.StatusCode == http.StatusOK {
		var body contractSingleResponse
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode 200 response: %v", err)
		}
		if body.Connection == nil {
			t.Error("connection field is required in 200 response")
		}
	} else if resp.StatusCode == http.StatusNotFound {
		// 404 is acceptable for a nonexistent ID - validate error schema
		var errBody contractErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
			t.Fatalf("failed to decode 404 response: %v", err)
		}
		if errBody.Error == nil {
			t.Error("error field is required in error response")
		}
		if errBody.Message == nil {
			t.Error("message field is required in error response")
		}
	} else {
		t.Errorf("expected status 200 or 404, got %d", resp.StatusCode)
	}
}

func TestContract_GetConnection_404ResponseSchema(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections/nonexistent-id-12345")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for nonexistent connection, got %d", resp.StatusCode)
	}

	var errBody contractErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errBody.Error == nil {
		t.Error("error field is required in 404 response")
	}
	if errBody.Message == nil {
		t.Error("message field is required in 404 response")
	}
}

// ---------- Contract Tests: Error Responses ----------

func TestContract_ListConnections_400_InvalidStatus(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections?status=bogus")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid status, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("error response Content-Type should be application/json, got %q", ct)
	}

	var errBody contractErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errBody.Error == nil {
		t.Error("error field is required")
	}
	if errBody.Message == nil {
		t.Error("message field is required")
	}
}

func TestContract_ListConnections_400_InvalidType(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	resp := mustGet(t, ts.URL+"/api/connections?type=invalid")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid type, got %d", resp.StatusCode)
	}

	var errBody contractErrorResponse
	json.NewDecoder(resp.Body).Decode(&errBody)
	if errBody.Error == nil || errBody.Message == nil {
		t.Error("400 response must have error and message fields")
	}
}

func TestContract_ListConnections_400_InvalidPage(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	tests := []struct {
		name  string
		query string
	}{
		{"page=0", "?page=0"},
		{"page=-1", "?page=-1"},
		{"page=abc", "?page=abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := mustGet(t, ts.URL+"/api/connections"+tt.query)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400 for %s, got %d", tt.name, resp.StatusCode)
			}
		})
	}
}

func TestContract_ListConnections_400_InvalidPageSize(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	tests := []struct {
		name  string
		query string
	}{
		{"page_size=0", "?page_size=0"},
		{"page_size=-1", "?page_size=-1"},
		{"page_size=101", "?page_size=101"},
		{"page_size=abc", "?page_size=abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := mustGet(t, ts.URL+"/api/connections"+tt.query)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400 for %s, got %d", tt.name, resp.StatusCode)
			}
		})
	}
}

func TestContract_GetConnection_400_EmptyID(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// Trailing slash with no ID
	resp := mustGet(t, ts.URL+"/api/connections/")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty connection ID, got %d", resp.StatusCode)
	}
}

// ---------- Contract Tests: Method Not Allowed ----------

func TestContract_ListConnections_405_POST(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/api/connections", "application/json", nil)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for POST, got %d", resp.StatusCode)
	}
}

func TestContract_GetConnection_405_DELETE(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/connections/some-id", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for DELETE, got %d", resp.StatusCode)
	}
}
