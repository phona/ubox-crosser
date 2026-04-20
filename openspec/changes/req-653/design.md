---
change_id: req-653
title: "Design: GET /version endpoint"
---

## Context

ubox-crosser is a proxy tunnel server. It currently exposes only a TCP listener for tunnel traffic. There is no HTTP endpoint to query runtime metadata such as version, commit hash, or build timestamp. This makes deployment verification and troubleshooting difficult across distributed nodes.

## Goals / Non-Goals

**Goals:**
- Expose build-time metadata (version, commit, build_time) via a lightweight HTTP endpoint
- Keep the implementation minimal with zero external dependencies
- Inject metadata at compile time via standard Go ldflags

**Non-Goals:**
- Health check endpoint (separate REQ scope)
- Authentication or rate-limiting on the version endpoint
- Metrics or observability endpoints
- Shared HTTP mux with other future endpoints

## Decisions

### 1. Package location: `version/` at repo root

**Chosen:** Top-level `version/` package.
**Alternative:** `internal/version/`. Rejected because version metadata may be consumed by other cmd binaries (client, auth_server), so it should not be internal-scoped.

### 2. HTTP server embedding in cmd/server

**Chosen:** Embed an `http.ServeMux` in `cmd/server` alongside the existing TCP listener, controlled by `--http-addr` flag (default `:8080`).
**Alternative:** Separate binary/service for HTTP endpoints. Rejected — overkill for a single endpoint.

### 3. Version source: constant + ldflags

**Chosen:** `const Version = "0.1.0"` in source, `Commit` and `BuildTime` injected via `go build -ldflags -X`. Standard Go pattern, zero runtime cost, no file I/O at startup.
**Alternative:** Read from embedded file or env var. More flexible but adds failure modes.

### 4. Method restriction in handler

**Chosen:** Handler-level `r.Method != http.MethodGet` check returning 405. No router framework needed for a single endpoint.
**Alternative:** Use Chi or gorilla/mux for method routing. Rejected — unnecessary dependency for one route.

### 5. JSON serialization

**Chosen:** `encoding/json.NewEncoder(w).Encode()` with struct tags matching the contract field names (`version`, `commit`, `build_time`).

## Risks / Trade-offs

- **[Port conflict]** `--http-addr :8080` may conflict with other services → Mitigation: flag is configurable, can be set to empty string to disable HTTP server entirely
- **[No graceful shutdown]** HTTP server goroutine has no graceful shutdown → Mitigation: acceptable for a metadata-only endpoint; can be addressed in future if more endpoints are added
- **[REQ-601 overlap]** A future health endpoint REQ may need a shared HTTP mux → Mitigation: current design uses `http.ServeMux` which is extensible; migration to shared mux is straightforward
