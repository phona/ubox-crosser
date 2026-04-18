package stats

import (
	"encoding/json"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Server serves stats over HTTP.
type Server struct {
	collector *Collector
	mux       *http.ServeMux
}

// NewServer creates a new stats HTTP server.
func NewServer(collector *Collector) *Server {
	s := &Server{
		collector: collector,
		mux:       http.NewServeMux(),
	}
	s.mux.HandleFunc("/api/stats", s.handleStats)
	return s
}

// ListenAndServe starts the stats HTTP server on the given address.
func (s *Server) ListenAndServe(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Infof("Stats server listening on %s", address)
	return http.Serve(listener, s.mux)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := map[string]string{"error": "method not allowed"}
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
		return
	}

	snap := s.collector.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snap) //nolint:errcheck
}
