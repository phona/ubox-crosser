---
change_id: REQ-642
title: "Tasks: GET /version endpoint"
---

# Tasks

## Stage 1: Core Implementation

- [x] Create `version/version.go` with Version constant and ldflags vars
- [x] Create `version/handler.go` with HTTP handler
- [x] Create `version/handler_test.go` with unit tests
- [x] Wire handler into `cmd/server` with `--http-addr` flag
- [x] Update Makefile ldflags for commit and build_time injection
- [x] Update Dockerfile with build args for ldflags

## Stage 2: Verification

- [x] Unit tests pass (`go test ./version/`)
- [x] `go vet ./...` clean
