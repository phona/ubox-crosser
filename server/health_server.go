package server

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type HealthResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp     int64  `json:"timestamp"`
}

type HealthServer struct {
	port      string
	startTime time.Time
	mu        sync.RWMutex
}

func NewHealthServer(port string) *HealthServer {
	if port == "" {
		port = "8080"
	}
	return &HealthServer{
		port:      port,
		startTime: time.Now(),
	}
}

func (h *HealthServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	startTime := h.startTime
	h.mu.RUnlock()

	now := time.Now()
	uptime := now.Unix() - startTime.Unix()

	response := HealthResponse{
		Status:        "healthy",
		UptimeSeconds: uptime,
		Timestamp:     now.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *HealthServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealthz)

	addr := ":" + h.port
	log.Infof("Starting health check server on %s", addr)
	return http.ListenAndServe(addr, mux)
}
