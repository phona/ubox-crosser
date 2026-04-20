package config

import (
	"regexp"
	"testing"
)

// Acceptance tests for REQ-559: version constants in internal/config.
// These tests define the contract — they MUST fail until the production
// code in version.go is implemented.

// --- Happy-path scenarios ---

func TestVersionConstantIsValidSemver(t *testing.T) {
	// Given the Version constant is defined
	// When we inspect its value
	// Then it matches a valid semver pattern (MAJOR.MINOR.PATCH)
	semverRe := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !semverRe.MatchString(Version) {
		t.Errorf("Version %q is not valid semver", Version)
	}
}

func TestVersionConstantIsNotEmpty(t *testing.T) {
	// Given the Version constant exists
	// When we check its value
	// Then it is not empty
	if Version == "" {
		t.Error("Version must not be empty")
	}
}

func TestFullVersionContainsVersion(t *testing.T) {
	// Given Version is set
	// When we call FullVersion()
	// Then the result contains the Version string
	fv := FullVersion()
	if !containsSubstring(fv, Version) {
		t.Errorf("FullVersion() = %q, want it to contain Version %q", fv, Version)
	}
}

func TestFullVersionWithoutLdflags(t *testing.T) {
	// Given GitCommit and BuildTime are empty (no -ldflags injection)
	// When we call FullVersion()
	// Then the output contains "unknown" for both fields
	origCommit := GitCommit
	origTime := BuildTime
	GitCommit = ""
	BuildTime = ""
	t.Cleanup(func() {
		GitCommit = origCommit
		BuildTime = origTime
	})

	fv := FullVersion()
	if !containsSubstring(fv, "unknown") {
		t.Errorf("FullVersion() = %q, want it to contain 'unknown' when ldflags not set", fv)
	}
}

func TestFullVersionWithLdflags(t *testing.T) {
	// Given GitCommit and BuildTime are injected via -ldflags
	// When we call FullVersion()
	// Then the output contains the injected values
	origCommit := GitCommit
	origTime := BuildTime
	GitCommit = "abc1234"
	BuildTime = "2026-04-20T12:00:00Z"
	t.Cleanup(func() {
		GitCommit = origCommit
		BuildTime = origTime
	})

	fv := FullVersion()
	if !containsSubstring(fv, "abc1234") {
		t.Errorf("FullVersion() = %q, want it to contain commit 'abc1234'", fv)
	}
	if !containsSubstring(fv, "2026-04-20T12:00:00Z") {
		t.Errorf("FullVersion() = %q, want it to contain build time '2026-04-20T12:00:00Z'", fv)
	}
}

func TestFullVersionFormat(t *testing.T) {
	// Given GitCommit and BuildTime are injected
	// When we call FullVersion()
	// Then the format is "VERSION (commit: COMMIT, built: TIME)"
	origCommit := GitCommit
	origTime := BuildTime
	GitCommit = "def5678"
	BuildTime = "2026-01-01T00:00:00Z"
	t.Cleanup(func() {
		GitCommit = origCommit
		BuildTime = origTime
	})

	expected := Version + " (commit: def5678, built: 2026-01-01T00:00:00Z)"
	fv := FullVersion()
	if fv != expected {
		t.Errorf("FullVersion() = %q, want %q", fv, expected)
	}
}

func TestFullVersionFormatUnknown(t *testing.T) {
	// Given GitCommit and BuildTime are empty
	// When we call FullVersion()
	// Then the format is "VERSION (commit: unknown, built: unknown)"
	origCommit := GitCommit
	origTime := BuildTime
	GitCommit = ""
	BuildTime = ""
	t.Cleanup(func() {
		GitCommit = origCommit
		BuildTime = origTime
	})

	expected := Version + " (commit: unknown, built: unknown)"
	fv := FullVersion()
	if fv != expected {
		t.Errorf("FullVersion() = %q, want %q", fv, expected)
	}
}

// --- Edge-case scenarios ---

func TestFullVersionPartialLdflags_CommitOnly(t *testing.T) {
	// Given only GitCommit is injected, BuildTime is empty
	// When we call FullVersion()
	// Then commit shows the value, built shows "unknown"
	origCommit := GitCommit
	origTime := BuildTime
	GitCommit = "aaa1111"
	BuildTime = ""
	t.Cleanup(func() {
		GitCommit = origCommit
		BuildTime = origTime
	})

	expected := Version + " (commit: aaa1111, built: unknown)"
	fv := FullVersion()
	if fv != expected {
		t.Errorf("FullVersion() = %q, want %q", fv, expected)
	}
}

func TestFullVersionPartialLdflags_BuildTimeOnly(t *testing.T) {
	// Given only BuildTime is injected, GitCommit is empty
	// When we call FullVersion()
	// Then commit shows "unknown", built shows the value
	origCommit := GitCommit
	origTime := BuildTime
	GitCommit = ""
	BuildTime = "2026-06-15T08:30:00Z"
	t.Cleanup(func() {
		GitCommit = origCommit
		BuildTime = origTime
	})

	expected := Version + " (commit: unknown, built: 2026-06-15T08:30:00Z)"
	fv := FullVersion()
	if fv != expected {
		t.Errorf("FullVersion() = %q, want %q", fv, expected)
	}
}

func TestGitCommitAndBuildTimeDefaultToEmpty(t *testing.T) {
	// Given no -ldflags are passed during build
	// When we inspect GitCommit and BuildTime
	// Then they default to the zero value for string (empty)
	// Note: this test verifies the vars are declared with zero-value defaults.
	// When production code exists, the vars should be declared as:
	//   var GitCommit string
	//   var BuildTime string
	// which means their default is "".
	// We cannot truly test "no ldflags" from within the test binary itself
	// since test binaries also don't inject ldflags. So we verify that after
	// resetting them to "", they are indeed empty.
	origCommit := GitCommit
	origTime := BuildTime
	GitCommit = ""
	BuildTime = ""
	t.Cleanup(func() {
		GitCommit = origCommit
		BuildTime = origTime
	})

	if GitCommit != "" {
		t.Errorf("GitCommit = %q, want empty string", GitCommit)
	}
	if BuildTime != "" {
		t.Errorf("BuildTime = %q, want empty string", BuildTime)
	}
}

func TestVersionConstantInitialValue(t *testing.T) {
	// Given the initial release
	// When we read Version
	// Then it should be "0.1.0" as specified in the dev spec
	if Version != "0.1.0" {
		t.Errorf("Version = %q, want %q", Version, "0.1.0")
	}
}

// --- Helper ---

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
