package server

import (
	"encoding/json"
	"net/http"
	"time"
)

var startTime = time.Now()

// BuildInfo holds build-time metadata injected via ldflags and environment.
type BuildInfo struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

// HealthzResponse is the JSON shape returned by /healthz.
type HealthzResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp     int64  `json:"timestamp"`
}

// NewBuildInfoHandler returns an http.HandlerFunc that writes build metadata JSON.
func NewBuildInfoHandler(info BuildInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info) //nolint:errcheck
	}
}

// NewHealthzHandler returns an http.HandlerFunc that writes service health JSON.
func NewHealthzHandler(start time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthzResponse{ //nolint:errcheck
			Status:        "healthy",
			UptimeSeconds: int64(now.Sub(start).Seconds()),
			Timestamp:     now.Unix(),
		})
	}
}

// StartHTTPServer starts an HTTP server on addr serving /healthz and /buildinfo.
// Returns immediately; the server runs in a background goroutine.
func StartHTTPServer(addr string, info BuildInfo) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", NewHealthzHandler(startTime))
	mux.HandleFunc("/buildinfo", NewBuildInfoHandler(info))
	go http.ListenAndServe(addr, mux) //nolint:errcheck
}
