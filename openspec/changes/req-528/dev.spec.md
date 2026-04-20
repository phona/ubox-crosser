# DEV-SPEC: REQ-528 — Add Version History Section to README

## File Structure

### New Files

None.

### Modified Files

| File | Change |
|------|--------|
| `README.md` | Append `## Version History` section with Markdown table at end of file |

## Change Detail

### `README.md`

Append the following block **after** the `## Architecture` section (after line 10, the last line of the current file). The `## Version History` heading must be the **last** H2 section in the file.

```markdown

## Version History

| Version | Date | Changes |
|---------|------|---------|
| v0.0.1 | 2026-04-20 | Initial sisyphus pipeline test |
```

**Format constraints** (enforced by acceptance tests in `tests/acceptance/readme_version_history_test.go`):

1. **Heading**: Exactly `## Version History` (level-2, matches existing sections).
2. **Table format**: Markdown table with three columns — `Version`, `Date`, `Changes`.
3. **Version column**: Semantic version `vX.Y.Z` (regex: `^v\d+\.\d+\.\d+$`).
4. **Date column**: ISO 8601 date `YYYY-MM-DD` (regex: `^\d{4}-\d{2}-\d{2}$`).
5. **Changes column**: Non-empty text describing the change.
6. **Ordering**: Reverse chronological — most recent version first (tested when ≥2 entries).
7. **Section position**: Must be the last `## ` heading in the file.
8. **One blank line** before the heading (standard Markdown section separator).
9. **Trailing newline**: File must end with a single `\n`.

### Do NOT

- Modify any existing README content.
- Add entries beyond `v0.0.1` (the initial seed entry).
- Use bullet list format (tests expect a table).

## Function Signatures & Responsibilities

N/A — documentation-only change, no code.

## Dependencies

None. No new packages, tools, or build changes.

## Error Handling Strategy

N/A — no runtime behavior.

## Boundary Conditions

1. **Duplicate detection**: If `## Version History` already exists (e.g., from a parallel branch merge), the dev agent must detect and skip rather than duplicate.
2. **Table structure**: The acceptance tests parse table rows by splitting on `|` and trimming whitespace. Ensure table cells have no leading/trailing pipes issues.
3. **Separator row**: The table must include the standard `|---|---|---|` separator between header and data rows. The acceptance tests skip rows containing `---`.

## Storage / DB Requirements

None.

## Acceptance Tests (Already Written)

Tests exist at `tests/acceptance/readme_version_history_test.go` (build tag: `acceptance`). They verify:

| Test | What it checks |
|------|---------------|
| `TestVersionHistorySectionExists` | `## Version History` heading present |
| `TestVersionHistorySectionIsLast` | Heading is the last H2 in the file |
| `TestVersionHistoryTableHeaderPresent` | Table header contains `Version`, `Date`, `Changes` columns |
| `TestVersionHistoryHasAtLeastOneEntry` | At least one data row in the table |
| `TestVersionEntryFormat` | Version matches `vX.Y.Z`, date matches `YYYY-MM-DD`, changes non-empty |
| `TestVersionEntriesDescendingOrder` | Versions descend top-to-bottom (skipped if <2 entries) |
| `TestRejectEmptyChanges` | No entry has empty Changes column |

Run with: `go test -tags acceptance ./tests/acceptance/ -run VersionHistory -v`

## Checklist for Dev Agent

- [ ] Append `## Version History` section with table to end of `README.md`
- [ ] Table has header row: `| Version | Date | Changes |`
- [ ] Table has separator row: `|---------|------|---------|`
- [ ] Table has one data row: `| v0.0.1 | 2026-04-20 | Initial sisyphus pipeline test |`
- [ ] No existing README content is removed or altered
- [ ] `## Version History` is the last H2 heading
- [ ] File ends with a single trailing newline
- [ ] All acceptance tests pass: `go test -tags acceptance ./tests/acceptance/ -run VersionHistory -v`
