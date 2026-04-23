# REQ-final7-1776959695: /buildinfo Endpoint

## Stage: contract / spec

- [x] Author `specs/buildinfo-endpoint/contract.spec.yaml` — OpenAPI 3.0 schema for /buildinfo
- [x] Author `specs/buildinfo-endpoint/spec.md` — 7 scenarios in openspec delta format

## Stage: implementation

- [x] Create `server/http.go` with `HealthzHandler`, `BuildInfoHandler(gitSHA)`, `StartHTTPServer`
- [x] Add `var GitSHA string` to `cmd/server/server.go`; call `server.StartHTTPServer(":8080", GitSHA)`
- [x] Update `Makefile` build target with `-X main.GitSHA=$(GIT_SHA)` ldflag for server binary
- [x] Add `ci-test`, `ci-accept-env-up`, `ci-accept-env-down` targets to Makefile
- [x] Update `Dockerfile` with `GIT_SHA` build arg and conditional ldflags for server
- [x] Update `tests/Dockerfile.test` with `GIT_SHA` build arg and ldflags for server binary

## Stage: testing

- [x] Unit tests in `server/http_test.go` (5 cases: healthz 200, healthz JSON, buildinfo 200, buildinfo fields, buildinfo default build_id, content-type)
- [x] Acceptance tests in `tests/acceptance/buildinfo_test.go` (5 scenarios: 200, three fields, go_version, build_id from env, content-type); `//go:build acceptance` tag
- [x] Update `tests/acceptance/healthz_test.go`: add `//go:build acceptance` tag; `getHealthCheckAddr()` reads `SERVER_ADDR` env
- [x] Update `tests/acceptance/docker-compose.yml`: add `BUILD_ID`, `GIT_SHA`, test-runner service with `-tags acceptance`

## Stage: PR

- [x] Commit all changes to `feat/REQ-final7-1776959695`
- [x] `git push origin feat/REQ-final7-1776959695`
- [x] `gh pr create` with full description
