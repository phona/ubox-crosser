//go:build acceptance

package acceptance

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func repoRoot() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..")
}

func readmeLines(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(), "README.md"))
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func findSectionIndex(lines []string, heading string) int {
	for i, l := range lines {
		if strings.TrimSpace(l) == heading {
			return i
		}
	}
	return -1
}

func allH2Indices(lines []string) []int {
	var indices []int
	for i, l := range lines {
		if strings.HasPrefix(l, "## ") {
			indices = append(indices, i)
		}
	}
	return indices
}

func versionHistoryTableRows(t *testing.T, lines []string, sectionIdx int) []string {
	t.Helper()
	var rows []string
	inSection := false
	pastHeader := false
	for i := sectionIdx + 1; i < len(lines); i++ {
		l := lines[i]
		if strings.HasPrefix(l, "## ") {
			break
		}
		if strings.HasPrefix(l, "|") {
			inSection = true
			if !pastHeader {
				pastHeader = true
				continue // header row
			}
			if strings.Contains(l, "---") {
				continue // separator row
			}
			rows = append(rows, l)
		} else if inSection && strings.TrimSpace(l) == "" {
			continue
		}
	}
	return rows
}

func TestVersionHistorySectionExists(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionIndex(lines, "## Version History")
	if idx == -1 {
		t.Fatal("README.md does not contain '## Version History' section")
	}
}

func TestVersionHistorySectionIsLast(t *testing.T) {
	lines := readmeLines(t)
	h2s := allH2Indices(lines)
	if len(h2s) == 0 {
		t.Fatal("no H2 headings found in README.md")
	}
	lastH2 := lines[h2s[len(h2s)-1]]
	if strings.TrimSpace(lastH2) != "## Version History" {
		t.Fatalf("expected '## Version History' to be the last H2 section, got %q", lastH2)
	}
}

func TestVersionHistoryTableHeaderPresent(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionIndex(lines, "## Version History")
	if idx == -1 {
		t.Fatal("section not found")
	}
	found := false
	for i := idx + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			break
		}
		normalized := strings.ReplaceAll(lines[i], " ", "")
		if strings.Contains(normalized, "|Version|Date|Changes|") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Version History table header '| Version | Date | Changes |' not found")
	}
}

func TestVersionHistoryHasAtLeastOneEntry(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionIndex(lines, "## Version History")
	if idx == -1 {
		t.Fatal("section not found")
	}
	rows := versionHistoryTableRows(t, lines, idx)
	if len(rows) == 0 {
		t.Fatal("Version History table has no data rows")
	}
}

var semverRe = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func TestVersionEntryFormat(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionIndex(lines, "## Version History")
	if idx == -1 {
		t.Fatal("section not found")
	}
	rows := versionHistoryTableRows(t, lines, idx)
	if len(rows) == 0 {
		t.Fatal("no version entries to validate")
	}
	for _, row := range rows {
		cols := splitTableRow(row)
		if len(cols) < 3 {
			t.Fatalf("expected at least 3 columns, got %d in row: %s", len(cols), row)
		}
		version := strings.TrimSpace(cols[0])
		date := strings.TrimSpace(cols[1])
		changes := strings.TrimSpace(cols[2])

		if !semverRe.MatchString(version) {
			t.Errorf("version %q does not match vX.Y.Z format", version)
		}
		if !dateRe.MatchString(date) {
			t.Errorf("date %q does not match YYYY-MM-DD format", date)
		}
		if changes == "" {
			t.Errorf("changes column is empty for version %s", version)
		}
	}
}

func TestVersionEntriesDescendingOrder(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionIndex(lines, "## Version History")
	if idx == -1 {
		t.Fatal("section not found")
	}
	rows := versionHistoryTableRows(t, lines, idx)
	if len(rows) < 2 {
		t.Skip("need at least 2 entries to verify ordering")
	}
	for i := 0; i < len(rows)-1; i++ {
		cols1 := splitTableRow(rows[i])
		cols2 := splitTableRow(rows[i+1])
		v1 := strings.TrimSpace(cols1[0])
		v2 := strings.TrimSpace(cols2[0])
		if compareVersions(v1, v2) <= 0 {
			t.Errorf("versions not in descending order: %s should be greater than %s", v1, v2)
		}
	}
}

func TestRejectEmptyChanges(t *testing.T) {
	lines := readmeLines(t)
	idx := findSectionIndex(lines, "## Version History")
	if idx == -1 {
		t.Fatal("section not found")
	}
	rows := versionHistoryTableRows(t, lines, idx)
	for _, row := range rows {
		cols := splitTableRow(row)
		if len(cols) >= 3 && strings.TrimSpace(cols[2]) == "" {
			t.Errorf("empty changes column found in row: %s", row)
		}
	}
}

func splitTableRow(row string) []string {
	row = strings.TrimSpace(row)
	row = strings.Trim(row, "|")
	parts := strings.Split(row, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func compareVersions(a, b string) int {
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")
	pa := strings.Split(a, ".")
	pb := strings.Split(b, ".")
	for i := 0; i < 3; i++ {
		na, _ := strconv.Atoi(pa[i])
		nb, _ := strconv.Atoi(pb[i])
		if na != nb {
			return na - nb
		}
	}
	return 0
}
