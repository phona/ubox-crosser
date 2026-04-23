package server

import (
	"encoding/json"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

// startTime records when the HTTP server was initialized.
var startTime = time.Now()

// BuildInfo holds version metadata injected at build time or read from env.
type BuildInfo struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

// healthzResponse is the JSON body for /healthz.
type healthzResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp     int64  `json:"timestamp"`
}

// uptimeSeconds exposes the elapsed seconds since startTime (used in tests).
var uptimeSeconds atomic.Int64

func buildInfoFromVars(gitSHA string) BuildInfo {
	buildID := os.Getenv("BUILD_ID")
	if buildID == "" {
		buildID = "dev"
	}
	return BuildInfo{
		GitSHA:    gitSHA,
		BuildID:   buildID,
		GoVersion: "go1.23",
	}
}

// StartHTTPServer starts the HTTP server on addr, serving /healthz and
// /buildinfo.  gitSHA is typically injected via ldflags at build time.
func StartHTTPServer(addr, gitSHA string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		resp := healthzResponse{
			Status:        "healthy",
			UptimeSeconds: int64(now.Sub(startTime).Seconds()),
			Timestamp:     now.Unix(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	})

	mux.HandleFunc("/buildinfo", func(w http.ResponseWriter, r *http.Request) {
		info := buildInfoFromVars(gitSHA)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info) //nolint:errcheck
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return srv.ListenAndServe()
}
