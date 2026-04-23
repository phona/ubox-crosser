# REQ-final2-1776948458: /buildinfo HTTP Endpoint — Tasks

## Stage: contract-tests (owner: spec-agent)

- [ ] TODO: author `specs/buildinfo-endpoint/contract.spec.yaml` describing
      the HTTP contract (path, method, response schema, status codes)
- [ ] TODO: cover no-auth requirement in the contract
- [ ] TODO: run `openspec validate` locally; fix any schema complaints

## Stage: acceptance-tests (owner: spec-agent)

- [ ] TODO: author `specs/buildinfo-endpoint/spec.md` with Given/When/Then
      scenarios (FEATURE-B1..FEATURE-Bn) for the three-field JSON, the env
      default behaviour, concurrency, and the 200 / content-type guarantees
- [ ] TODO: confirm scenario IDs don't collide with FEATURE-A* reserved by
      the `/healthz` REQ
- [ ] TODO: run `check-scenario-refs.sh` on the change dir

## Stage: implementation (owner: dev-agent)

- [ ] TODO: add `buildinfoHandler` + `BuildInfo` struct
- [ ] TODO: register `/buildinfo` on the admin HTTP mux (piggyback on the
      `/healthz` listener if already merged; otherwise bootstrap the listener
      here — see design.md §"Listener placement")
- [ ] TODO: wire `-X main.GitSHA=…` through `Makefile` and `Dockerfile`
- [ ] TODO: unit test covering default env, explicit env, and the JSON shape
- [ ] TODO: integration test under `tests/acceptance/` that hits
      `http://ubox-crosser:8080/buildinfo` via docker-compose and asserts
      status + all three fields
- [ ] TODO: drift-guard unit test comparing the hard-coded `go_version`
      literal to the `go` directive in `go.mod`
