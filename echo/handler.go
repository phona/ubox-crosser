package echo

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
