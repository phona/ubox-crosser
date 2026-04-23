package server

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type HealthzResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Timestamp     int64  `json:"timestamp"`
}

type HealthChecker struct {
	startTime time.Time
	mu        sync.RWMutex
	addr      string
}

func NewHealthChecker(addr string) *HealthChecker {
	return &HealthChecker{
		startTime: time.Now(),
		addr:      addr,
	}
}

func (h *HealthChecker) Start() error {
	http.HandleFunc("/healthz", h.handleHealthz)
	log.Infof("Starting health check server on %s", h.addr)
	return http.ListenAndServe(h.addr, nil)
}

func (h *HealthChecker) handleHealthz(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	uptime := time.Since(h.startTime).Seconds()
	h.mu.RUnlock()

	response := HealthzResponse{
		Status:        "healthy",
		UptimeSeconds: int64(uptime),
		Timestamp:     time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
