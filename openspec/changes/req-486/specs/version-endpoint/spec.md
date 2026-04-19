## ADDED Requirements

### Requirement: Version HTTP endpoint
proxy server SHALL provide a `GET /version` HTTP endpoint that returns the current build version information as JSON.

The response SHALL have Content-Type `application/json` and HTTP status 200, with the following fields:
- `version`: the release version string (e.g. `"2.1.0"`)
- `git_commit`: the short git commit hash at build time
- `build_time`: the build timestamp in RFC 3339 format

If version information was not injected at build time, `version` SHALL default to `"dev"`.

#### Scenario: Successful version query
- **WHEN** a client sends `GET /version` to the HTTP listener
- **THEN** the server responds with HTTP 200 and a JSON body containing `version`, `git_commit`, and `build_time` fields

#### Scenario: Version not injected at build time
- **WHEN** the binary was built without `-ldflags` version injection
- **THEN** the `version` field SHALL be `"dev"`, and `git_commit` and `build_time` SHALL be empty strings

### Requirement: HTTP listener configuration
The proxy server SHALL accept an optional `http_address` configuration field. When `http_address` is non-empty, the server SHALL start an HTTP listener on that address. When `http_address` is empty or omitted, the server SHALL NOT start any HTTP listener.

#### Scenario: HTTP listener enabled
- **WHEN** `http_address` is set to `":8080"` in the server configuration
- **THEN** the server starts an HTTP listener on port 8080 serving the `/version` endpoint

#### Scenario: HTTP listener disabled by default
- **WHEN** `http_address` is not specified in the configuration
- **THEN** no HTTP listener is started and the server operates as before

### Requirement: Build-time version injection
The Makefile SHALL inject version, git commit, and build time into the binary via `-ldflags -X` flags. The version SHALL be derived from `git describe --tags --always`.

#### Scenario: Build with version injection
- **WHEN** running `make build`
- **THEN** the resulting binaries contain the git-derived version, commit hash, and build timestamp
