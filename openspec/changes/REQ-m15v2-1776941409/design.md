# Design — REQ-m15v2-1776941409: /buildinfo Endpoint

## Approach (high level)

Mount one new handler on the existing health-check HTTP mux in `cmd/server`. The handler reads three values:

| Field        | Source                                          | Default   |
|--------------|-------------------------------------------------|-----------|
| `git_sha`    | package-level `var GitSHA string` (ldflag-set)  | `""` (build that forgets the ldflag will surface as empty string — acceptance test asserts non-empty) |
| `build_id`   | `os.Getenv("BUILD_ID")`                         | `"dev"`   |
| `go_version` | hardcoded string literal `"go1.23"`             | n/a       |

Then `json.Marshal` → `w.Write`. ~25 LOC of handler + ~15 LOC of unit test.

## Why hardcoded `go_version`

The prompt explicitly says hardcode `"go1.23"`. We do **not** call `runtime.Version()` because the spec wants this string stable across patch releases (e.g. `go1.23.4` vs `go1.23.7`) and independent of which Go toolchain happens to build the image.

## Why no new listener

`/healthz` already runs an `http.Server` on a configurable health-check port inside `cmd/server`. Reusing the same `*http.ServeMux` keeps the change to a single `mux.HandleFunc("/buildinfo", ...)` line at registration time — no port plumbing, no config addition, no new lifecycle to manage.

## ldflag wiring

`dev-agent` must update whatever currently builds the server binary (Dockerfile / Makefile / both — to be confirmed during dev) to pass:

```bash
go build -ldflags "-X main.GitSHA=$(git rev-parse --short HEAD)" ./cmd/server
```

The integration test runs against the docker-compose stack from `tests/acceptance/docker-compose.yml`, so the Dockerfile that compose builds **is** the wire under test. If the Makefile builds outside docker, that path must also be updated for the unit test toolchain.

## Risks

1. **Empty GitSHA on plain `go build`** — without the ldflag, `GitSHA` is `""`. Acceptance test must check `len(git_sha) == 7` (not just non-nil) so an unwired build path fails loudly rather than silently returning `""`.
2. **`/healthz` mux collision** — none expected (different path), but `dev-agent` should add the route on the same `mux` value, not register a competing `http.HandleFunc` against `http.DefaultServeMux`.
3. **Build cache vs. SHA injection** — Go caches by package source + flags, so changing `GitSHA` re-links but does not recompile. Acceptable; integration test only requires a 7-char value, not freshness across commits within a single CI run.

## Parallelism evaluation (M16-style)

Total work is roughly:
- 1 handler + registration  (~25 LOC, 1 file)
- 1 unit test               (~30 LOC, 1 new file)
- 1 acceptance test         (~50 LOC, 1 new file, reuses healthz docker-compose)
- 1 ldflag tweak in build script

Critical path is short (well under 100 LOC across 3-4 files). **Strong dependencies dominate:** the acceptance test cannot be authored against a non-existent route handler shape; the ldflag change is meaningless without the package-level `var GitSHA`. Splitting into parallel dev sub-issues would create more coordination overhead than it saves.

**Decision: do NOT split.** Single `tag=dev` issue, single `tag=spec` issue. (Acknowledging the M16 default leans toward splitting — this REQ is the small-change exception.)

## Stage handoff contract

- **spec stage** writes `openspec/changes/REQ-m15v2-1776941409/specs/buildinfo-endpoint/spec.md` with Given/When/Then scenarios mirroring the `/healthz` spec format. Acceptance test file `tests/acceptance/buildinfo_test.go` is authored here too (test exists, fails because handler doesn't exist yet).
- **dev stage** implements the handler, registers it on the health-check mux, wires the ldflag, and adds `cmd/server/buildinfo_test.go` (unit test). Acceptance test from spec stage flips to passing.
- **accept stage** runs `tests/acceptance/docker-compose.yml` and asserts `curl -fsS http://server:<healthz-port>/buildinfo` returns 200 with all three fields populated.
