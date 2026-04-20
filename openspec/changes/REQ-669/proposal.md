---
id: REQ-669
title: GET /version endpoint
layers:
  - backend
status: approved
---

# Proposal: GET /version Endpoint

## Summary

Add an HTTP `GET /version` endpoint to the proxy server that returns build metadata (version, git commit, build time) as JSON.

## Motivation

Operators and monitoring systems need a lightweight way to verify which version of the server is running without inspecting the binary directly.

## Scope

- New `version` package holding build-time variables
- New `api` package with HTTP handler
- HTTP listener on configurable `--api-addr` (default `:8080`)
- Makefile ldflags injection of commit and build time
- No authentication required
