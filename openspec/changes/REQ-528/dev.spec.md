# DEV-SPEC: REQ-528 — Add Version History Section to README

## File Structure

### New Files

None.

### Modified Files

| File | Change |
|------|--------|
| `README.md` | Append `## Version History` section after `## Architecture` |

## Change Detail

### `README.md`

Append the following block **after** the `## Architecture` section (after line 10, the last line of the current file):

```markdown

## Version History

- **v0.0.1** — sisyphus pipeline test
```

**Constraints**:
- One blank line before the new `## Version History` heading (standard Markdown section separator).
- Each version entry is a bulleted list item: `- **vX.Y.Z** — description`.
- The version tag is bold (`**v0.0.1**`).
- Use an em-dash (`—`) between version and description, not a hyphen.
- Do **not** modify any existing content in the file.

## Function Signatures & Responsibilities

N/A — documentation-only change, no code.

## Dependencies

None. No new packages, tools, or build changes.

## Error Handling Strategy

N/A — no runtime behavior.

## Boundary Conditions

1. **Markdown rendering**: The `## Version History` heading must be a level-2 heading to match existing sections (`## What is Ubox-crosser?`, `## Architecture`).
2. **Trailing newline**: File must end with a single trailing newline (Unix convention, matches existing file).
3. **No duplicate sections**: If a `## Version` or `## Version History` section already exists (e.g., from a parallel branch), the dev agent should detect and skip rather than duplicate.

## Storage / DB Requirements

None.

## Checklist for Dev Agent

- [ ] Append `## Version History` section to `README.md` after `## Architecture`
- [ ] Verify the section contains `- **v0.0.1** — sisyphus pipeline test`
- [ ] Verify no existing README content was removed or altered
- [ ] Verify file ends with a single trailing newline
- [ ] Verify Markdown renders correctly (heading level, bold, bullet)
