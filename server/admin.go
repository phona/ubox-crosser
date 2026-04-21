package server

import (
	"net/http"

	"github.com/phona/ubox-crosser/webhookdebug"
)

func NewAdminMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook-debug", webhookdebug.Handler)
	return mux
}
