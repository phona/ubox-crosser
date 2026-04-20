---
id: REQ-669
title: Implementation tasks
---

# Tasks

## Stage 1: Core Implementation

- [x] Create `version/version.go` with build variables
- [x] Create `api/version.go` HTTP handler
- [x] Wire HTTP server into `cmd/server/server.go`
- [x] Add `--api-addr` flag
- [x] Update Makefile ldflags for version injection

## Stage 2: Testing

- [x] Unit tests for handler (status, content-type, body fields)
- [ ] Contract test against OpenAPI schema

## Stage 3: Documentation

- [x] OpenSpec proposal, design, spec, tasks
- [x] Contract spec (OpenAPI YAML)
