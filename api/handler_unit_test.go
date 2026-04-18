package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"ubox-crosser/models/connection"
)

// testRegistry implements RegistryReader for unit tests
type testRegistry struct {
	conns map[string]*connection.Connection
}

func newTestRegistry(conns ...*connection.Connection) *testRegistry {
	r := &testRegistry{conns: make(map[string]*connection.Connection)}
	for _, c := range conns {
		r.conns[c.ID] = c
	}
	return r
}

func (r *testRegistry) Get(id string) (*connection.Connection, bool) {
	c, ok := r.conns[id]
	return c, ok
}

func (r *testRegistry) List(filter FilterOptions) ([]*connection.Connection, int) {
	var filtered []*connection.Connection
	for _, c := range r.conns {
		if filter.ServeName != "" && c.ServeName != filter.ServeName {
			continue
		}
		if filter.Status != "" && c.Status != filter.Status {
			continue
		}
		if filter.Type != "" && c.Type != filter.Type {
			continue
		}
		filtered = append(filtered, c)
	}
	total := len(filtered)
	start := (filter.Page - 1) * filter.PageSize
	if start >= total {
		return []*connection.Connection{}, total
	}
	end := start + filter.PageSize
	if end > total {
		end = total
	}
	return filtered[start:end], total
}

func makeConn(id, serveName, addr string, status connection.Status, typ connection.Type) *connection.Connection {
	return &connection.Connection{
		ID:          id,
		ServeName:   serveName,
		RemoteAddr:  addr,
		Status:      status,
		Type:        typ,
		ConnectedAt: time.Now(),
	}
}

func TestUnit_ListConnections_WithPopulatedRegistry(t *testing.T) {
	reg := newTestRegistry(
		makeConn("c1", "svc-a", "1.2.3.4:100", connection.StatusActive, connection.TypeControl),
		makeConn("c2", "svc-b", "5.6.7.8:200", connection.StatusIdle, connection.TypeWorker),
	)
	handler := NewHandler(reg)
	mux := NewServeMux(handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/connections")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body ConnectionListResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if len(body.Connections) != 2 {
		t.Errorf("expected 2 connections, got %d", len(body.Connections))
	}
	if body.Pagination.Total != 2 {
		t.Errorf("expected total=2, got %d", body.Pagination.Total)
	}
}

func TestUnit_GetConnection_Found(t *testing.T) {
	reg := newTestRegistry(
		makeConn("conn-42", "svc-x", "10.0.0.1:999", connection.StatusActive, connection.TypeControl),
	)
	handler := NewHandler(reg)
	mux := NewServeMux(handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/connections/conn-42")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body ConnectionResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Connection == nil || body.Connection.ID != "conn-42" {
		t.Errorf("expected connection with ID 'conn-42'")
	}
}

func TestUnit_ListConnections_FilterByServeName(t *testing.T) {
	reg := newTestRegistry(
		makeConn("c1", "svc-a", "1.1.1.1:1", connection.StatusActive, connection.TypeControl),
		makeConn("c2", "svc-b", "2.2.2.2:2", connection.StatusActive, connection.TypeWorker),
	)
	handler := NewHandler(reg)
	mux := NewServeMux(handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/connections?serve_name=svc-a")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body ConnectionListResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Pagination.Total != 1 {
		t.Errorf("expected total=1 for serve_name=svc-a, got %d", body.Pagination.Total)
	}
}

func TestUnit_ListConnections_PaginationMath(t *testing.T) {
	conns := make([]*connection.Connection, 25)
	for i := range conns {
		conns[i] = makeConn("c"+string(rune('A'+i)), "svc", "1.1.1.1:1", connection.StatusActive, connection.TypeControl)
	}
	reg := newTestRegistry(conns...)
	handler := NewHandler(reg)
	mux := NewServeMux(handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/connections?page=1&page_size=10")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body ConnectionListResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Pagination.TotalPages != 3 {
		t.Errorf("expected total_pages=3 for 25 items/page_size=10, got %d", body.Pagination.TotalPages)
	}
	if len(body.Connections) != 10 {
		t.Errorf("expected 10 connections on page 1, got %d", len(body.Connections))
	}
}
