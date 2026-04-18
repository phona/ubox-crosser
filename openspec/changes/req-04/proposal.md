## Why

The ubox-crosser proxy server exposes only custom TCP protocol listeners. Container orchestrators (Kubernetes, Docker) and load balancers require an HTTP endpoint to perform health checks. Without one, operators must rely on raw TCP port probes, which cannot distinguish a listening socket from a functioning service. Adding a standard HTTP health endpoint enables proper liveness monitoring.

## What Changes

- Add a new HTTP listener to the proxy server dedicated to health checking.
- Expose a `GET /health` endpoint that returns HTTP 200 with `{"status":"ok"}`.
- Add a configurable `health_address` field to `ServerConfig` for specifying the HTTP listen address (host:port).
- Add a `--health-address` CLI flag to the server command.
- Non-GET requests to `/health` return HTTP 405 Method Not Allowed.
- Requests to any path other than `/health` return HTTP 404 Not Found.

## Capabilities

### New Capabilities

- `health-endpoint`: HTTP health-check endpoint on the proxy server returning liveness status.

### Modified Capabilities

(none)

## Impact

- **Code**: `server/` package gains a new HTTP server; `models/config/` adds `health_address` field; `cmd/server/` wires the new flag and starts the health server.
- **APIs**: New HTTP API surface (`GET /health`) on a separate port from existing TCP listeners.
- **Dependencies**: None new — uses Go stdlib `net/http`.
- **Systems**: Deployment manifests and Docker health-check configs can now target the HTTP endpoint.
