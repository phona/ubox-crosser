package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var startTime time.Time

// UptimeResponse is the JSON response for GET /api/uptime.
type UptimeResponse struct {
	StartedAt     string `json:"started_at"`
	UptimeSeconds int    `json:"uptime_seconds"`
}

// Init records the process start time. Must be called once at startup.
func Init() {
	startTime = time.Now()
}

// UptimeHandler returns an http.HandlerFunc for the /api/uptime endpoint.
func UptimeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := UptimeResponse{
			StartedAt:     startTime.Format(time.RFC3339),
			UptimeSeconds: int(time.Since(startTime).Seconds()),
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// StartServer starts an HTTP server on the given address with /api/uptime registered.
// Unknown routes return 404. Blocks, so call in a goroutine.
func StartServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/uptime", UptimeHandler())
	log.Infof("Starting HTTP API server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Errorf("HTTP API server error: %s", err)
	}
}
