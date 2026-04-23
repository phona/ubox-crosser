# Tasks: REQ-final10-1776966282

## Stage: contract / spec

- [x] author specs/buildinfo/contract.spec.yaml — define /buildinfo JSON contract
- [x] author specs/buildinfo/spec.md — write scenarios in openspec delta format

## Stage: implementation

- [x] create server/health.go — BuildInfo struct, NewBuildInfoHandler, NewHealthzHandler, StartHTTPServer
- [x] modify cmd/server/server.go — add var GitSHA, read BUILD_ID env, call StartHTTPServer
- [x] create server/health_test.go — pure unit tests for both handlers (no build tag)
- [x] add Makefile targets ci-test and dev-cross-check
- [x] add //go:build integration tag to tests/acceptance/healthz_test.go

## Stage: PR

- [x] git push feat/REQ-final10-1776966282
- [x] gh pr create
