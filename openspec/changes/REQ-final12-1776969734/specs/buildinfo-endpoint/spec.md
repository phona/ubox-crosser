## ADDED Requirements

### Requirement: The system SHALL expose /buildinfo returning git_sha, build_id, and go_version

The `/buildinfo` HTTP endpoint SHALL be available on the management server and SHALL return build metadata
as a JSON object. No authentication SHALL be required. The endpoint SHALL be served on the same HTTP server
as `/healthz` (default `:8080`).

#### Scenario: UBOX-S1 returns 200 with all three fields on GET

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
