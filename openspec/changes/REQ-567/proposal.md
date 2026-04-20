## Why

Kubernetes liveness/readiness probes and container orchestrators expect an HTTP health endpoint to determine whether the proxy server process is alive. The conventional Kubernetes path is `/healthz`. Without it, operators rely on raw TCP port probes that cannot distinguish a listening socket from a functioning service.

## What Changes

- Add a new HTTP listener to the proxy server dedicated to health checking.
- Expose a `GET /healthz` endpoint that returns HTTP 200 with `{"status":"ok"}` and `Content-Type: application/json`.
- Add a configurable `health_address` field to `ServerConfig` for specifying the HTTP listen address (host:port).
- Add a `--health-address` CLI flag to the server command.
- Non-GET requests to `/healthz` return HTTP 405 Method Not Allowed with `Allow: GET` header.
- Requests to any path other than `/healthz` return HTTP 404 Not Found.

## Capabilities

### New Capabilities

- `healthz-endpoint`: HTTP health-check endpoint on the proxy server returning liveness status at `/healthz`.

### Modified Capabilities

(none)

## Impact

- **Code**: `server/` package gains a new HTTP server; `models/config/` adds `health_address` field; `cmd/server/` wires the new flag and starts the health server.
- **APIs**: New HTTP API surface (`GET /healthz`) on a separate port from existing TCP listeners.
- **Dependencies**: None new — uses Go stdlib `net/http`.
- **Systems**: Deployment manifests and Docker/Kubernetes health-check configs can now target the HTTP endpoint.
