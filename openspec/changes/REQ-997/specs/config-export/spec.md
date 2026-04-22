---
capability: config-export
change_id: REQ-997
status: LOCKED
---

## ADDED

### Scenario: FEATURE-A13 — Config export produces server.json-compatible output

```gherkin
Given the crosser-api server is running
  and the admin is authenticated with a valid JWT
  and a service named "test-svc" exists with key, method, address, login_password, and auth_password set
When the admin sends a GET request to /api/v1/services/test-svc/config
Then the response status code is 200
  and the response JSON contains a "common" object with fields: key, method, address, login_password, auth_password, log_file, log_level
  and the response JSON contains a "test-svc" object with field: key
  and the JSON structure is byte-level compatible with the existing server.json format
```

### Scenario: FEATURE-A14 — Config export for non-existent service returns 404

```gherkin
Given the crosser-api server is running
  and the admin is authenticated with a valid JWT
When the admin sends a GET request to /api/v1/services/nonexistent/config
Then the response status code is 404
  and the response JSON field "code" is not 0
```

### Scenario: FEATURE-A15 — Exported config can be consumed by existing proxy-server

```gherkin
Given the crosser-api server is running
  and a service named "test-svc" exists with full configuration
  and the admin exports the config via GET /api/v1/services/test-svc/config
When the exported JSON is written to a file and fed to the existing proxy-server as server.json
Then the proxy-server starts without configuration errors
```
