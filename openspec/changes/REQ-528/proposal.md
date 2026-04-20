# REQ-528: Add Version History Section to README

## What

Add a "Version History" section to `README.md` that records each released version with a one-line description. The initial entry is `v0.0.1 — sisyphus pipeline test`.

## Why

The project has no visible changelog or release history. Contributors and users cannot tell what changed between versions by reading the README. A lightweight version history section in the README gives immediate visibility without introducing a separate CHANGELOG file.

## Scope

- Documentation only — no code changes.
- Single file modified: `README.md`.
- Add a new `## Version History` section after the existing `## Architecture` section.
- Seed with the initial version entry.

## Out of Scope

- Automated changelog generation.
- Separate CHANGELOG.md file.
- Linking versions to git tags.
