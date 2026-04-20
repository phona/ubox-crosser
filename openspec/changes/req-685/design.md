---
change_id: req-685
title: "GET /healthz endpoint — Design"
---

## Context

ubox-crosser already has an admin HTTP listener (`:8080`) serving `GET /version`. Operators need a standard health-check endpoint for liveness probes. The `/healthz` path is the Kubernetes convention and widely recognized by infrastructure tooling.

## Goals / Non-Goals

**Goals:**
- Expose `GET /healthz` on the existing admin HTTP listener returning `{"status":"ok"}`
- Keep the response static and zero-dependency (no deep checks)
- Accept only GET; reject other methods with 405

**Non-Goals:**
- Readiness checks (database, upstream connectivity)
- Metrics endpoint (`/metrics`)
- Configurable health conditions or degraded states

## Decisions

### Decision 1: Static JSON response (liveness only)

A liveness probe only needs to confirm the process is running and the HTTP listener is accepting connections. Deep health checks belong in a readiness probe (future work).

**Chosen:** Return a fixed `{"status":"ok"}` with no runtime checks.

| Option | Pros | Cons |
|--------|------|------|
| Static response | Simple, fast, no false negatives | Does not detect degraded state |
| Check tunnel goroutine | Detects tunnel failures | Couples health to tunnel internals, risk of false negatives |

### Decision 2: Separate `health` package

Follow the same pattern as the `version` package — a small, self-contained package with its own handler and tests. Keeps the admin endpoint handlers decoupled and independently testable.

### Decision 3: Use `/healthz` path (not `/health`)

`/healthz` is the Kubernetes liveness convention. Using the standard path avoids needing custom probe configuration in pod specs.

## Risks / Trade-offs

- **[Static liveness only]** → This endpoint does not detect degraded tunnel state. Mitigation: this is intentional for v1; readiness checks are future work.
- **[No auth]** → Same as `/version` — relies on network-level access control.
