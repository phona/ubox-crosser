package api

import (
	"encoding/json"
	"net/http"
)

type echoResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func EchoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(errorResponse{Error: "method not allowed"})
			return
		}

		query := r.URL.Query()
		if !query.Has("msg") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: "msg parameter is required"})
			return
		}

		msg := query.Get("msg")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(echoResponse{Message: msg})
	}
}
