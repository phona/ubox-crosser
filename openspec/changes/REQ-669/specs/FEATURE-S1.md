---
id: FEATURE-S1
title: Version endpoint returns correct JSON
---

# Spec: GET /version returns build metadata

## Scenario

Given the server is running with API enabled,
When a client sends `GET /version`,
Then the response has status 200, Content-Type `application/json`, and body contains `version`, `commit`, `build_time` fields.

## Acceptance Criteria

1. HTTP 200 status code
2. `Content-Type: application/json` header
3. Response body is valid JSON with exactly three string fields: `version`, `commit`, `build_time`
4. `version` matches the compiled constant (default "0.1.0")
5. Non-GET methods receive 405 Method Not Allowed
