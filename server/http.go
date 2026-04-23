package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var httpStartTime = time.Now()

type healthzResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp     int64  `json:"timestamp"`
}

type buildInfoResponse struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

// HealthzHandler returns service health and uptime.
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	resp := healthzResponse{
		Status:        "healthy",
		UptimeSeconds: int64(time.Since(httpStartTime).Seconds()),
		Timestamp:     time.Now().Unix(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BuildInfoHandler returns a handler that reports build metadata.
// gitSHA is injected at link time via -X main.GitSHA=…; buildID is read
// from the BUILD_ID env var at request time, defaulting to "dev".
func BuildInfoHandler(gitSHA string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildID := os.Getenv("BUILD_ID")
		if buildID == "" {
			buildID = "dev"
		}
		resp := buildInfoResponse{
			GitSHA:    gitSHA,
			BuildID:   buildID,
			GoVersion: "go1.23",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// StartHTTPServer registers /healthz and /buildinfo and begins listening.
func StartHTTPServer(addr, gitSHA string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthzHandler)
	mux.HandleFunc("/buildinfo", BuildInfoHandler(gitSHA))
	go http.ListenAndServe(addr, mux) //nolint:errcheck
}
