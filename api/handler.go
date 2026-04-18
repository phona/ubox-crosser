package api

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"ubox-crosser/models/connection"
)

type FilterOptions struct {
	ServeName string
	Status    connection.Status
	Type      connection.Type
	Page      int
	PageSize  int
}

type RegistryReader interface {
	Get(id string) (*connection.Connection, bool)
	List(filter FilterOptions) ([]*connection.Connection, int)
}

type Handler struct {
	registry RegistryReader
}

func NewHandler(registry RegistryReader) *Handler {
	return &Handler{registry: registry}
}

func (h *Handler) ListConnections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "method_not_allowed",
			Message: "Only GET method is allowed",
		})
		return
	}

	q := r.URL.Query()

	// Validate status
	statusStr := q.Get("status")
	if statusStr != "" {
		validStatuses := map[string]bool{"active": true, "idle": true, "terminated": true}
		if !validStatuses[statusStr] {
			writeError(w, http.StatusBadRequest, "invalid_parameter",
				"Invalid status value. Must be one of: active, idle, terminated")
			return
		}
	}

	// Validate type
	typeStr := q.Get("type")
	if typeStr != "" {
		validTypes := map[string]bool{"control": true, "worker": true}
		if !validTypes[typeStr] {
			writeError(w, http.StatusBadRequest, "invalid_parameter",
				"Invalid type value. Must be one of: control, worker")
			return
		}
	}

	// Validate page
	page := 1
	if pageStr := q.Get("page"); pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			writeError(w, http.StatusBadRequest, "invalid_parameter",
				"Invalid page value. Must be a positive integer")
			return
		}
	}

	// Validate page_size
	pageSize := 20
	if pageSizeStr := q.Get("page_size"); pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 || pageSize > 100 {
			writeError(w, http.StatusBadRequest, "invalid_parameter",
				"Invalid page_size value. Must be between 1 and 100")
			return
		}
	}

	filter := FilterOptions{
		ServeName: q.Get("serve_name"),
		Status:    connection.Status(statusStr),
		Type:      connection.Type(typeStr),
		Page:      page,
		PageSize:  pageSize,
	}

	var conns []*connection.Connection
	var total int

	if h.registry != nil {
		conns, total = h.registry.List(filter)
	}

	if conns == nil {
		conns = make([]*connection.Connection, 0)
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	resp := ConnectionListResponse{
		Connections: conns,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "method_not_allowed",
			Message: "Only GET method is allowed",
		})
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/connections/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid_parameter", "Connection ID is required")
		return
	}

	if h.registry == nil {
		writeError(w, http.StatusNotFound, "not_found", "Connection not found")
		return
	}

	conn, ok := h.registry.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "Connection not found")
		return
	}

	resp := ConnectionResponse{Connection: conn}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, status int, errCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errCode,
		Message: message,
	})
}
