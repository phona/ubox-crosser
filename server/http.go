package server

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const defaultHTTPAddr = ":8080"

func (p *ProxyServer) StartHTTPServer(addr string) {
	if addr == "" {
		addr = defaultHTTPAddr
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", handlePing)

	log.Infof("Starting HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Errorf("HTTP server error: %s", err)
		p.errs <- err
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"message": "pong"}
	json.NewEncoder(w).Encode(resp)
}
