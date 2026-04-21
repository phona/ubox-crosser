package server

import (
	"net/http"

	"github.com/phona/ubox-crosser/healthz"
	"github.com/phona/ubox-crosser/webhookdebug"
)

func NewAdminMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook-debug", webhookdebug.Handler)
	mux.HandleFunc("/api/healthz", healthz.Handler)
	return mux
}
