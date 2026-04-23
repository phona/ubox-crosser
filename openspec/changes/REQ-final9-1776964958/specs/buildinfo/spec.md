## ADDED Requirements

### Requirement: The system SHALL expose /buildinfo returning git_sha, build_id, and go_version

The server SHALL provide a `/buildinfo` endpoint over HTTP that returns a JSON object with
three fields: `git_sha` (7-char git commit SHA injected via ldflags, empty string if not set),
`build_id` (value of `BUILD_ID` environment variable, defaulting to `"dev"`), and
`go_version` (hardcoded to `"go1.23"`). The endpoint SHALL require no authentication.

#### Scenario: UBOX-S1 returns 200 with all three fields on bare GET
- **GIVEN** the server is started (BUILD_ID env may be unset)
- **WHEN** a client sends GET /buildinfo
- **THEN** the response status is 200
- **AND** the response body is valid JSON containing fields git_sha, build_id, and go_version

#### Scenario: UBOX-S2 build_id defaults to "dev" when BUILD_ID env is unset
- **GIVEN** the server is started without the BUILD_ID environment variable
- **WHEN** a client sends GET /buildinfo
- **THEN** the response JSON field build_id equals "dev"

#### Scenario: UBOX-S3 build_id reflects BUILD_ID env when set
- **GIVEN** the server is started with BUILD_ID=42
- **WHEN** a client sends GET /buildinfo
- **THEN** the response JSON field build_id equals "42"

#### Scenario: UBOX-S4 go_version is always "go1.23"
- **GIVEN** the server is running
- **WHEN** a client sends GET /buildinfo
- **THEN** the response JSON field go_version equals "go1.23"

#### Scenario: UBOX-S5 endpoint requires no authentication
- **GIVEN** the server is running
- **WHEN** a client sends GET /buildinfo without any credentials
- **THEN** the response status is 200 (not 401 or 403)

#### Scenario: UBOX-S6 endpoint is served on the health-check port alongside /healthz
- **GIVEN** the server is started with --health-addr=:8080
- **WHEN** a client sends GET /buildinfo to port 8080
- **THEN** the response status is 200
