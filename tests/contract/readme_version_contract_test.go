//go:build contract

package contract

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func repoRoot() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..")
}

func readmeLines(t *testing.T) []string {
	t.Helper()
	path := filepath.Join(repoRoot(), "README.md")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("cannot open README.md: %v", err)
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("error reading README.md: %v", err)
	}
	return lines
}

func findSectionStart(lines []string, heading string) int {
	for i, l := range lines {
		if strings.TrimSpace(l) == heading {
			return i
		}
	}
	return -1
}

func sectionBody(lines []string, start int) []string {
	var body []string
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			break
		}
		body = append(body, lines[i])
	}
	return body
}

func nonEmptyLines(lines []string) []string {
	var out []string
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}

var versionEntryRe = regexp.MustCompile(`^- v\d+\.\d+\.\d+\s+-\s+.+$`)

func TestReadmeContainsVersionSection(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionStart(lines, "## Version")
	if idx == -1 {
		t.Fatal("README.md must contain a '## Version' section")
	}
}

func TestVersionSectionHasAtLeastOneEntry(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionStart(lines, "## Version")
	if idx == -1 {
		t.Fatal("README.md must contain a '## Version' section")
	}
	body := nonEmptyLines(sectionBody(lines, idx))
	if len(body) == 0 {
		t.Fatal("Version section must contain at least one entry")
	}
}

func TestVersionEntryFormat(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionStart(lines, "## Version")
	if idx == -1 {
		t.Fatal("README.md must contain a '## Version' section")
	}
	body := nonEmptyLines(sectionBody(lines, idx))
	if len(body) == 0 {
		t.Fatal("Version section must contain at least one entry")
	}
	for _, entry := range body {
		if !versionEntryRe.MatchString(entry) {
			t.Errorf("version entry does not match format '- vX.Y.Z - description': %q", entry)
		}
	}
}

func TestVersionSectionIsLastSection(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionStart(lines, "## Version")
	if idx == -1 {
		t.Fatal("README.md must contain a '## Version' section")
	}
	for i := idx + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			t.Fatal("Version section must be the last section in README.md")
		}
	}
}

func TestVersionSectionAppearsAfterArchitecture(t *testing.T) {
	lines := readmeLines(t)
	archIdx := findSectionStart(lines, "## Architecture")
	verIdx := findSectionStart(lines, "## Version")
	if verIdx == -1 {
		t.Fatal("README.md must contain a '## Version' section")
	}
	if archIdx == -1 {
		t.Fatal("README.md must contain a '## Architecture' section")
	}
	if verIdx <= archIdx {
		t.Fatal("Version section must appear after Architecture section")
	}
}

func TestVersionSectionContentType(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionStart(lines, "## Version")
	if idx == -1 {
		t.Fatal("README.md must contain a '## Version' section")
	}
	body := nonEmptyLines(sectionBody(lines, idx))
	for _, entry := range body {
		if !strings.HasPrefix(entry, "- ") {
			t.Errorf("version entries must be unordered list items (start with '- '): %q", entry)
		}
	}
}

func TestVersionEntrySemver(t *testing.T) {
	semverInEntry := regexp.MustCompile(`v(\d+)\.(\d+)\.(\d+)`)
	lines := readmeLines(t)
	idx := findSectionStart(lines, "## Version")
	if idx == -1 {
		t.Fatal("README.md must contain a '## Version' section")
	}
	body := nonEmptyLines(sectionBody(lines, idx))
	if len(body) == 0 {
		t.Fatal("Version section must contain at least one entry")
	}
	for _, entry := range body {
		if !semverInEntry.MatchString(entry) {
			t.Errorf("version entry must contain a semver (vX.Y.Z): %q", entry)
		}
	}
}

func TestVersionEntryHasDescription(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionStart(lines, "## Version")
	if idx == -1 {
		t.Fatal("README.md must contain a '## Version' section")
	}
	body := nonEmptyLines(sectionBody(lines, idx))
	if len(body) == 0 {
		t.Fatal("Version section must contain at least one entry")
	}
	for _, entry := range body {
		parts := strings.SplitN(entry, " - ", 3)
		if len(parts) < 3 {
			t.Errorf("version entry must have format '- vX.Y.Z - description': %q", entry)
			continue
		}
		desc := strings.TrimSpace(parts[2])
		if desc == "" {
			t.Errorf("version entry description must not be empty: %q", entry)
		}
	}
}

func TestReadmeFileExists(t *testing.T) {
	path := filepath.Join(repoRoot(), "README.md")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("README.md must exist: %v", err)
	}
	if info.IsDir() {
		t.Fatal("README.md must be a file, not a directory")
	}
}

func TestReadmeIsNotEmpty(t *testing.T) {
	path := filepath.Join(repoRoot(), "README.md")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("README.md must exist: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("README.md must not be empty")
	}
}
