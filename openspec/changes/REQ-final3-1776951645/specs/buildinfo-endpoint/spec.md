# /buildinfo Endpoint — Build Identification

## Overview

The `ubox-crosser server` binary exposes an unauthenticated HTTP endpoint
`GET /buildinfo` on its admin HTTP listener (port `8080`). It returns a
fixed-shape JSON document identifying the build — git SHA, build
identifier, and Go toolchain version — so operators and CI can verify at
a glance which build is running in a given container or cluster.

## Acceptance scenarios

### [UBOX-CROSSER-S1] `/buildinfo` returns 200 and JSON

**Given** the ubox-crosser server is running and its admin HTTP listener
is bound to port `8080`  
**When** a client issues `GET http://<server>:8080/buildinfo`  
**Then** the response status is `200 OK`  
**And** the `Content-Type` header is `application/json`  
**And** the response body parses as a JSON object.

### [UBOX-CROSSER-S2] Response shape is exactly three fields

**Given** the server is running  
**When** a client issues `GET /buildinfo`  
**Then** the JSON body has exactly the keys `git_sha`, `build_id`, and
`go_version`  
**And** every value is a non-empty string.

### [UBOX-CROSSER-S3] `git_sha` reflects the ldflag-injected commit

**Given** the server binary was built with
`-ldflags "-X main.GitSHA=<7-char sha>"`  
**When** a client issues `GET /buildinfo`  
**Then** the response's `git_sha` field equals that 7-character SHA  
**And** when no ldflag was passed, `git_sha` equals the literal
`"unknown"`.

### [UBOX-CROSSER-S4] `build_id` respects the `BUILD_ID` env var

**Given** the server process has `BUILD_ID=ci-12345` in its environment  
**When** a client issues `GET /buildinfo`  
**Then** the response's `build_id` field equals `"ci-12345"`  
**And** when `BUILD_ID` is unset or empty, `build_id` equals the literal
`"dev"`  
**And** updating `BUILD_ID` without restarting the process is reflected
on the next request (read-on-each-request semantics).

### [UBOX-CROSSER-S5] `go_version` is the literal `"go1.23"`

**Given** the server is running  
**When** a client issues `GET /buildinfo`  
**Then** the response's `go_version` field is exactly the string
`"go1.23"` — not `runtime.Version()`, not a patch-level string.

## Non-functional requirements

- The endpoint is unauthenticated. It exposes only build-identifying
  metadata, no secrets, no PII, no runtime state beyond what an operator
  picks via `BUILD_ID`.
- The endpoint responds in O(1) with no filesystem or network I/O
  (one `os.Getenv` read per request).
- Concurrent requests are supported — no shared mutable state in the
  handler.
- The admin listener must not block startup of the TCP proxy loop; it
  runs in a separate goroutine.
