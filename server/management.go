package server

import (
	"encoding/json"
	"net/http"

	"github.com/phona/ubox-crosser/version"
)

func NewManagementServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", handleVersion)
	return &http.Server{Addr: addr, Handler: mux}
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Version string `json:"version"`
	}{Version: version.Version})
}
