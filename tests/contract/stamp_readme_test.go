package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var timestampPattern = regexp.MustCompile(`^Built at: (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)$`)

func projectRoot(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Fatalf("cannot find project root: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func lastLine(content string) string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	if len(lines) == 0 {
		return ""
	}
	return lines[len(lines)-1]
}

func countBuiltAtLines(content string) int {
	count := 0
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "Built at:") {
			count++
		}
	}
	return count
}

// REQ-975-S1: make stamp-readme appends a valid Built at timestamp line to README.md
func TestStampReadmeAppendsTimestamp(t *testing.T) {
	root := projectRoot(t)
	readmePath := filepath.Join(root, "README.md")

	orig, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("README.md must exist for this test: %v", err)
	}
	t.Cleanup(func() { _ = os.WriteFile(readmePath, orig, 0644) })

	before := time.Now().UTC().Add(-2 * time.Second)
	cmd := exec.Command("make", "stamp-readme")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make stamp-readme failed: %v\noutput: %s", err, out)
	}
	after := time.Now().UTC().Add(2 * time.Second)

	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read README.md after stamp: %v", err)
	}

	last := lastLine(string(content))
	matches := timestampPattern.FindStringSubmatch(last)
	if matches == nil {
		t.Fatalf("last line does not match expected pattern.\ngot: %q\nwant: Built at: YYYY-MM-DDTHH:MM:SSZ", last)
	}

	ts, err := time.Parse(time.RFC3339, matches[1])
	if err != nil {
		t.Fatalf("timestamp is not valid ISO 8601: %v", err)
	}

	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v is not within expected range [%v, %v]", ts, before, after)
	}
}

// REQ-975-S2: running stamp-readme twice yields exactly one Built at line
func TestStampReadmeIdempotent(t *testing.T) {
	root := projectRoot(t)
	readmePath := filepath.Join(root, "README.md")

	orig, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("README.md must exist for this test: %v", err)
	}
	t.Cleanup(func() { _ = os.WriteFile(readmePath, orig, 0644) })

	for i := 0; i < 2; i++ {
		cmd := exec.Command("make", "stamp-readme")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("make stamp-readme (run %d) failed: %v\noutput: %s", i+1, err, out)
		}
	}

	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}

	count := countBuiltAtLines(string(content))
	if count != 1 {
		t.Errorf("expected exactly 1 Built at line, found %d", count)
	}
}

// REQ-975-S3: stamp-readme exits 0 when README.md does not exist
func TestStampReadmeNoReadme(t *testing.T) {
	root := projectRoot(t)
	readmePath := filepath.Join(root, "README.md")

	orig, err := os.ReadFile(readmePath)
	if err != nil {
		t.Skipf("README.md not found, cannot test removal scenario: %v", err)
	}

	if err := os.Remove(readmePath); err != nil {
		t.Fatalf("failed to remove README.md: %v", err)
	}
	t.Cleanup(func() { _ = os.WriteFile(readmePath, orig, 0644) })

	cmd := exec.Command("make", "stamp-readme")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make stamp-readme should exit 0 when README.md absent, got error: %v\noutput: %s", err, out)
	}

	if _, err := os.Stat(readmePath); err == nil {
		t.Error("README.md should not be created when it was absent")
	}
}
