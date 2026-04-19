package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/phona/ubox-crosser/internal/version"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(version.GetInfo().JSON())
}

func startHTTPServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", versionHandler)
	log.Infof("Starting HTTP listener on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Errorf("HTTP server error: %v", err)
	}
}
