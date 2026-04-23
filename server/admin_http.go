package server

import (
	"net/http"
	"time"
)

// StartAdminHTTP starts a minimal unauthenticated admin HTTP server on addr
// serving the supplied mux, and blocks until it exits. The caller is
// responsible for registering handlers on mux before invocation, and for
// running this inside a goroutine so the main proxy loop is not blocked.
//
// The admin listener intentionally has no TLS and no middleware — it exposes
// only read-only identification endpoints (e.g. /buildinfo, /healthz) and is
// expected to bind to an interface reachable only from the operations
// network.
func StartAdminHTTP(addr string, mux *http.ServeMux) error {
	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s.ListenAndServe()
}
