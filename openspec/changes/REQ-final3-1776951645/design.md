# REQ-final3-1776951645: Design Notes

## High-level approach

A single `http.HandlerFunc` in `cmd/server/buildinfo.go` marshals a
`BuildInfo` struct to JSON. The handler is registered on a new admin
`*http.ServeMux` that `cmd/server/server.go` starts in a goroutine on
port `8080` before the proxy loop begins blocking. The admin listener is
deliberately minimal — no TLS, no auth, no middleware — because it is a
read-only identification endpoint and the spec explicitly says it is
unauthenticated.

`GitSHA` lives as a package-level `var` in `main` (so `-X main.GitSHA=…`
lands it directly) and stays `"unknown"` when the ldflag is absent (local
`go run` / `go test`). `BUILD_ID` is read via `os.Getenv` on each request
so operators can redeploy with a new `BUILD_ID` without a container
restart. `go_version` is a hard-coded literal `"go1.23"` (not
`runtime.Version()` — that returns `"go1.23.12"` or similar, and the
acceptance spec wants a stable `"go1.23"` value).

## Data model

```go
package main

var GitSHA = "unknown" // injected via -ldflags -X main.GitSHA=…

type BuildInfo struct {
    GitSHA    string `json:"git_sha"`
    BuildID   string `json:"build_id"`
    GoVersion string `json:"go_version"`
}
```

## Admin HTTP listener

New file `server/admin_http.go`:

```go
package server

// StartAdminHTTP starts a minimal admin HTTP listener and blocks.
// Callers register handlers on mux before calling.
func StartAdminHTTP(addr string, mux *http.ServeMux) error {
    s := &http.Server{
        Addr:              addr,
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
    }
    return s.ListenAndServe()
}
```

`cmd/server/server.go` wires it:

```go
mux := http.NewServeMux()
mux.HandleFunc("/buildinfo", BuildInfoHandler)
go func() {
    if err := server.StartAdminHTTP(":8080", mux); err != nil {
        logrus.Errorf("admin http: %v", err)
    }
}()
```

This keeps the mux construction in `main()` so `/healthz` (from
REQ-e2e-1776916220) can join the same mux on merge without touching the
`server/` package.

## Build-time injection

`Makefile`:

```make
GIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
SERVER_LDFLAGS := -s -w -X main.GitSHA=$(GIT_SHA)
```

`build` and `ci-build` targets pass `-ldflags="$(SERVER_LDFLAGS)"` only
for `./cmd/server` (client / auth_server don't need it; their main
packages don't declare `GitSHA`).

`Dockerfile` accepts `ARG GIT_SHA=unknown` and threads it into the
`go build -ldflags` invocation. Docker doesn't carry `.git` into the
build context reliably (and `.dockerignore` may strip it), so passing the
SHA as a build arg from outside is the only safe option.

`tests/Dockerfile.test` also accepts `ARG GIT_SHA=unknown` and threads it
into the server build so the integration test can assert a non-empty
`git_sha` against the baked-in image.

## Integration with docker-compose

`tests/docker-compose.yml` changes to the `proxy-server` service:

- Add `environment: BUILD_ID`.
- Add `build.args: GIT_SHA: "${GIT_SHA:-unknown}"`.

The test-runner does not need a port mapping — it is inside the same
docker network and addresses the admin listener as `proxy-server:8080`.

## Testing

**Unit** (`cmd/server/buildinfo_test.go`, default build):

- `httptest.NewRecorder()` + direct handler call.
- Assert status 200, `Content-Type: application/json`, valid JSON.
- Assert `git_sha` defaults to `"unknown"` when `GitSHA` package var is
  not overridden; override via `GitSHA = "abc1234"` in the test and
  assert it.
- Use `t.Setenv("BUILD_ID", "test-123")` and assert `build_id` matches.
- Assert `go_version == "go1.23"` always.

**Integration** (`tests/integration/buildinfo_test.go`,
`//go:build integration`):

- `http.Get("http://" + proxyAddr + "/buildinfo")` where proxyAddr is
  derived from `PROXY_SERVER_ADDR` env (stripping the `:7000` proxy port
  and appending `:8080`).
- Status 200.
- JSON unmarshals into `BuildInfo`.
- `GoVersion == "go1.23"`.
- `GitSHA != ""` and `GitSHA != "unknown"` (proves ldflag threading).
- `BuildID == "ci-test"` (set via docker-compose env).

## Risks / trade-offs

- **Shallow clones**: `git rev-parse --short HEAD` on a `--depth=1` clone
  still resolves HEAD fine — no unshallow needed.
- **.git absent from Docker build context**: explicitly mitigated by the
  `GIT_SHA` build arg. Do not `COPY .git` — image bloat + breaks on
  builders that strip it.
- **Env read per request**: `os.Getenv("BUILD_ID")` is ~100ns, dwarfed by
  network; benefit is operators can change `BUILD_ID` without restart.
- **Hard-coded go_version**: diverges from `runtime.Version()` if the
  toolchain is bumped. A unit test cross-checks the literal against the
  `go` directive in `go.mod` to fail CI on drift.
- **Port 8080 collision**: If the host binds something else to 8080 in
  production, the admin listener fails. For this REQ the server logs
  the error and keeps the proxy loop alive — admin endpoints are
  optional for the proxy's primary function.

## Parallel split assessment

**Not split.** Single repo, < 150 LOC, one handler + one unit test +
one integration test + two build-infra edits (Makefile + Dockerfile).
Fan-out overhead would exceed the whole implementation cost. One
linear pass.
