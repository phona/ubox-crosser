package webhookdebug

import (
	"encoding/json"
	"io"
	"net/http"
)

type RequestInfo struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	info := RequestInfo{
		Method:  r.Method,
		Path:    r.URL.Path,
		Query:   r.URL.Query(),
		Headers: r.Header,
		Body:    string(body),
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(info)
}
