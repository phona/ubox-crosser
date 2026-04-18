package server

import (
	"encoding/json"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// NewStatsHandler returns an http.Handler that serves GET /api/stats.
// Non-GET methods receive 405 Method Not Allowed.
func NewStatsHandler(collector *Collector) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		snap := collector.Snapshot()
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			log.Errorf("stats handler encode error: %v", err)
		}
	})
	return mux
}

// StartStatsServer starts the stats HTTP server on the given address.
// It blocks until the listener fails or the server is shut down.
func StartStatsServer(address string, collector *Collector) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Infof("Stats server listening on %s", address)
	return http.Serve(ln, NewStatsHandler(collector))
}
