# Tasks: REQ-final8-1776962929

## Stage: contract / spec

- [x] author specs/buildinfo/contract.spec.yaml — define /buildinfo JSON contract
- [x] author specs/buildinfo/spec.md — write scenarios in openspec delta format

## Stage: implementation

- [x] create server/health.go — BuildInfo struct, NewBuildInfoHandler, NewHealthzHandler, StartHTTPServer
- [x] modify cmd/server/server.go — add var GitSHA, read BUILD_ID env, call StartHTTPServer
- [x] create server/health_test.go — pure unit tests for both handlers (no build tag)
- [x] add Makefile targets ci-test and dev-cross-check

## Stage: PR

- [x] git push feat/REQ-final8-1776962929
- [x] gh pr create
