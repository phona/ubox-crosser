---
change_id: req-889
title: "Add GET /buildinfo endpoint"
repos: [ubox-crosser]
layers:
  - backend
status: draft
---

## Why

Provide a `/buildinfo` endpoint that returns the running binary's version, git commit, and build timestamp. This is a common DevOps convention — some tooling expects `/buildinfo` specifically, distinct from `/version`.

## What Changes

- Register `GET /buildinfo` on the admin `http.ServeMux`, reusing the existing `version.Handler`
- The response is identical to `GET /version`: `{"version":"...","commit":"...","build_time":"..."}`
- No new packages or dependencies needed — this is a one-line route addition

### Design Decision: Reuse vs. New Package

The existing `version.Handler` already returns exactly the requested data (version + commit + build_time). Two options:

1. **Reuse `version.Handler`** (recommended): Add `mux.HandleFunc("GET /buildinfo", version.Handler)` — zero new code, zero maintenance cost
2. **New `buildinfo/` package**: Separate handler with its own struct — adds duplication for no functional benefit

Option 1 is recommended. If the response format needs to diverge in the future, a dedicated package can be extracted at that point.

## Capabilities

### New Capabilities
- `buildinfo-endpoint`: HTTP GET /buildinfo returning JSON `{"version","commit","build_time"}` with HTTP 200

### Modified Capabilities
- None

## Impact

- `cmd/server/server.go` — add one route registration line
- No new files, packages, or dependencies
- Fully covered by existing `version.Handler` unit tests; only mux-level routing test needed
