package api

import (
	"net/http"
)

func NewServeMux(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/connections", handler.ListConnections)
	mux.HandleFunc("/api/connections/", handler.GetConnection)
	return mux
}

func StartManagementServer(addr string, registry RegistryReader) error {
	handler := NewHandler(registry)
	mux := NewServeMux(handler)
	return http.ListenAndServe(addr, mux)
}
