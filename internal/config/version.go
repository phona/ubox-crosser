package config

import "fmt"

const Version = "0.1.0"

var (
	GitCommit string
	BuildTime string
)

func FullVersion() string {
	commit := GitCommit
	if commit == "" {
		commit = "unknown"
	}
	built := BuildTime
	if built == "" {
		built = "unknown"
	}
	return fmt.Sprintf("%s (commit: %s, built: %s)", Version, commit, built)
}
