package version

import (
	"encoding/json"
	"net/http"
)

var Commit string

type versionResponse struct {
	Commit string `json:"commit"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	commit := Commit
	if commit == "" {
		commit = "unknown"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(versionResponse{Commit: commit})
}
