## ADDED Requirements

### Requirement: Update mutable configuration fields
The system SHALL expose an HTTP endpoint `PUT /api/config` that accepts a JSON body to update mutable configuration fields at runtime without restarting the process.

#### Scenario: Update log level successfully
- **WHEN** a `PUT` request is sent to `/api/config` with body `{"common": {"log_level": "info"}}`
- **THEN** the response status code SHALL be `200`
- **THEN** the response body SHALL contain `{"updated": ["common.log_level"]}`
- **THEN** the running log level SHALL be changed to `info` immediately

#### Scenario: Update log level for a specific service
- **WHEN** a `PUT` request is sent to `/api/config` with body `{"service1": {"log_level": "warn"}}`
- **THEN** the response status code SHALL be `200`
- **THEN** the response body SHALL contain `{"updated": ["service1.log_level"]}`
- **THEN** the configuration for `service1` SHALL reflect the new log level

#### Scenario: Reject update of immutable field
- **WHEN** a `PUT` request is sent to `/api/config` with body `{"common": {"key": "new-key"}}`
- **THEN** the response status code SHALL be `400`
- **THEN** the response body SHALL contain an error message identifying `key` as an immutable field
- **THEN** no configuration fields SHALL be modified

#### Scenario: Reject update with invalid JSON body
- **WHEN** a `PUT` request is sent to `/api/config` with a malformed JSON body
- **THEN** the response status code SHALL be `400`
- **THEN** the response body SHALL contain an error message about invalid JSON

#### Scenario: Reject update for non-existent service
- **WHEN** a `PUT` request is sent to `/api/config` with body `{"unknown_service": {"log_level": "info"}}`
- **THEN** the response status code SHALL be `404`
- **THEN** the response body SHALL contain an error message indicating the service was not found

### Requirement: Hot-reload without restart
Configuration updates applied via `PUT /api/config` SHALL take effect immediately in the running process without requiring a restart.

#### Scenario: Log level change takes effect immediately
- **WHEN** a `PUT` request updates `log_level` to `"error"`
- **THEN** subsequent log output SHALL respect the new `error` level
- **THEN** debug and info log messages SHALL no longer appear

### Requirement: Atomicity of configuration updates
The system SHALL apply configuration updates atomically. If any field in a request is invalid or immutable, the entire request SHALL be rejected and no fields SHALL be modified.

#### Scenario: Mixed mutable and immutable fields rejected entirely
- **WHEN** a `PUT` request is sent with body `{"common": {"log_level": "info", "key": "new-key"}}`
- **THEN** the response status code SHALL be `400`
- **THEN** `log_level` SHALL NOT be updated
- **THEN** the error message SHALL identify `key` as the immutable field
