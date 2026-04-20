---
id: REQ-653
title: GET /version endpoint – Tasks
---

## Stage 1: Backend implementation

- [x] Create `version/version.go` with Version constant and Commit/BuildTime vars
- [x] Create `version/handler.go` with HTTP handler
- [x] Create `version/handler_test.go` with unit tests
- [x] Update `Makefile` ldflags to inject Commit and BuildTime
- [x] Update `Dockerfile` to accept and pass build args
- [x] Integrate admin HTTP server in `cmd/server/server.go`

## Stage 2: Verification

- [x] Unit tests pass (`go test ./version/...`)
- [x] Build compiles (`go build ./cmd/server/`)
- [x] `go vet ./...` passes
