package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type HealthCheckServer struct {
	address   string
	startTime time.Time
}

type HealthResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp     int64  `json:"timestamp"`
}

func NewHealthCheckServer(address string) *HealthCheckServer {
	return &HealthCheckServer{
		address:   address,
		startTime: time.Now(),
	}
}

func (h *HealthCheckServer) Start() error {
	http.HandleFunc("/healthz", h.handleHealthz)
	log.Infof("Starting health check server on %s", h.address)
	return http.ListenAndServe(h.address, nil)
}

func (h *HealthCheckServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	uptime := int64(now.Sub(h.startTime).Seconds())

	response := HealthResponse{
		Status:        "healthy",
		UptimeSeconds: uptime,
		Timestamp:     now.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("Error writing health check response: %v", err)
	}
}

func (h *HealthCheckServer) GetUptime() time.Duration {
	return time.Since(h.startTime)
}
