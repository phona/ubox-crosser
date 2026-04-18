## Context

Ubox-crosser is a Go 1.23 SOCKS5 reverse proxy that communicates over raw TCP with a custom JSON message protocol. It has no HTTP API surface. The project uses cobra for CLI, logrus for logging, and follows a `cmd/` + `server/` + `client/` layout.

## Goals / Non-Goals

**Goals:**
- Add a lightweight HTTP server exposing `GET /api/echo?msg=X`
- Return JSON responses with appropriate content types and status codes
- Use only Go stdlib (`net/http`) — no external HTTP framework

**Non-Goals:**
- Authentication or authorization on the echo endpoint
- HTTPS/TLS termination (assume reverse proxy or internal network)
- Replacing the existing TCP protocol with HTTP
- Adding middleware, logging interceptors, or metrics to the HTTP server

## Decisions

### 1. Use Go stdlib `net/http` for routing and serving

**Rationale:** The echo endpoint is trivial — a single route with no path parameters. Adding a framework (Chi, Gin) would be over-engineering. Stdlib keeps the dependency tree unchanged.

**Alternative considered:** Chi router — rejected because a single endpoint doesn't justify the dependency.

### 2. Separate HTTP handler package

**Rationale:** Place the handler in `server/api/` to keep HTTP concerns isolated from the existing SOCKS5/TCP code. This follows the project's existing separation (`server/` for server-side logic).

### 3. JSON response format

**Rationale:** Return `{"message": "<msg>"}` with `Content-Type: application/json`. JSON is the standard for REST APIs and is trivially parseable by any client.

### 4. HTTP listen address via configuration

**Rationale:** The HTTP server port should be configurable (default `:8080`) to avoid port conflicts. Follow the existing config pattern in `models/config/`.

## Risks / Trade-offs

- **Port conflict** → Make the HTTP listen address configurable; document the default.
- **Lifecycle coupling** → The HTTP server should start/stop cleanly alongside the existing services. Use `context.Context` for graceful shutdown.
- **Scope creep** → This is intentionally minimal. Future endpoints should follow the same pattern but are out of scope.
