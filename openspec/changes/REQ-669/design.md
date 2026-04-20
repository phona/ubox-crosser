---
id: REQ-669
title: GET /version design
---

# Design: GET /version

## Approach

Standalone `net/http` server running alongside the TCP proxy. No external framework needed for a single endpoint.

## Trade-offs

| Option | Pros | Cons |
|--------|------|------|
| Embed in TCP protocol | No extra port | Breaks protocol simplicity, clients must understand it |
| Separate HTTP server (chosen) | Standard tooling, curl-friendly, extensible | Extra port to expose |

## Decision

Separate HTTP server on `--api-addr`. Minimal blast radius — if the API goroutine panics, the proxy continues.

## Build-time injection

Use `go build -ldflags -X` to set `version.Commit` and `version.BuildTime` at compile time. Fallback to "unknown" when built without flags (e.g., `go run`).
