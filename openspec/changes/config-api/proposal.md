## Why

ubox-crosser currently loads configuration from JSON files at startup with no way to inspect or modify the running configuration without restarting the process. Operators need runtime visibility into the active configuration and the ability to adjust settings (e.g., log level) without service interruption. An HTTP-based config API enables this and lays the groundwork for future management endpoints.

## What Changes

- Add an embedded HTTP management server to the proxy server process, listening on a configurable address (e.g., `127.0.0.1:8080`).
- Expose `GET /api/config` to return the current running configuration as JSON.
- Expose `PUT /api/config` to update mutable configuration fields at runtime without restart.
- Only a subset of fields are hot-updatable (e.g., `log_level`). Immutable fields (e.g., `address`, `key`, `method`) are rejected with an error if modification is attempted.

## Capabilities

### New Capabilities
- `config-read`: HTTP endpoint to retrieve the current running configuration via `GET /api/config`.
- `config-update`: HTTP endpoint to update mutable configuration fields via `PUT /api/config` with hot-reload (no restart).

### Modified Capabilities
<!-- No existing specs to modify -->

## Impact

- **Code**: New `api/` package for HTTP handler and server. Changes to `cmd/server/server.go` to start the management HTTP server alongside the proxy. Changes to `models/config/config.go` to support concurrency-safe reads and mutable field identification.
- **APIs**: New HTTP REST API surface (`/api/config`). The existing TCP protocol is unchanged.
- **Dependencies**: Uses Go stdlib `net/http` only — no new external dependencies.
- **Systems**: The management HTTP port must be firewall-restricted to localhost/internal networks to prevent unauthorized configuration changes.
