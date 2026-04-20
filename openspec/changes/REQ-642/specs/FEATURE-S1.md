---
spec_id: FEATURE-S1
change_id: REQ-642
title: "GET /version returns 200 with version info"
layer: backend
---

# FEATURE-S1: GET /version returns 200 with version info

## Scenario

**Given** the server is running with HTTP enabled
**When** a client sends `GET /version`
**Then** the server responds with HTTP 200, `Content-Type: application/json`, and a body containing `version`, `commit`, and `build_time` fields.

## Acceptance Criteria

- Status code is 200
- Content-Type header is `application/json`
- Response body is valid JSON with exactly three fields: `version` (string), `commit` (string), `build_time` (string)
- `version` matches the package constant
- Non-GET methods return 405
