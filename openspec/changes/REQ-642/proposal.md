---
change_id: REQ-642
title: "Add GET /version endpoint"
layers:
  - backend
status: implementing
---

# REQ-642: Add GET /version endpoint

## Problem

There is no way to query the running server's version, git commit, or build time. This makes it hard to verify deployments and troubleshoot version mismatches.

## Solution

Add an HTTP `GET /version` endpoint to the proxy server that returns version metadata as JSON. The server already runs a TCP listener for its tunnel protocol; a lightweight HTTP server is added alongside it on a configurable `--http-addr` (default `:8080`).

## Scope

- New `version` package with build-time injected `Commit` and `BuildTime` vars
- HTTP handler for `GET /version`
- Makefile and Dockerfile updated with `-ldflags` for compile-time injection
- Unit tests for the handler
- No authentication required
