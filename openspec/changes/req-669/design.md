---
change_id: req-669
title: "GET /version endpoint — Technical Design"
---

## Context

ubox-crosser is a proxy tunnel server. Operators need runtime visibility into which version is deployed on each node. The server currently has no admin/health endpoints. This design adds a minimal admin HTTP listener with a single `GET /version` endpoint.

## Goals / Non-Goals

**Goals:**
- Expose version, commit SHA, and build timestamp via a single JSON endpoint
- Zero external dependencies (stdlib `net/http` only)
- Compile-time injection via standard `-ldflags -X` pattern

**Non-Goals:**
- Authentication on the admin endpoint (internal network only)
- Health checks, metrics, or readiness probes (future work)
- Graceful shutdown of the admin listener

## Decisions

### Use Go 1.22+ method-based routing
Register `"GET /version"` pattern on `http.ServeMux`. This rejects non-GET methods with 405 automatically without a custom middleware.

**Alternative**: Manual method check in handler — rejected because it duplicates what the stdlib already provides.

### Separate admin listener from proxy traffic
The admin HTTP server runs on its own goroutine and address (`--admin-addr`, default `:8080`), isolated from the proxy data path.

**Alternative**: Multiplex on the proxy port — rejected because it complicates protocol detection and exposes admin endpoints to untrusted traffic.

### Hardcoded version constant + injected build vars
`Version` is a compile-time constant in `version/version.go`. `Commit` and `BuildTime` are `var` with defaults, overridden by ldflags at build time.

**Alternative**: Read from embedded file — rejected for simplicity; ldflags is the standard Go pattern.

## Risks / Trade-offs

- [Admin port exposed without auth] → Mitigated by binding to internal/loopback by default in production deployments; document in README.
- [Admin server crash kills process] → `logrus.Fatalf` on listen failure is intentional; if the port is occupied, the operator should know immediately.
