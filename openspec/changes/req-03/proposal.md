## Why

ubox-crosser currently exposes only raw TCP listeners for its proxy protocol — there is no lightweight way for orchestrators (Docker, Kubernetes, load balancers) to check whether the server process is alive and ready. A standard HTTP health endpoint removes the need for custom TCP probes and integrates with every major container-health and service-mesh system out of the box.

## What Changes

- Introduce a minimal HTTP listener on the proxy server that serves a single `GET /health` endpoint.
- The endpoint returns HTTP 200 with a JSON body `{"status":"ok"}` when the server is running.
- The HTTP listen address is configurable (default `:8080`) so it does not collide with the existing TCP control/data ports.

## Capabilities

### New Capabilities
- `health-endpoint`: HTTP health-check endpoint (`GET /health`) on the proxy server, returning 200 / `{"status":"ok"}`.

### Modified Capabilities
<!-- None — this is a purely additive feature with no changes to existing proxy behaviour. -->

## Impact

- **Code**: new HTTP server startup logic in the `server` package (or a thin wrapper in `cmd/server`); new config field for the health listen address.
- **APIs**: one new HTTP endpoint (`GET /health`). No changes to the existing TCP-based proxy protocol.
- **Dependencies**: uses Go stdlib `net/http` only — no new third-party dependencies.
- **Systems**: deployment manifests / Docker HEALTHCHECK directives can be updated to point at `/health`.
