---
change_id: req-889
title: "GET /buildinfo endpoint — design"
---

## Context

The admin HTTP server already exposes `GET /version` returning `{"version":"0.1.0","commit":"<sha>","build_time":"<iso>"}` via `version.Handler`. The request is to add `GET /buildinfo` with the same data.

## Goals

1. Expose build metadata at `/buildinfo` for tooling that expects this path
2. Minimize code duplication — reuse existing handler

## Decision

**Route alias approach**: Register `version.Handler` on both `/version` and `/buildinfo`. This avoids any new packages, structs, or build-flag injection.

```go
mux.HandleFunc("GET /buildinfo", version.Handler)
```

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| Two paths returning identical data may confuse consumers | Document both in API docs; `/version` is the canonical path |
| Future divergence (e.g., `/buildinfo` adds extra fields) | Extract to its own package at that point — trivial refactor |

## Dependencies

- `version` package (existing) — no changes needed
- Build-time ldflags injection (existing) — no changes needed
