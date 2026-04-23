## ADDED Requirements

### Requirement: The server SHALL expose a /buildinfo endpoint returning git_sha, build_id, and go_version

The `/buildinfo` HTTP endpoint MUST be available on the server's HTTP management port (same port as `/healthz`). It MUST return a JSON object with exactly three fields: `git_sha` (7-character git commit hash injected at build time), `build_id` (value of `$BUILD_ID` env var, defaulting to `"dev"`), and `go_version` (hardcoded `"go1.23"`). The endpoint SHALL require no authentication.

#### Scenario: UBOX-S1 returns 200 with all three fields on bare GET

- **GIVEN** the server is started with default configuration
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response status is 200
- **AND** the response body is valid JSON containing `git_sha`, `build_id`, and `go_version` fields

#### Scenario: UBOX-S2 build_id defaults to dev when BUILD_ID env var is unset

- **GIVEN** the server is started without `BUILD_ID` set in the environment
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response JSON has `"build_id": "dev"`

#### Scenario: UBOX-S3 build_id reflects BUILD_ID env var when set

- **GIVEN** the server is started with `BUILD_ID=ci-run-42` in the environment
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response JSON has `"build_id": "ci-run-42"`

#### Scenario: UBOX-S4 go_version is always go1.23

- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response JSON has `"go_version": "go1.23"`

#### Scenario: UBOX-S5 git_sha matches the 7-char short hash from build ldflags

- **GIVEN** the server binary was built with `-ldflags "-X main.GitSHA=abc1234"`
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response JSON has `"git_sha": "abc1234"`

#### Scenario: UBOX-S6 response Content-Type is application/json

- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response `Content-Type` header is `application/json`

#### Scenario: UBOX-S7 endpoint is unauthenticated

- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo` without any credentials or tokens
- **THEN** the response status is 200 (not 401 or 403)
