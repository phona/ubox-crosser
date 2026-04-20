---
change_id: req-686
title: "GET /ping endpoint — Design"
---

## Context

ubox-crosser already has an admin HTTP listener (`:8080`) serving `GET /version` and `GET /healthz`. Operators and clients need the simplest possible endpoint to confirm the server process is reachable at the network level — a classic ping/pong pattern.

## Goals / Non-Goals

**Goals:**
- Expose `GET /ping` on the existing admin HTTP listener returning plain-text `pong`
- Minimal overhead: no JSON encoding, no Content-Type negotiation
- Accept only GET; reject other methods with 405

**Non-Goals:**
- Health semantics (that is `/healthz`)
- Build metadata (that is `/version`)
- Round-trip latency measurement or timestamped responses

## Decisions

### Decision 1: Plain-text response (not JSON)

A ping/pong endpoint is an echo check. Plain text avoids unnecessary JSON encoding overhead and is the idiomatic format for this pattern.

**Chosen:** Return `pong` as `text/plain; charset=utf-8`.

| Option | Pros | Cons |
|--------|------|------|
| Plain text `pong` | Minimal, idiomatic, zero-alloc | Not JSON like other endpoints |
| JSON `{"ping":"pong"}` | Consistent with other endpoints | Unnecessary overhead for an echo |

### Decision 2: Separate `ping` package

Follow the same pattern as `version` and `health` packages — a small, self-contained package with its own handler and tests.

## Risks / Trade-offs

- **[Plain text vs JSON]** → Breaks JSON consistency with other admin endpoints. Acceptable because ping/pong is universally understood as plain text and the simplicity is the point.
- **[No auth]** → Same as `/version` and `/healthz` — relies on network-level access control.
