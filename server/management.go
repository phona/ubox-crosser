package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type buildInfoResponse struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

type healthzResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp     int64  `json:"timestamp"`
}

// ManagementServer serves HTTP management endpoints on a dedicated address.
type ManagementServer struct {
	startTime time.Time
	gitSHA    string
	buildID   string
	mux       *http.ServeMux
}

// NewManagementServer creates a ManagementServer. gitSHA is injected via ldflags at build time.
// BUILD_ID is read from the environment; defaults to "dev".
func NewManagementServer(gitSHA string) *ManagementServer {
	buildID := os.Getenv("BUILD_ID")
	if buildID == "" {
		buildID = "dev"
	}
	ms := &ManagementServer{
		startTime: time.Now(),
		gitSHA:    gitSHA,
		buildID:   buildID,
		mux:       http.NewServeMux(),
	}
	ms.mux.HandleFunc("/healthz", ms.handleHealthz)
	ms.mux.HandleFunc("/buildinfo", ms.handleBuildinfo)
	return ms
}

func (ms *ManagementServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := healthzResponse{
		Status:        "healthy",
		UptimeSeconds: int64(time.Since(ms.startTime).Seconds()),
		Timestamp:     time.Now().Unix(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

func (ms *ManagementServer) handleBuildinfo(w http.ResponseWriter, r *http.Request) {
	resp := buildInfoResponse{
		GitSHA:    ms.gitSHA,
		BuildID:   ms.buildID,
		GoVersion: "go1.23",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

// ListenAndServe starts the HTTP management server on addr (e.g. ":8080").
func (ms *ManagementServer) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, ms.mux)
}
