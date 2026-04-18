package server

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// StartMetricsServer starts an HTTP server exposing /api/metrics.
// It blocks, so it should be called in a goroutine.
func StartMetricsServer(addr string, metrics *Metrics) {
	mux := http.NewServeMux()
	mux.Handle("/api/metrics", metricsHandler(metrics))

	log.Infof("Starting metrics HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Errorf("Metrics HTTP server error: %s", err)
	}
}

func metricsHandler(metrics *Metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		snapshot := metrics.Snapshot()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(snapshot); err != nil {
			log.Errorf("Error encoding metrics response: %s", err)
		}
	})
}
