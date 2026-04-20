---
layers:
  - backend
---

## Why

ubox-crosser currently has no way to report its build version at runtime. Operators and CI pipelines need a lightweight health-check endpoint that also surfaces the exact binary revision, enabling fast rollback triage and deployment verification. This is the first HTTP endpoint in the project; the existing transport is pure TCP.

## What Changes

- Add a new `internal/version` package exposing `Version`, `Commit`, and `BuildTime` variables (injected via `-ldflags` at compile time).
- Register `GET /version` on REQ-601's shared health HTTP mux (`--health-address`), no new listener or flag.
- Update the `Makefile` build target to inject git SHA and ISO 8601 timestamp via `-ldflags -X`.
- Add unit tests for the `/version` handler (200, correct Content-Type, all three JSON fields present).
- Add contract tests using `net/http/httptest` with hand-written JSON assertions (no third-party validator).

## Capabilities

### New Capabilities
- `version-endpoint`: HTTP GET /version returning JSON with version, commit, and build_time fields. No authentication required.

### Modified Capabilities

_(none)_

## Impact

- **Code**: New `internal/version/` package; route registration in `server/health.go` (REQ-601's mux).
- **Build**: Makefile ldflags change (non-breaking, additive).
- **Dependencies**: stdlib only (`net/http`, `encoding/json`). No new third-party deps.
- **Operations**: No new port — shares REQ-601's `--health-address` listener.
