## ADDED Requirements

### Requirement: Health endpoint returns OK on GET

The proxy server SHALL expose an HTTP endpoint at the path `/health` that returns HTTP status 200 with a JSON body `{"status":"ok"}` and `Content-Type: application/json` header.

#### Scenario: Successful health check

- **WHEN** a client sends `GET /health` to the health HTTP listener
- **THEN** the server responds with HTTP 200, header `Content-Type: application/json`, and body `{"status":"ok"}`

### Requirement: Health endpoint rejects non-GET methods

The health endpoint SHALL return HTTP 405 Method Not Allowed for any HTTP method other than GET on the `/health` path. The response SHALL include an `Allow: GET` header.

#### Scenario: POST request to health endpoint

- **WHEN** a client sends `POST /health` to the health HTTP listener
- **THEN** the server responds with HTTP 405 and an `Allow: GET` header

#### Scenario: PUT request to health endpoint

- **WHEN** a client sends `PUT /health` to the health HTTP listener
- **THEN** the server responds with HTTP 405 and an `Allow: GET` header

### Requirement: Unknown paths return 404

The health HTTP listener SHALL return HTTP 404 Not Found for any request path other than `/health`.

#### Scenario: Request to root path

- **WHEN** a client sends `GET /` to the health HTTP listener
- **THEN** the server responds with HTTP 404

#### Scenario: Request to unknown path

- **WHEN** a client sends `GET /metrics` to the health HTTP listener
- **THEN** the server responds with HTTP 404

### Requirement: Health address is configurable

The server SHALL accept a `health_address` configuration field (JSON key: `"health_address"`) and a `--health-address` CLI flag specifying the `host:port` on which the health HTTP listener binds.

#### Scenario: Health address set via config file

- **WHEN** the server config file contains `"health_address": ":9090"` in the common section
- **THEN** the health HTTP listener binds to `:9090`

#### Scenario: Health address set via CLI flag

- **WHEN** the server is started with `--health-address :9090`
- **THEN** the health HTTP listener binds to `:9090`

### Requirement: Health server disabled when address is empty

The server SHALL NOT start the health HTTP listener when `health_address` is empty or unset.

#### Scenario: No health address configured

- **WHEN** the server starts without a `health_address` value in config or CLI
- **THEN** no health HTTP listener is started and the proxy server operates normally

### Requirement: Health server startup errors are reported

The server SHALL report health HTTP listener startup failures (e.g., port already in use) through the existing error channel so they surface in the server's error logging.

#### Scenario: Port conflict on startup

- **WHEN** the configured `health_address` port is already in use
- **THEN** the health server reports the bind error through the server's error channel
