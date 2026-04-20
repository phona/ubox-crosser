package config

import (
	"regexp"
	"strings"
	"testing"
)

var semverRegex = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)` +
	`(-((0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?` +
	`(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$`)

var gitHashRegex = regexp.MustCompile(`^[0-9a-f]{7,40}$`)

var iso8601Regex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

func TestVersionConstantNotEmpty(t *testing.T) {
	if Version == "" {
		t.Fatal("Version constant must not be empty")
	}
}

func TestVersionConstantSemver(t *testing.T) {
	if !semverRegex.MatchString(Version) {
		t.Fatalf("Version %q does not match semver format", Version)
	}
}

func TestFullVersionContainsVersion(t *testing.T) {
	fv := FullVersion()
	if !strings.HasPrefix(fv, Version) {
		t.Fatalf("FullVersion() = %q does not start with Version %q", fv, Version)
	}
}

func TestFullVersionFormat(t *testing.T) {
	fv := FullVersion()
	if !strings.Contains(fv, "commit:") {
		t.Fatalf("FullVersion() = %q missing 'commit:' substring", fv)
	}
	if !strings.Contains(fv, "built:") {
		t.Fatalf("FullVersion() = %q missing 'built:' substring", fv)
	}
}

func TestFullVersionWithoutLdflags(t *testing.T) {
	origCommit := GitCommit
	origBuild := BuildTime
	GitCommit = ""
	BuildTime = ""
	defer func() {
		GitCommit = origCommit
		BuildTime = origBuild
	}()

	fv := FullVersion()
	if !strings.Contains(fv, "unknown") {
		t.Fatalf("FullVersion() with empty ldflags = %q, want 'unknown' placeholder", fv)
	}
}

func TestFullVersionWithLdflags(t *testing.T) {
	origCommit := GitCommit
	origBuild := BuildTime
	GitCommit = "abc1234"
	BuildTime = "2026-04-20T12:00:00Z"
	defer func() {
		GitCommit = origCommit
		BuildTime = origBuild
	}()

	fv := FullVersion()
	if !strings.Contains(fv, "abc1234") {
		t.Fatalf("FullVersion() = %q does not contain GitCommit 'abc1234'", fv)
	}
	if !strings.Contains(fv, "2026-04-20T12:00:00Z") {
		t.Fatalf("FullVersion() = %q does not contain BuildTime '2026-04-20T12:00:00Z'", fv)
	}
}

func TestGitCommitDefaultEmpty(t *testing.T) {
	if GitCommit != "" {
		t.Fatalf("GitCommit default should be empty, got %q", GitCommit)
	}
}

func TestBuildTimeDefaultEmpty(t *testing.T) {
	if BuildTime != "" {
		t.Fatalf("BuildTime default should be empty, got %q", BuildTime)
	}
}

func TestGitCommitFormatWhenSet(t *testing.T) {
	origCommit := GitCommit
	GitCommit = "abc1234"
	defer func() { GitCommit = origCommit }()

	if !gitHashRegex.MatchString(GitCommit) {
		t.Fatalf("GitCommit %q does not match git short hash format", GitCommit)
	}
}

func TestBuildTimeFormatWhenSet(t *testing.T) {
	origBuild := BuildTime
	BuildTime = "2026-04-20T12:00:00Z"
	defer func() { BuildTime = origBuild }()

	if !iso8601Regex.MatchString(BuildTime) {
		t.Fatalf("BuildTime %q does not match ISO 8601 format", BuildTime)
	}
}

func TestFullVersionOutputFormat(t *testing.T) {
	origCommit := GitCommit
	origBuild := BuildTime
	GitCommit = "def5678"
	BuildTime = "2026-01-01T00:00:00Z"
	defer func() {
		GitCommit = origCommit
		BuildTime = origBuild
	}()

	fv := FullVersion()
	expected := Version + " (commit: def5678, built: 2026-01-01T00:00:00Z)"
	if fv != expected {
		t.Fatalf("FullVersion() = %q, want %q", fv, expected)
	}
}

func TestFullVersionUnknownFormat(t *testing.T) {
	origCommit := GitCommit
	origBuild := BuildTime
	GitCommit = ""
	BuildTime = ""
	defer func() {
		GitCommit = origCommit
		BuildTime = origBuild
	}()

	fv := FullVersion()
	expected := Version + " (commit: unknown, built: unknown)"
	if fv != expected {
		t.Fatalf("FullVersion() = %q, want %q", fv, expected)
	}
}
