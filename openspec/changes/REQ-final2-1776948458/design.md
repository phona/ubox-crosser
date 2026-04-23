# REQ-final2-1776948458: Design Notes

## High-level approach

Add a single HTTP handler `buildinfoHandler` that serializes a fixed struct to
JSON and writes it with `Content-Type: application/json`. Register it on the
existing admin mux next to `/healthz`. Inject `GitSHA` via linker at build time.
Read `BUILD_ID` from env on each request (cheap; no need to snapshot at startup,
and it lets operators surface per-run identifiers without rebuilding).

## Data model

```go
type BuildInfo struct {
    GitSHA    string `json:"git_sha"`    // 7-char, injected via -ldflags
    BuildID   string `json:"build_id"`   // env BUILD_ID, default "dev"
    GoVersion string `json:"go_version"` // hard-coded "go1.23"
}
```

`GitSHA` is a package-level `var GitSHA = "unknown"` in `main` (or a dedicated
`buildinfo` package — dev-agent picks placement). It stays `"unknown"` if the
ldflag isn't set, which is fine for local `go run`.

## Build-time injection

`Makefile` `build` + `ci-build` targets get the ldflag:

```make
GIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS := -s -w -X main.GitSHA=$(GIT_SHA)
```

`Dockerfile` needs to pass `GIT_SHA` in as a build arg because the `.git`
directory isn't copied into the build context in the current `COPY . .` flow
(and shallow-clone builds may not have full history anyway). Concretely:

```dockerfile
ARG GIT_SHA=unknown
RUN go build -ldflags="-s -w -X main.GitSHA=${GIT_SHA}" -o /app/crosser ./cmd/${BINARY}
```

`ci-build` / `test-integration` in `Makefile` passes `--build-arg GIT_SHA=$(GIT_SHA)`.

## Listener placement

The admin HTTP listener (port `8080`) is introduced by REQ-e2e-1776916220
(`/healthz`). Two dev paths, chosen at impl time based on branch state:

1. **`/healthz` already merged**: just `mux.HandleFunc("/buildinfo", …)` next to
   the existing `/healthz` registration.
2. **`/healthz` not yet merged**: dev-agent brings up the admin listener here
   (new file e.g. `server/admin_http.go`) and registers both handlers, staying
   compatible with the `/healthz` spec so the two REQs compose cleanly on merge.

Either way, `cmd/server/server.go` calls the admin-listener bootstrapper once,
in a goroutine, before blocking on the proxy event loop.

## Testing

- **Unit** (`server/admin_http_test.go` or equivalent, default build tag): drive
  the handler with `httptest.NewRecorder`, assert 200, JSON parse, field values
  (`BuildID` respects env with `t.Setenv`; `GitSHA` via overriding the package
  var in the test).
- **Integration** (`tests/acceptance/buildinfo_test.go`, reuses the docker
  stack defined in `tests/acceptance/docker-compose.yml`): `http.Get` against
  `ubox-crosser:8080/buildinfo`, assert status + JSON shape + non-empty
  `git_sha` and `go_version == "go1.23"`.

The acceptance docker-compose already maps `8080:8080` and sets a healthcheck,
so the integration test just piggybacks.

## Risks / trade-offs

- **Git SHA in a shallow clone**: `git rev-parse --short HEAD` in a
  `--depth=1` runner is fine (HEAD is resolvable); no unshallow needed.
- **Docker build context missing `.git`**: mitigated by passing `GIT_SHA` as
  a build arg from outside (Makefile / CI). Don't try to `COPY .git` — bloats
  the image.
- **Env read on every request**: `os.Getenv("BUILD_ID")` per request is
  ~100ns; negligible vs network. Keeps container redeploys with different
  `BUILD_ID` values picking up without restart. If profiling ever flags it,
  snapshot at `init`.
- **Hard-coded `go_version`**: simpler than `runtime.Version()` (which
  returns `"go1.23.x"` with a patch suffix the acceptance spec doesn't want).
  Trade-off: the string lies if someone bumps `go.mod` without touching this
  literal. Mitigation: a unit test compares the literal to `go.mod`'s
  `go` directive so drift fails CI.

## Parallel split assessment

**Not splitting.** Single-repo, < 100 LOC touching 3-4 files, one handler +
one unit test + one integration test. Splitting into parallel sub-issues
would be pure overhead; the critical-path wall-clock is the docker-compose
rebuild, which splitting doesn't help.
