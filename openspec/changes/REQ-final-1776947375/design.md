# REQ-final-1776947375 — Design

## High-level approach

Add a new `GET /buildinfo` HTTP handler on the server's existing control-plane
HTTP listener (port 8080, the same one used by the in-flight `/healthz`
endpoint). The handler marshals a fixed 3-field JSON struct:

```go
type buildInfoResponse struct {
    GitSHA    string `json:"git_sha"`
    BuildID   string `json:"build_id"`
    GoVersion string `json:"go_version"`
}
```

- `GitSHA` is read from a package-level `var GitSHA = "dev"` in the
  `github.com/phona/ubox-crosser/server` package. At build time the Makefile
  overrides it via `-ldflags "-X github.com/phona/ubox-crosser/server.GitSHA=<7-char>"`.
- `BuildID` is read at **request time** (not program-start time) from
  `os.Getenv("BUILD_ID")`. If empty, respond with `"dev"`. Reading per-request
  is fine — the call is cheap and lets ops rotate the env var without a
  restart if ever needed. (We accept the minor tradeoff: the per-request
  read means a test can `os.Setenv` mid-test.)
- `GoVersion` is the hardcoded string `"go1.23"`. The prompt is explicit:
  do not derive from `runtime.Version()`. If the Go toolchain is upgraded,
  this string must be bumped in the same PR.

## Nuance on the "existing /version" reference

The intent issue says "reference the existing `/version` endpoint". No
`/version` handler exists in `master` as of commit `233e8fa`. The closest
precedent is `/healthz`, whose **acceptance spec** landed in `233e8fa`
(REQ-e2e-1776916220) but whose **implementation has not yet merged**.

**Decision (no escalate, per KPI):** treat the phrase "existing /version"
as loose shorthand for "the ldflags-injection pattern commonly used for
`/version`". This REQ is the first place the pattern lands in the repo;
`/buildinfo` introduces:
- the `GitSHA` package variable + `-X` ldflags on server builds
- a convention future `/version`-style endpoints can copy from

If `/healthz` lands first, its HTTP mux is reused; if `/buildinfo` lands
first, this REQ brings up a minimal `net/http` server on :8080 and `/healthz`
layers onto it. See §"Dependency on /healthz" for coordination.

## Dependency on /healthz (REQ-e2e-1776916220)

`/healthz` and `/buildinfo` both need an HTTP listener on port 8080. Only
one REQ should bring it up. The two can land in either order:

- **If /healthz merges first** — its HTTP mux already exists.  
  `/buildinfo` dev stage just registers an additional handler.
- **If /buildinfo merges first** — this REQ's dev stage must bring up the
  shared HTTP mux (small `http.ServeMux` on :8080, started in a goroutine
  from `server.NewProxyServer`'s bootstrap path).  
  `/healthz` dev then registers onto the same mux.

**Recommendation:** the dev-agent should make the HTTP-server bring-up code
a small, named helper (e.g. `server/http_control.go` with a
`StartControlPlaneHTTP(mux *http.ServeMux, addr string)` func) so the
ordering does not matter and both REQs converge cleanly. Worst-case merge
conflict is in `cmd/server/server.go` where handlers are registered; both
sides are additive so resolution is trivial.

## Ldflags wiring

Current `Makefile:15-17`:

```make
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/server
```

Target state (server target only; client / auth_server unaffected):

```make
GIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
SERVER_LDFLAGS := -s -w -X github.com/phona/ubox-crosser/server.GitSHA=$(GIT_SHA)

CGO_ENABLED=0 go build -ldflags "$(SERVER_LDFLAGS)" -o bin/server ./cmd/server
```

Current `Dockerfile:15-16`:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/crosser ./cmd/${BINARY}
```

Target: add `ARG GIT_SHA=dev` and thread it into the same `-X` flag:

```dockerfile
ARG GIT_SHA=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w -X github.com/phona/ubox-crosser/server.GitSHA=${GIT_SHA}" \
    -o /app/crosser ./cmd/${BINARY}
```

`BUILD_ID` stays a pure env var — no Makefile/Dockerfile change required
for it (it's read at request time).

## Risks and tradeoffs

| Risk | Likelihood | Mitigation |
|---|---|---|
| Hardcoded `"go1.23"` string drifts from actual toolchain | medium | CI runs on Go 1.23 (Makefile `ci-env`). Add a unit-test assertion `runtime.Version() hasPrefix "go1.23"` to catch toolchain bumps. |
| `GitSHA` left as literal string `"dev"` when ldflags missing | low | Default `var GitSHA = "dev"` — cosmetic, does not break the contract (still 3-char or 7-char string). Acceptance test asserts field presence, not format. |
| Docker build does not inject `GIT_SHA` in CI | medium | `ci-build` target should pass `--build-arg GIT_SHA=$$(git rev-parse --short HEAD)`. Flagged as explicit task in `tasks.md`. |
| HTTP listener not brought up by this REQ alone | see §Dependency | Dev helper `StartControlPlaneHTTP` keeps ordering agnostic. |

## Parallel-split evaluation

Candidate task breakdown (per the analyze-agent heuristic):

1. **handler + unit test** (`server/buildinfo.go`, `server/buildinfo_test.go`)
2. **ldflags wiring** (Makefile + Dockerfile)
3. **acceptance test** (`tests/acceptance/buildinfo_test.go` + compose tweak)

Strong dependency: (3) asserts output fields, which requires (1). (2)
only makes one field non-empty — (1) defaults `GitSHA` to `"dev"` so (3)
can still pass without (2).

**Decision: do NOT split.** Total surface is <150 LOC across ~5 files, all
cohesive. A split would add coordination overhead (two dev agents, two
PRs, two CI runs) that dwarfs the wall-clock savings. Keep as a single
`tag=dev` subissue. This is the "single block, < 100 LOC" case in the
heuristic.

## Open questions (resolved without escalate)

1. **"existing /version"** — addressed in §"Nuance on the 'existing /version' reference".
2. **`go_version` field** — prompt hardcodes `"go1.23"`. If the intent was "the
   actual runtime Go version", the prompt would have said `runtime.Version()`.
   Following the letter of the spec.
3. **7-char short SHA** — `git rev-parse --short` defaults to 7 chars on
   this repo. No explicit `--short=7` needed, but we pin it in the Makefile
   recipe as `--short=7` for determinism.
