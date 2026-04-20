package api

import (
	"encoding/json"
	"net/http"

	"github.com/phona/ubox-crosser/version"
)

type versionResponse struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionResponse{
		Version:   version.Version,
		Commit:    version.Commit,
		BuildTime: version.BuildTime,
	})
}
