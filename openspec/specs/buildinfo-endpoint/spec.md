# buildinfo-endpoint Specification

## Purpose
TBD - created by archiving change REQ-final15-1776989948. Update Purpose after archive.
## Requirements
### Requirement: /buildinfo endpoint exposes build metadata

The system SHALL expose `GET /buildinfo` on the management HTTP server. The endpoint MUST
return a JSON object containing exactly three fields: `git_sha` (7-char short SHA injected
at build time via `-ldflags "-X main.GitSHA=..."`), `build_id` (value of the `BUILD_ID`
environment variable at server startup, defaulting to `"dev"` when absent), and `go_version`
(the literal string `"go1.23"`). The endpoint MUST NOT require authentication.

#### Scenario: UBOX-S1 returns 200 with all three fields on bare GET

- **GIVEN** the server is started with GitSHA injected via ldflags and BUILD_ID unset (defaults to "dev")
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response status is 200
- **AND** the response body is JSON containing `git_sha`, `build_id`, and `go_version`

#### Scenario: UBOX-S2 git_sha reflects the value injected at build time

- **GIVEN** the server binary was built with `-ldflags "-X main.GitSHA=abc1234"`
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `git_sha` field in the response equals `"abc1234"`

#### Scenario: UBOX-S3 build_id reflects BUILD_ID env when set

- **GIVEN** the server is started with environment variable `BUILD_ID=ci-42`
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `build_id` field in the response equals `"ci-42"`

#### Scenario: UBOX-S4 build_id defaults to "dev" when BUILD_ID env is absent

- **GIVEN** the server is started without the `BUILD_ID` environment variable
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `build_id` field in the response equals `"dev"`

#### Scenario: UBOX-S5 go_version is always "go1.23"

- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `go_version` field in the response equals `"go1.23"`

#### Scenario: UBOX-S6 response Content-Type is application/json

- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `Content-Type` response header contains `application/json`

#### Scenario: UBOX-S7 endpoint requires no authentication

- **GIVEN** the server is running without any auth headers or tokens
- **WHEN** a client sends `GET /buildinfo` with no Authorization header
- **THEN** the response status is 200

