package server

import (
	"encoding/json"
	"net/http"
)

// newStatsHandler returns an http.Handler that serves GET /api/stats.
func newStatsHandler(collector *Collector) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		snapshot := collector.Snapshot()
		json.NewEncoder(w).Encode(snapshot)
	})
	return mux
}
