---
layers:
  - backend
---

## Why

ubox-crosser currently has no way to report its build version at runtime. Operators and CI pipelines need a lightweight health-check endpoint that also surfaces the exact binary revision, enabling fast rollback triage and deployment verification. This is the first HTTP endpoint in the project; the existing transport is pure TCP.

## What Changes

- Add a new `internal/version` package exposing `Version`, `Commit`, and `BuildTime` variables (injected via `-ldflags` at compile time).
- Introduce a minimal `net/http` server bound to a configurable HTTP port on the proxy-server binary, serving `GET /version`.
- Update the `Makefile` build target to inject git SHA and ISO 8601 timestamp via `-ldflags -X`.
- Add unit tests for the `/version` handler (200, correct Content-Type, all three JSON fields present).
- Add a contract test validating the response against the OpenAPI schema.

## Capabilities

### New Capabilities
- `version-endpoint`: HTTP GET /version returning JSON with version, commit, and build_time fields. No authentication required.

### Modified Capabilities

_(none)_

## Impact

- **Code**: New `internal/version/` package; new HTTP listener in `cmd/server`.
- **Build**: Makefile ldflags change (non-breaking, additive).
- **Dependencies**: stdlib only (`net/http`, `encoding/json`). No new third-party deps.
- **Operations**: A new TCP port will be opened for HTTP; must be documented and configurable.
