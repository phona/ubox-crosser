## ADDED Requirements

### Requirement: Retrieve running configuration
The system SHALL expose an HTTP endpoint `GET /api/config` that returns the current running configuration of all services as a JSON object.

#### Scenario: Successful config retrieval
- **WHEN** a `GET` request is sent to `/api/config`
- **THEN** the response status code SHALL be `200`
- **THEN** the response `Content-Type` header SHALL be `application/json`
- **THEN** the response body SHALL be a JSON object with service names as keys and their configuration objects as values

#### Scenario: Sensitive fields are masked
- **WHEN** a `GET` request is sent to `/api/config`
- **THEN** the fields `key`, `login_password`, and `auth_password` in each service configuration SHALL have their values replaced with `"***"`

#### Scenario: Method not allowed
- **WHEN** a `POST` request is sent to `/api/config`
- **THEN** the response status code SHALL be `405`

### Requirement: Management server lifecycle
The management HTTP server SHALL start only when `management_address` is configured in the server's common configuration. When not configured, no HTTP listener SHALL be created.

#### Scenario: Management server starts when configured
- **WHEN** the server configuration contains `"management_address": "127.0.0.1:8080"` in the common section
- **THEN** the management HTTP server SHALL listen on `127.0.0.1:8080`

#### Scenario: Management server does not start when unconfigured
- **WHEN** the server configuration does not contain a `management_address` field
- **THEN** no management HTTP listener SHALL be created
- **THEN** the proxy server SHALL operate normally without the management API
