# REQ-final5-1776956833: Add /buildinfo HTTP Endpoint

## Summary

Expose build information via an HTTP `/buildinfo` endpoint that returns JSON containing git commit SHA, build ID, and Go version. This enables monitoring and deployment verification.

## Motivation

- Enable quick verification of deployed version
- Support CI/CD pipelines that need to confirm correct build was deployed
- Provide diagnostics endpoint for troubleshooting

## Scope

### In Scope
- Add HTTP server alongside existing proxy server
- Expose `/buildinfo` endpoint returning JSON with:
  - `git_sha`: 7-character git commit SHA
  - `build_id`: Environment variable `BUILD_ID` or "dev"
  - `go_version`: hardcoded "go1.23"
- No authentication required
- Use existing configured port

### Out of Scope
- Additional diagnostics endpoints
- Metrics or monitoring
- Authentication mechanisms

## Implementation Approach

1. Add HTTP server to the proxy server startup
2. Register `/buildinfo` handler that returns structured JSON
3. Use ldflags to inject git SHA at build time
4. Read `BUILD_ID` environment variable at runtime
5. Add comprehensive unit and integration tests

## Acceptance Criteria

- ✅ Endpoint returns 200 HTTP status
- ✅ Response is valid JSON with three required fields
- ✅ git_sha is 7 characters
- ✅ build_id defaults to "dev" when env var not set
- ✅ go_version is always "go1.23"
- ✅ Works via curl: `curl http://server:<port>/buildinfo`
- ✅ Unit tests cover JSON structure and field values
- ✅ Integration tests verify HTTP response

