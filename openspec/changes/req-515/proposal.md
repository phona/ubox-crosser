---
layers:
  - api
  - server
---

## Why

The system currently has no HTTP API surface — all communication is via custom JSON-over-TCP messaging. Operators and monitoring tools need a simple HTTP health/version endpoint to verify which release is running. A lightweight `GET /version` endpoint is the first step toward an HTTP management plane.

## What Changes

- Add a new HTTP server that listens on a configurable port alongside the existing TCP listeners.
- Expose a single `GET /version` endpoint that returns the hardcoded version string `v3` as a JSON response.

## Capabilities

### New Capabilities
- `version-endpoint`: HTTP `GET /version` endpoint that returns the application version in JSON format.

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- **Code**: New HTTP handler and server startup logic added to the `server` package (or a new `api` package); `cmd/server/server.go` wired up to start the HTTP listener.
- **Dependencies**: Uses only Go stdlib `net/http` — no new third-party dependencies.
- **APIs**: Introduces the project's first HTTP endpoint (`GET /version`).
- **Deployment**: A new port must be exposed; defaults should be documented.
