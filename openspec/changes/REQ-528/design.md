# Design: REQ-528 — Version History in README

## Decision: Inline in README vs Separate File

**Choice**: Inline `## Version History` section in README.md.

**Rationale**: The project is small, has only 3 git tags, and no existing changelog infrastructure. A separate CHANGELOG.md adds ceremony without value at this scale. The README is the single entry point for project documentation.

## Section Placement

After `## Architecture`, before any future sections. This follows a natural reading order: what → how → history.

## Format

Markdown bulleted list with version and description:

```markdown
## Version History

- **v0.0.1** — sisyphus pipeline test
```

**Rationale**: Bullet list is scannable, easy to maintain, and doesn't require table formatting overhead. Bold version tag aids visual scanning.
