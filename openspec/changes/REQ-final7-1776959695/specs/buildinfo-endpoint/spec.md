# /buildinfo Endpoint — Build Metadata

## ADDED Requirements

### Requirement: The system SHALL expose GET /buildinfo returning git_sha, build_id, and go_version

#### Scenario: UBOX-BINFO-S1 returns 200 with three required JSON fields on bare GET
- **GIVEN** the server is running (BUILD_ID env may or may not be set)
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response status is 200
- **AND** the response body is valid JSON containing `git_sha`, `build_id`, and `go_version`

#### Scenario: UBOX-BINFO-S2 build_id reflects BUILD_ID environment variable
- **GIVEN** the server is started with `BUILD_ID=ci-run-42` in its environment
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `build_id` field equals `"ci-run-42"`

#### Scenario: UBOX-BINFO-S3 build_id defaults to "dev" when BUILD_ID is absent
- **GIVEN** the server is started without a `BUILD_ID` environment variable
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `build_id` field equals `"dev"`

#### Scenario: UBOX-BINFO-S4 go_version is always "go1.23"
- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `go_version` field equals `"go1.23"`

#### Scenario: UBOX-BINFO-S5 git_sha reflects the value injected via ldflags
- **GIVEN** the server binary was built with `-X main.GitSHA=abc1234`
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `git_sha` field equals `"abc1234"`

#### Scenario: UBOX-BINFO-S6 endpoint is available without authentication
- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo` with no credentials
- **THEN** the response status is 200 (no 401/403)

#### Scenario: UBOX-BINFO-S7 endpoint uses Content-Type application/json
- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the `Content-Type` response header is `application/json`
