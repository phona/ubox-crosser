package contract

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var timestampRe = regexp.MustCompile(`^Built at: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

func stampReadme(dir string, timestamp string) error {
	path := filepath.Join(dir, "README.md")

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "Built at: ") {
			lines = append(lines, line)
		}
	}

	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	lines = append(lines, "Built at: "+timestamp, "")
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

// REQ-983-S1: make stamp-readme appends Built-at timestamp to README.md
func TestStampReadmeAppendsTimestamp(t *testing.T) {
	dir := t.TempDir()
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# My Project\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ts := "2026-04-21T12:00:00Z"
	if err := stampReadme(dir, ts); err != nil {
		t.Fatalf("stampReadme failed: %v", err)
	}

	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	builtLine := "Built at: " + ts
	if !strings.Contains(content, builtLine) {
		t.Fatalf("README should contain %q, got:\n%s", builtLine, content)
	}

	if !timestampRe.MatchString(builtLine) {
		t.Fatalf("timestamp line %q does not match ISO-8601 UTC format", builtLine)
	}

	count := strings.Count(content, "Built at: ")
	if count != 1 {
		t.Fatalf("expected exactly 1 Built at line, found %d", count)
	}
}

// REQ-983-S2: stamp-readme is idempotent (no duplicate lines)
func TestStampReadmeIdempotent(t *testing.T) {
	dir := t.TempDir()
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# My Project\nBuilt at: 2026-04-20T10:00:00Z\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ts := "2026-04-21T12:00:00Z"
	if err := stampReadme(dir, ts); err != nil {
		t.Fatalf("stampReadme failed: %v", err)
	}

	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	count := strings.Count(content, "Built at: ")
	if count != 1 {
		t.Fatalf("expected exactly 1 Built at line after re-run, found %d", count)
	}

	if !strings.Contains(content, "Built at: "+ts) {
		t.Fatalf("expected new timestamp %q, got:\n%s", ts, content)
	}

	if strings.Contains(content, "Built at: 2026-04-20T10:00:00Z") {
		t.Fatal("old timestamp should have been removed")
	}
}

// REQ-983-S3: stamp-readme exits 0 when README.md does not exist
func TestStampReadmeNoReadme(t *testing.T) {
	dir := t.TempDir()

	if err := stampReadme(dir, "2026-04-21T12:00:00Z"); err != nil {
		t.Fatalf("stampReadme should succeed when README.md is missing, got: %v", err)
	}

	readme := filepath.Join(dir, "README.md")
	if _, err := os.Stat(readme); !os.IsNotExist(err) {
		t.Fatal("README.md should not be created when it did not exist")
	}
}
