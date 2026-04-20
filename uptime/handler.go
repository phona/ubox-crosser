package uptime

import (
	"encoding/json"
	"net/http"
	"time"
)

var startTime time.Time

func Init() {
	startTime = time.Now()
}

type response struct {
	UptimeSeconds int `json:"uptime_seconds"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{
		UptimeSeconds: int(time.Since(startTime).Seconds()),
	})
}
