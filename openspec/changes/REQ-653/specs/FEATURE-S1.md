---
id: FEATURE-S1
title: GET /version returns build metadata
---

## Scenario: GET /version returns 200 with JSON build info

**Given** the server is running with admin HTTP enabled
**When** a client sends `GET /version`
**Then** the response status is `200 OK`
**And** the `Content-Type` header is `application/json`
**And** the body contains `version`, `commit`, and `build_time` fields
**And** `version` equals the hardcoded constant (currently `0.1.0`)
