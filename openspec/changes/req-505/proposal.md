---
layers:
  - api
  - config
  - build
---

## Why

The ubox-crosser server currently exposes no HTTP endpoints, making it difficult for operations tooling (health checks, monitoring dashboards, deployment verification) to query the running service version. A simple `GET /version` HTTP endpoint enables load balancers, CI/CD pipelines, and ops scripts to confirm which build is deployed without inspecting binaries or logs.

## What Changes

- Introduce a lightweight HTTP management API on the proxy server that listens on a configurable address.
- Add a `GET /version` endpoint that returns the application version string (`v2`) in a JSON response.
- Inject the version at build time via `-ldflags` so the binary carries the correct version without source changes per release.

## Capabilities

### New Capabilities
- `version-endpoint`: HTTP endpoint that reports the running application version.

### Modified Capabilities

_(none — no existing specs to modify)_

## Impact

- **Code**: New HTTP listener setup in `server` package; new `version` package or constant for build-time injection.
- **Build**: Makefile `-ldflags` updated to embed version string.
- **Config**: Optional `management-address` field in `ServerConfig` for the HTTP listen address.
- **Dependencies**: Uses Go stdlib `net/http` only — no new external dependencies.
- **Deployment**: An additional port is exposed; documentation/Docker config may need updating.
