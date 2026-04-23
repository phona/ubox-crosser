package main

import (
	"encoding/json"
	"net/http"
	"os"
)

// GitSHA is the 7-character git SHA the binary was built from. It is
// overridden at build time via:
//
//	go build -ldflags "-X main.GitSHA=$(git rev-parse --short HEAD)" ./cmd/server
//
// When the ldflag is not passed (e.g. `go run`, `go test`), it stays
// "unknown".
var GitSHA = "unknown"

// goVersionLiteral is the hard-coded Go toolchain string returned by
// /buildinfo. Kept as a literal rather than runtime.Version() because the
// acceptance spec wants a stable "go1.23" value, not a patch-level string
// like "go1.23.12". A unit test cross-checks this against go.mod to catch
// drift when the toolchain is bumped.
const goVersionLiteral = "go1.23"

// BuildInfo is the response shape returned by GET /buildinfo.
type BuildInfo struct {
	GitSHA    string `json:"git_sha"`
	BuildID   string `json:"build_id"`
	GoVersion string `json:"go_version"`
}

// BuildInfoHandler serves GET /buildinfo. It reads BUILD_ID from the
// environment on each request (so operators can change it without a
// restart) and emits a fixed-shape JSON document.
func BuildInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := BuildInfo{
		GitSHA:    GitSHA,
		BuildID:   buildIDFromEnv(),
		GoVersion: goVersionLiteral,
	}
	body, _ := json.Marshal(info)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func buildIDFromEnv() string {
	if v := os.Getenv("BUILD_ID"); v != "" {
		return v
	}
	return "dev"
}
