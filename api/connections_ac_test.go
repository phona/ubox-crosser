package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =============================================================================
// Acceptance Tests for /api/connections
//
// Written in Given/When/Then style to validate business logic and end-to-end
// behavior. These tests verify the API does the right thing.
// =============================================================================

// ---------- helpers ----------

func setupACServer(t *testing.T) *httptest.Server {
	t.Helper()
	handler := NewHandler(nil) // will use real registry once implemented
	mux := NewServeMux(handler)
	return httptest.NewServer(mux)
}

func decodeList(t *testing.T, resp *http.Response) contractListResponse {
	t.Helper()
	var body contractListResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	return body
}

func decodeSingle(t *testing.T, resp *http.Response) contractSingleResponse {
	t.Helper()
	var body contractSingleResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode single response: %v", err)
	}
	return body
}

func decodeError(t *testing.T, resp *http.Response) contractErrorResponse {
	t.Helper()
	var body contractErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	return body
}

// =============================================================================
// Scenario: List connections with no active connections
// =============================================================================

func TestAC_ListConnections_EmptyState(t *testing.T) {
	// GIVEN: no connections have been established
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request the connection list
	resp := mustGet(t, ts.URL+"/api/connections")
	defer resp.Body.Close()

	// THEN: I receive an empty list with correct pagination
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := decodeList(t, resp)

	if len(body.Connections) != 0 {
		t.Errorf("expected 0 connections, got %d", len(body.Connections))
	}
	if body.Pagination != nil && body.Pagination.Total != nil && *body.Pagination.Total != 0 {
		t.Errorf("expected total=0, got %d", *body.Pagination.Total)
	}
	if body.Pagination != nil && body.Pagination.TotalPages != nil && *body.Pagination.TotalPages != 0 {
		t.Errorf("expected total_pages=0, got %d", *body.Pagination.TotalPages)
	}
}

// =============================================================================
// Scenario: List connections returns all active connections
// =============================================================================

func TestAC_ListConnections_ReturnsActiveConnections(t *testing.T) {
	// GIVEN: there are active connections in the registry
	// (This test will rely on the registry being populated;
	//  once implemented, we inject a pre-populated registry)
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request the connection list
	resp := mustGet(t, ts.URL+"/api/connections")
	defer resp.Body.Close()

	// THEN: the response includes connection objects with all required fields
	body := decodeList(t, resp)

	for i, conn := range body.Connections {
		if conn.ID == "" {
			t.Errorf("connection[%d].id must not be empty", i)
		}
		if conn.ServeName == "" {
			t.Errorf("connection[%d].serve_name must not be empty", i)
		}
		if conn.RemoteAddr == "" {
			t.Errorf("connection[%d].remote_addr must not be empty", i)
		}
		if conn.Status == "" {
			t.Errorf("connection[%d].status must not be empty", i)
		}
		if conn.Type == "" {
			t.Errorf("connection[%d].type must not be empty", i)
		}
		if conn.ConnectedAt.IsZero() {
			t.Errorf("connection[%d].connected_at must not be zero", i)
		}
	}
}

// =============================================================================
// Scenario: Filter connections by serve_name
// =============================================================================

func TestAC_ListConnections_FilterByServeName(t *testing.T) {
	// GIVEN: connections exist for multiple serve_names
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I filter by serve_name=test-service
	resp := mustGet(t, ts.URL+"/api/connections?serve_name=test-service")
	defer resp.Body.Close()

	// THEN: only connections for that serve_name are returned
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := decodeList(t, resp)
	for i, conn := range body.Connections {
		if conn.ServeName != "test-service" {
			t.Errorf("connection[%d] has serve_name=%q, expected 'test-service'", i, conn.ServeName)
		}
	}
}

// =============================================================================
// Scenario: Filter connections by status
// =============================================================================

func TestAC_ListConnections_FilterByStatus(t *testing.T) {
	// GIVEN: connections exist with various statuses
	ts := setupACServer(t)
	defer ts.Close()

	for _, status := range []string{"active", "idle", "terminated"} {
		t.Run("status="+status, func(t *testing.T) {
			// WHEN: I filter by status
			resp := mustGet(t, ts.URL+"/api/connections?status="+status)
			defer resp.Body.Close()

			// THEN: only connections with that status are returned
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}
			body := decodeList(t, resp)
			for i, conn := range body.Connections {
				if conn.Status != status {
					t.Errorf("connection[%d] has status=%q, expected %q", i, conn.Status, status)
				}
			}
		})
	}
}

// =============================================================================
// Scenario: Filter connections by type
// =============================================================================

func TestAC_ListConnections_FilterByType(t *testing.T) {
	// GIVEN: connections exist with both control and worker types
	ts := setupACServer(t)
	defer ts.Close()

	for _, connType := range []string{"control", "worker"} {
		t.Run("type="+connType, func(t *testing.T) {
			// WHEN: I filter by type
			resp := mustGet(t, ts.URL+"/api/connections?type="+connType)
			defer resp.Body.Close()

			// THEN: only connections with that type are returned
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}
			body := decodeList(t, resp)
			for i, conn := range body.Connections {
				if conn.Type != connType {
					t.Errorf("connection[%d] has type=%q, expected %q", i, conn.Type, connType)
				}
			}
		})
	}
}

// =============================================================================
// Scenario: Pagination works correctly
// =============================================================================

func TestAC_ListConnections_Pagination(t *testing.T) {
	// GIVEN: there are connections in the system
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request page 1 with page_size=5
	resp := mustGet(t, ts.URL+"/api/connections?page=1&page_size=5")
	defer resp.Body.Close()

	// THEN: at most 5 connections are returned and pagination reflects the total
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := decodeList(t, resp)

	if len(body.Connections) > 5 {
		t.Errorf("expected at most 5 connections, got %d", len(body.Connections))
	}
	if body.Pagination != nil && body.Pagination.PageSize != nil && *body.Pagination.PageSize != 5 {
		t.Errorf("expected page_size=5 in response, got %d", *body.Pagination.PageSize)
	}
	if body.Pagination != nil && body.Pagination.Page != nil && *body.Pagination.Page != 1 {
		t.Errorf("expected page=1 in response, got %d", *body.Pagination.Page)
	}
}

func TestAC_ListConnections_PaginationTotalPagesCalculation(t *testing.T) {
	// GIVEN: there are connections in the system
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request with a specific page_size
	resp := mustGet(t, ts.URL+"/api/connections?page_size=10")
	defer resp.Body.Close()

	// THEN: total_pages = ceil(total / page_size)
	body := decodeList(t, resp)
	if body.Pagination != nil && body.Pagination.Total != nil && body.Pagination.TotalPages != nil && body.Pagination.PageSize != nil {
		total := *body.Pagination.Total
		pageSize := *body.Pagination.PageSize
		totalPages := *body.Pagination.TotalPages

		expectedPages := (total + pageSize - 1) / pageSize
		if total == 0 {
			expectedPages = 0
		}
		if totalPages != expectedPages {
			t.Errorf("expected total_pages=%d (total=%d, page_size=%d), got %d",
				expectedPages, total, pageSize, totalPages)
		}
	}
}

func TestAC_ListConnections_PageBeyondRange(t *testing.T) {
	// GIVEN: there are few or no connections
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request a page far beyond the available data
	resp := mustGet(t, ts.URL+"/api/connections?page=9999")
	defer resp.Body.Close()

	// THEN: I get 200 with an empty connections list, not a 404
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for out-of-range page, got %d", resp.StatusCode)
	}

	body := decodeList(t, resp)
	if len(body.Connections) != 0 {
		t.Errorf("expected 0 connections on out-of-range page, got %d", len(body.Connections))
	}
}

// =============================================================================
// Scenario: Get specific connection by ID
// =============================================================================

func TestAC_GetConnection_ExistingConnection(t *testing.T) {
	// GIVEN: a connection exists with a known ID
	// (once implemented, the test should pre-populate the registry)
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request that specific connection
	resp := mustGet(t, ts.URL+"/api/connections/test-conn-1")
	defer resp.Body.Close()

	// THEN: I receive the connection details (or 404 if not populated)
	if resp.StatusCode == http.StatusOK {
		body := decodeSingle(t, resp)
		if body.Connection == nil {
			t.Fatal("connection must not be nil in 200 response")
		}
		if body.Connection.ID != "test-conn-1" {
			t.Errorf("expected connection ID 'test-conn-1', got %q", body.Connection.ID)
		}
	} else if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 200 or 404, got %d", resp.StatusCode)
	}
}

func TestAC_GetConnection_NotFound(t *testing.T) {
	// GIVEN: no connection exists with ID "does-not-exist"
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request that connection
	resp := mustGet(t, ts.URL+"/api/connections/does-not-exist")
	defer resp.Body.Close()

	// THEN: I receive a 404 with a meaningful error message
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	errBody := decodeError(t, resp)
	if errBody.Error == nil || *errBody.Error != "not_found" {
		t.Errorf("expected error='not_found', got %v", errBody.Error)
	}
	if errBody.Message == nil || *errBody.Message == "" {
		t.Error("error message must not be empty")
	}
}

// =============================================================================
// Scenario: Combining multiple filters
// =============================================================================

func TestAC_ListConnections_CombinedFilters(t *testing.T) {
	// GIVEN: connections exist with various properties
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I apply serve_name + status + type filters together
	resp := mustGet(t, ts.URL+"/api/connections?serve_name=svc1&status=active&type=control")
	defer resp.Body.Close()

	// THEN: only connections matching ALL filters are returned
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := decodeList(t, resp)
	for i, conn := range body.Connections {
		if conn.ServeName != "svc1" {
			t.Errorf("connection[%d] serve_name=%q, want 'svc1'", i, conn.ServeName)
		}
		if conn.Status != "active" {
			t.Errorf("connection[%d] status=%q, want 'active'", i, conn.Status)
		}
		if conn.Type != "control" {
			t.Errorf("connection[%d] type=%q, want 'control'", i, conn.Type)
		}
	}
}

// =============================================================================
// Scenario: Error handling for invalid parameters
// =============================================================================

func TestAC_ListConnections_InvalidStatusReturnsHelpfulError(t *testing.T) {
	// GIVEN: nothing special
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I provide an invalid status value
	resp := mustGet(t, ts.URL+"/api/connections?status=unknown")
	defer resp.Body.Close()

	// THEN: I receive a 400 with a message explaining valid values
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	errBody := decodeError(t, resp)
	if errBody.Message == nil || *errBody.Message == "" {
		t.Error("error message should explain valid status values")
	}
}

func TestAC_ListConnections_InvalidTypeReturnsHelpfulError(t *testing.T) {
	// GIVEN: nothing special
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I provide an invalid type value
	resp := mustGet(t, ts.URL+"/api/connections?type=unknown")
	defer resp.Body.Close()

	// THEN: I receive a 400 with a message explaining valid values
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	errBody := decodeError(t, resp)
	if errBody.Message == nil || *errBody.Message == "" {
		t.Error("error message should explain valid type values")
	}
}

func TestAC_ListConnections_NonIntegerPageReturnsError(t *testing.T) {
	// GIVEN: nothing special
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I provide a non-integer page value
	resp := mustGet(t, ts.URL+"/api/connections?page=abc")
	defer resp.Body.Close()

	// THEN: I receive a 400 error
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestAC_ListConnections_PageSizeTooLargeReturnsError(t *testing.T) {
	// GIVEN: nothing special
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request page_size > 100
	resp := mustGet(t, ts.URL+"/api/connections?page_size=200")
	defer resp.Body.Close()

	// THEN: I receive a 400 error
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// =============================================================================
// Scenario: Connections array is never null
// =============================================================================

func TestAC_ListConnections_ConnectionsArrayNeverNull(t *testing.T) {
	// GIVEN: no connections exist
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request the list
	resp := mustGet(t, ts.URL+"/api/connections")
	defer resp.Body.Close()

	// THEN: connections is an empty array [], not null
	var raw map[string]json.RawMessage
	json.NewDecoder(resp.Body).Decode(&raw)

	connsRaw := raw["connections"]
	if string(connsRaw) == "null" {
		t.Error("connections must be an empty array [], not null")
	}
}

// =============================================================================
// Scenario: Method not allowed
// =============================================================================

func TestAC_ListConnections_PostNotAllowed(t *testing.T) {
	// GIVEN: the API only supports GET
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I send a POST request
	resp, err := http.Post(ts.URL+"/api/connections", "application/json", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// THEN: I receive 405 Method Not Allowed
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestAC_GetConnection_PutNotAllowed(t *testing.T) {
	// GIVEN: the API only supports GET
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I send a PUT request
	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/connections/some-id", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// THEN: I receive 405 Method Not Allowed
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

// =============================================================================
// Scenario: Empty connection ID in path
// =============================================================================

func TestAC_GetConnection_EmptyIDReturnsBadRequest(t *testing.T) {
	// GIVEN: nothing special
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I request /api/connections/ with no ID
	resp := mustGet(t, ts.URL+"/api/connections/")
	defer resp.Body.Close()

	// THEN: I receive a 400 explaining that an ID is required
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}

	errBody := decodeError(t, resp)
	if errBody.Error == nil {
		t.Error("error field is required")
	}
}

// =============================================================================
// Scenario: Serve_name filter with no matches returns empty list
// =============================================================================

func TestAC_ListConnections_NoMatchingServeName(t *testing.T) {
	// GIVEN: no connections match the given serve_name
	ts := setupACServer(t)
	defer ts.Close()

	// WHEN: I filter by a serve_name that doesn't exist
	resp := mustGet(t, ts.URL+fmt.Sprintf("/api/connections?serve_name=%s", "nonexistent-service-xyz"))
	defer resp.Body.Close()

	// THEN: I get 200 with empty list, not 404
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body := decodeList(t, resp)
	if len(body.Connections) != 0 {
		t.Errorf("expected 0 connections, got %d", len(body.Connections))
	}
}
