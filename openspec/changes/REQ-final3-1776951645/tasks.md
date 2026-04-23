# REQ-final3-1776951645: /buildinfo HTTP Endpoint — Tasks

owner: analyze-agent

## Stage: contract / spec

- [x] Author `specs/buildinfo-endpoint/contract.spec.yaml` with the
  `GET /buildinfo` contract (status codes, headers, JSON schema).
- [x] Author `specs/buildinfo-endpoint/spec.md` with acceptance
  scenarios `[UBOX-CROSSER-S1]` .. `[UBOX-CROSSER-S5]`.

## Stage: implementation

- [x] Add `server/admin_http.go` with `StartAdminHTTP(addr, mux)`
  helper bringing up a minimal `http.Server` on a supplied mux.
- [x] Add `cmd/server/buildinfo.go` declaring `var GitSHA = "unknown"`,
  the `BuildInfo` struct, and `BuildInfoHandler`.
- [x] Wire the admin listener into `cmd/server/server.go` — new mux,
  register `/buildinfo`, `go server.StartAdminHTTP(":8080", mux)`.

## Stage: build plumbing

- [x] Add `GIT_SHA` variable + `SERVER_LDFLAGS` to `Makefile`; thread
  into `build` and `ci-build` targets for `./cmd/server` only.
- [x] Thread `ARG GIT_SHA` through the root `Dockerfile` and the server
  `go build` invocation.
- [x] Thread `ARG GIT_SHA` through `tests/Dockerfile.test` for the
  coverage-instrumented server build.

## Stage: tests

- [x] Unit tests `cmd/server/buildinfo_test.go`: status 200, JSON
  unmarshals, `BUILD_ID` env override, `GitSHA` var override,
  `go_version` literal, `Content-Type` header.
- [x] Unit test cross-checking the `go_version` literal against
  `go.mod`'s `go` directive to fail CI on toolchain drift.
- [x] Integration test `tests/integration/buildinfo_test.go`
  (`//go:build integration`): GET against `proxy-server:8080/buildinfo`,
  assert status + shape + ldflag-injected `GitSHA` + env-injected
  `BuildID` + literal `GoVersion`.
- [x] Update `tests/docker-compose.yml` — proxy-server gains `BUILD_ID`
  env and `GIT_SHA` build arg; `ci-integration-test` Makefile target
  passes both down.

## Stage: PR

- [x] Commit on `feat/REQ-final3-1776951645`.
- [x] `git push -u origin feat/REQ-final3-1776951645`.
- [x] `gh pr create` with title + body explaining motivation and test
  plan.
- [x] Verify `make ci-unit-test` passes in runner pod.
- [x] Verify `make ci-integration-test` passes in runner pod.
