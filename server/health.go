package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func newHealthMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Set("Allow", "GET")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	return mux
}

func (p *ProxyServer) startHealthServer(addr string) {
	log.Infof("Health server listening on %s", addr)
	srv := &http.Server{Addr: addr, Handler: newHealthMux()}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		p.errs <- err
	}
}
