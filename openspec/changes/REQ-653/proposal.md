---
id: REQ-653
title: GET /version endpoint
layers:
  - backend
status: proposed
---

## Summary

Add a `GET /version` HTTP endpoint to the ubox-crosser server that returns build metadata (version, git commit, build time) as JSON.

## Motivation

Operators and CI pipelines need a way to verify which version of ubox-crosser is deployed without inspecting the binary directly. A standard `/version` endpoint provides this introspection over HTTP.

## Scope

- New `version` package with build-time injected variables
- HTTP handler returning `{"version","commit","build_time"}`
- Admin HTTP server on configurable `--admin-addr` (default `:8080`)
- Makefile and Dockerfile ldflags for build-time injection
- Unit tests for the handler
