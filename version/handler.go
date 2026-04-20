package version

import (
	"encoding/json"
	"net/http"
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
}

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Info{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
	})
}
