# REQ-m15v2-1776941409: /buildinfo Endpoint — Tasks

## Stage: contract-tests (owner: spec-agent)
- [ ] TODO: Author `openspec/changes/REQ-m15v2-1776941409/specs/buildinfo-endpoint/spec.md` describing `GET /buildinfo` contract: 200 response, `Content-Type: application/json`, body shape `{git_sha, build_id, go_version}`. Use Given/When/Then in the same style as `specs/healthz-endpoint/spec.md`.
- [ ] TODO: Add field-level constraints in spec — `git_sha` is exactly 7 chars `[0-9a-f]`, `build_id` is non-empty string, `go_version == "go1.23"`.
- [ ] TODO: Run `openspec validate REQ-m15v2-1776941409` (or repo equivalent) to confirm spec parses.

## Stage: acceptance-tests (owner: spec-agent)
- [ ] TODO: Add `tests/acceptance/buildinfo_test.go` mirroring `healthz_test.go` shape: hit `http://server:<healthz-port>/buildinfo`, assert status 200, assert all three JSON fields populated per the spec constraints above.
- [ ] TODO: Confirm the test reuses the existing `tests/acceptance/docker-compose.yml` (no new compose stack).
- [ ] TODO: Verify the test currently FAILS against the current binary (handler not yet implemented) — this is the red baseline dev-agent flips to green.

## Stage: implementation (owner: dev-agent)
- [ ] TODO: Add `var GitSHA string` at package level in `cmd/server/server.go` (or a sibling `cmd/server/buildinfo.go`).
- [ ] TODO: Implement `buildinfoHandler(w http.ResponseWriter, r *http.Request)` returning the JSON body specified in the spec.
- [ ] TODO: Register the handler on the same `*http.ServeMux` that already serves `/healthz` — do NOT open a new listener and do NOT use `http.DefaultServeMux`.
- [ ] TODO: Add `cmd/server/buildinfo_test.go` unit test calling the handler via `httptest.NewRecorder()` — assert status, content-type, all three fields, default `build_id == "dev"` when env unset, custom `build_id` when env set.
- [ ] TODO: Wire `-ldflags "-X main.GitSHA=$(git rev-parse --short HEAD)"` into the build path that produces the docker image used by `tests/acceptance/docker-compose.yml` (Dockerfile and/or Makefile — confirm during dev).
- [ ] TODO: Run `go test ./cmd/server/...` locally → green. Run `go test ./tests/acceptance/...` against docker-compose stack → green.
