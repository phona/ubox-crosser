package server

import (
	"net/http"

	"github.com/phona/ubox-crosser/version"
)

func NewAdminMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/version", version.Handler)
	return mux
}
