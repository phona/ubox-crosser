package server

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	httpStartTime = time.Now()
	httpOnce      sync.Once
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

// StartHTTPServer starts the HTTP management server (idempotent).
func StartHTTPServer(addr, gitSHA string) {
	httpOnce.Do(func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/buildinfo", buildInfoHandler(gitSHA))
		mux.HandleFunc("/healthz", healthzHandler)
		go func() {
			if err := http.ListenAndServe(addr, mux); err != nil {
				// non-fatal: HTTP server failure doesn't kill proxy
				_ = err
			}
		}()
	})
}

func buildInfoHandler(gitSHA string) http.HandlerFunc {
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
		data, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	resp := healthzResponse{
		Status:        "healthy",
		UptimeSeconds: int64(time.Since(httpStartTime).Seconds()),
		Timestamp:     time.Now().Unix(),
	}
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
