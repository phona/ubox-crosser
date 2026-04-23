# REQ-final-1776947375: /buildinfo endpoint — tasks

Owner stages below are the standard analyze→spec→dev handoff. Analyze-agent
only fills in TODO bullets; spec-agent and dev-agent concretize their
respective stages when they pick up.

## Stage: contract-tests (owner: spec-agent)

- [ ] TODO: write `specs/buildinfo-endpoint/spec.md` with response schema
      (fields: `git_sha` string 7-char, `build_id` string non-empty,
      `go_version` string equal to `"go1.23"`) and HTTP contract
      (method `GET`, status 200, `Content-Type: application/json`,
      no auth)
- [ ] TODO: document the three failure modes that are **out of scope** and
      should NOT be tested as contracts (non-GET methods, auth, arbitrary
      future fields)

## Stage: acceptance-tests (owner: spec-agent)

- [ ] TODO: write user-facing acceptance scenarios in Given/When/Then
      form — minimum coverage:
  - FEATURE-A1: `GET /buildinfo` returns 200 with valid JSON body
  - FEATURE-A2: response contains `git_sha`, `build_id`, `go_version`
    string fields (presence + non-nil)
  - FEATURE-A3: `build_id` defaults to `"dev"` when `BUILD_ID` env var
    is unset at server startup
  - FEATURE-A4: `go_version` is the literal string `"go1.23"`
  - FEATURE-A5: endpoint is unauthenticated (no `Authorization` header
    required)
- [ ] TODO: integration test `tests/acceptance/buildinfo_test.go`
      reuses `getHealthCheckAddr()` helper (port 8080) — same pattern
      as `healthz_test.go`
- [ ] TODO: extend `tests/acceptance/docker-compose.yml` `test-client`
      to also `curl -fsS http://ubox-crosser:8080/buildinfo` after
      `/healthz`

## Stage: implementation (owner: dev-agent)

**Scope is single, cohesive, < 150 LOC. Do NOT split into parallel
subissues** (see `design.md` §"Parallel-split evaluation").

### Code
- [ ] TODO: `server/buildinfo.go` — declare `var GitSHA = "dev"`;
      implement `BuildInfoHandler(w http.ResponseWriter, r *http.Request)`
      that writes the 3-field JSON response (only accept GET; method
      check can be loose per acceptance spec)
- [ ] TODO: `cmd/server/server.go` — register `BuildInfoHandler` on the
      control-plane HTTP mux on port 8080. If `/healthz` has not landed
      yet, bring up the shared mux (see `design.md` §"Dependency on
      /healthz"); use the helper `StartControlPlaneHTTP(mux, ":8080")`
- [ ] TODO: `server/buildinfo_test.go` — unit test:
  - covers `BUILD_ID` unset → `"dev"` branch (use `t.Setenv`)
  - covers `BUILD_ID` set → echoes it
  - asserts `go_version == "go1.23"` literal
  - asserts response `Content-Type: application/json`

### Build plumbing
- [ ] TODO: `Makefile` — introduce `GIT_SHA := $(shell git rev-parse --short=7 HEAD 2>/dev/null || echo dev)` and add
      `-X github.com/phona/ubox-crosser/server.GitSHA=$(GIT_SHA)` to the
      `bin/server` build only (not client / auth_server)
- [ ] TODO: `Dockerfile` — add `ARG GIT_SHA=dev` and thread it into the
      `go build -ldflags` flag
- [ ] TODO: `Makefile` `ci-build` — pass `--build-arg GIT_SHA=$$(git rev-parse --short=7 HEAD)` to the
      `docker build` for the server image
- [ ] TODO: verify `make build` locally produces a binary whose
      `/buildinfo` returns the real 7-char SHA (not `"dev"`)

### CI + regression
- [ ] TODO: confirm `ci-unit-test` covers the new `buildinfo_test.go`
      (no Makefile change expected — it globs `./...`)
- [ ] TODO: confirm `ci-integration-test` covers the new acceptance
      test (docker-compose up flow already wired)
- [ ] TODO: add a toolchain-drift guard unit test:
      `assert strings.HasPrefix(runtime.Version(), "go1.23")` so a
      future Go bump surfaces as a failing unit test, not a silent
      lie in the JSON response

## Acceptance check (final)

- [ ] TODO: `curl -fsS http://<server>:8080/buildinfo` returns 200
      with JSON containing `git_sha`, `build_id`, `go_version`
      (the HARD acceptance bar from the intent issue)
