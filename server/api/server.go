package api

import (
	"context"
	"net/http"
)

func NewHTTPServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/echo", EchoHandler())
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func ListenAndServe(srv *http.Server) error {
	return srv.ListenAndServe()
}

func Shutdown(ctx context.Context, srv *http.Server) error {
	return srv.Shutdown(ctx)
}
