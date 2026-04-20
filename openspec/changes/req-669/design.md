---
change_id: req-669
title: "GET /version endpoint (v4) — Design"
---

## Context

ubox-crosser is a TCP-based proxy tunnel using a custom binary protocol. There is no built-in HTTP surface. Operators need a lightweight way to query the running build version, commit, and build timestamp without SSH access to verify deployments across proxy nodes.

The `version` package and admin HTTP listener already exist from REQ-653. This change formalizes the v4 specification for the same capability.

## Goals / Non-Goals

**Goals:**
- Expose `GET /version` on a dedicated admin HTTP listener returning JSON build metadata
- Inject `Commit` and `BuildTime` at compile time via `-ldflags -X`
- Keep the admin server decoupled from the tunnel protocol
- Accept only GET; reject other methods with 405
- Return 404 for unknown paths

**Non-Goals:**
- Health-check, readiness, or metrics endpoints (future work)
- Authentication or TLS on the admin listener
- Dynamic version bumping or runtime-mutable metadata

## Decisions

### Decision 1: Separate HTTP admin server on its own port

The tunnel uses raw TCP with a custom binary protocol. Mixing HTTP into the TCP handler would break protocol compatibility and couple unrelated concerns.

**Chosen:** Standalone `net/http.ServeMux` on `--admin-addr` (default `:8080`).

| Option | Pros | Cons |
|--------|------|------|
| Embed in TCP protocol | No extra port | Breaks binary protocol, clients must parse HTTP |
| Separate HTTP server | Clean separation, standard tooling | Extra port to expose |

### Decision 2: Go 1.22+ method-based mux routing

Register `"GET /version"` on `http.ServeMux` so the stdlib returns 405 for non-GET methods automatically — no manual method checking required.

### Decision 3: Build-time injection via ldflags

Use `go build -ldflags "-X pkg.Var=val"` for `Commit` and `BuildTime`. The `Version` constant is hardcoded (`0.1.0`) and changes only with code, not builds. Default values are `"unknown"` when ldflags are omitted.

### Decision 4: Makefile and Dockerfile wire ldflags

`Makefile` computes `GIT_COMMIT` and `BUILD_TIME` from shell commands and passes them as `-X` flags. `Dockerfile` forwards build args for the same purpose.

## Risks / Trade-offs

- **[Unprotected admin port]** → The admin listener has no auth. Mitigation: bind to localhost or use network-level access control. Auth is a non-goal for this iteration.
- **[Port conflict]** → Default `:8080` may conflict with other services. Mitigation: configurable via `--admin-addr` flag.
- **[Stale metadata]** → If the binary is copied without rebuilding, `commit`/`build_time` reflect the original build. Mitigation: this is expected behavior; documented via default `"unknown"` values.
