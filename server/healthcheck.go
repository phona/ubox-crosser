package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type BuildInfo struct {
	GitSha    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

type HealthCheckServer struct {
	httpServer *http.Server
	startTime  time.Time
}

func NewHealthCheckServer(listenAddr string, gitSha string) *HealthCheckServer {
	hc := &HealthCheckServer{
		startTime: time.Now(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", hc.handleHealthz)
	mux.HandleFunc("/buildinfo", hc.makeBuildInfoHandler(gitSha))

	hc.httpServer = &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	return hc
}

func (hc *HealthCheckServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uptime := int64(time.Since(hc.startTime).Seconds())
	response := map[string]interface{}{
		"status":          "healthy",
		"uptime_seconds":  uptime,
		"timestamp":       time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (hc *HealthCheckServer) makeBuildInfoHandler(gitSha string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		buildID := os.Getenv("BUILD_ID")
		if buildID == "" {
			buildID = "dev"
		}

		info := BuildInfo{
			GitSha:    gitSha,
			BuildID:   buildID,
			GoVersion: "go1.23",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

func (hc *HealthCheckServer) Start() error {
	return hc.httpServer.ListenAndServe()
}

func (hc *HealthCheckServer) Stop() error {
	return hc.httpServer.Close()
}
