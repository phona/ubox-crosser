package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

var healthBody = []byte(`{"status":"ok"}`)

func newHealthHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(healthBody)
	})
	return mux
}

func StartHealthServer(addr string, errs chan<- error) {
	if addr == "" {
		return
	}
	go func() {
		if err := http.ListenAndServe(addr, newHealthHandler()); err != nil {
			log.Errorf("Health server error: %v", err)
			errs <- err
		}
	}()
}
