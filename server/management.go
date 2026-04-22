package server

import (
	"encoding/json"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
)

type VersionInfo struct {
	Version string `json:"version"`
	Module  string `json:"module"`
	GoOS    string `json:"go_os"`
	GoArch  string `json:"go_arch"`
	GitSha  string `json:"git_sha"`
}

type HealthzResponse struct {
	Status        string `json:"status"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

type ManagementServer struct {
	address   string
	mux       *http.ServeMux
	errs      chan error
	startTime time.Time
}

func NewManagementServer(address string) *ManagementServer {
	m := &ManagementServer{
		address:   address,
		mux:       http.NewServeMux(),
		errs:      make(chan error, 10),
		startTime: time.Now(),
	}
	m.registerHandlers()
	return m
}

func (m *ManagementServer) registerHandlers() {
	m.mux.HandleFunc("/version", m.handleVersion)
	m.mux.HandleFunc("/health", m.handleHealth)
	m.mux.HandleFunc("/healthz", m.handleHealthz)
}

func (m *ManagementServer) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	info := VersionInfo{
		Module:  "github.com/phona/ubox-crosser",
		GoOS:    runtime.GOOS,
		GoArch:  runtime.GOARCH,
		Version: "unknown",
		GitSha:  "unknown",
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		info.Version = buildInfo.Main.Version
		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" {
				info.GitSha = setting.Value
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func (m *ManagementServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (m *ManagementServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uptime := int64(time.Since(m.startTime).Seconds())
	response := HealthzResponse{
		Status:        "ok",
		UptimeSeconds: uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (m *ManagementServer) Listen() {
	log.Infof("Starting management server on %s", m.address)
	if err := http.ListenAndServe(m.address, m.mux); err != nil {
		m.errs <- err
	}
}

func (m *ManagementServer) Err() error {
	select {
	case err := <-m.errs:
		return err
	default:
		return nil
	}
}
