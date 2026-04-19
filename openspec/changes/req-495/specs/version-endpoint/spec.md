## ADDED Requirements

### Requirement: Version HTTP endpoint
The system SHALL expose an HTTP `GET /version` endpoint on a configurable address that returns the current build version information in JSON format.

#### Scenario: Successful version query
- **WHEN** a client sends `GET /version` to the HTTP listen address
- **THEN** the server SHALL respond with HTTP 200 and a JSON body containing `version`, `commit`, and `buildTime` fields

#### Scenario: Version injected at build time
- **WHEN** the binary is built with `-ldflags "-X version.Version=v1.0.0 -X version.Commit=abc123 -X version.BuildTime=2026-04-19T00:00:00Z"`
- **THEN** the `GET /version` response SHALL reflect those exact values

#### Scenario: Default values when not injected
- **WHEN** the binary is built without ldflags version injection
- **THEN** the `version` field SHALL be `"dev"`, `commit` SHALL be `"unknown"`, and `buildTime` SHALL be `"unknown"`

### Requirement: HTTP listener configuration
The proxy server SHALL accept an optional HTTP listen address via configuration (config file field or CLI flag). The HTTP server SHALL run concurrently with the existing TCP proxy without blocking it.

#### Scenario: HTTP address specified
- **WHEN** the server is started with HTTP address `:8080`
- **THEN** the HTTP server SHALL listen on port 8080 and serve the `/version` endpoint

#### Scenario: HTTP address not specified
- **WHEN** the server is started without an HTTP address
- **THEN** no HTTP server SHALL be started, and the proxy SHALL operate as before

### Requirement: Version response format
The `GET /version` response body SHALL be a JSON object with the following structure:
```json
{
  "version": "<string>",
  "commit": "<string>",
  "buildTime": "<string>"
}
```
The `Content-Type` header SHALL be `application/json`.

#### Scenario: Response content type
- **WHEN** a client sends `GET /version`
- **THEN** the response `Content-Type` header SHALL be `application/json`

### Requirement: Build system version injection
The Makefile SHALL inject version, commit hash, and build time into the binary via `-ldflags -X` flags. The `VERSION` variable SHALL default to the output of `git describe --tags --always` if available, falling back to `"dev"`.

#### Scenario: Build with git tag
- **WHEN** `make build` is run in a repository with tag `v1.2.3` on HEAD
- **THEN** the built binary SHALL report `version` as `v1.2.3`

#### Scenario: Build without git tag
- **WHEN** `make build` is run in a repository without tags
- **THEN** the built binary SHALL report `version` based on `git describe --always` (short commit hash)
