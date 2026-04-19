package version

import (
	"encoding/json"
	"testing"
)

func TestGetInfoDefaults(t *testing.T) {
	info := GetInfo()
	if info.Version != "dev" {
		t.Errorf("expected default version 'dev', got %q", info.Version)
	}
	if info.GitCommit != "" {
		t.Errorf("expected empty git_commit, got %q", info.GitCommit)
	}
	if info.BuildTime != "" {
		t.Errorf("expected empty build_time, got %q", info.BuildTime)
	}
}

func TestInfoJSON(t *testing.T) {
	info := Info{Version: "1.0.0", GitCommit: "abc1234", BuildTime: "2026-01-01T00:00:00Z"}
	b := info.JSON()

	var parsed map[string]string
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if parsed["version"] != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", parsed["version"])
	}
	if parsed["git_commit"] != "abc1234" {
		t.Errorf("expected git_commit 'abc1234', got %q", parsed["git_commit"])
	}
	if parsed["build_time"] != "2026-01-01T00:00:00Z" {
		t.Errorf("expected build_time '2026-01-01T00:00:00Z', got %q", parsed["build_time"])
	}
}
