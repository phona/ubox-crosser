# Contract Testing: /version Endpoint Git SHA

## Context

This is the **contract-spec stage** for REQ-e2e-1776862572, following acceptance testing of the /version endpoint feature (REQ-rbac-1776859928).

The /version endpoint returns server metadata including the git commit SHA. This contract-spec stage formalizes cross-component test contracts to ensure the feature behaves reliably under real deployment conditions.

## Feature Summary

- **Endpoint**: `GET /version` on HTTP admin server (default :8080, configurable)
- **Response**: JSON object with commit field containing 40-character hex git SHA
- **Build Integration**: Git SHA injected via LDFLAGS at compile time
- **Related Specs**: See REQ-rbac-1776859928/specs/version_endpoint_git_sha/ for detailed scenarios

## Contract Testing Approach

The contract tests use **docker-compose based blackbox testing** to validate:

1. **Real Stack Integration**: proxy-server (ubox-crosser) running in container with actual binaries
2. **Port Configuration**: Admin server responds on configured port with correct git SHA
3. **Immutability Contract**: Commit SHA remains consistent across requests during server lifetime
4. **Security Contract**: No sensitive metadata leakage in response
5. **Response Format Contract**: JSON with specific fields matches specification

## Test Artifacts

- **Integration Tests**: `tests/integration/version_endpoint_test.go` (existing)
- **Docker Compose**: `tests/docker-compose.yml` (existing, with ADMIN_SERVER_ADDR env var)
- **Build Flags**: Makefile/Dockerfile inject COMMIT via LDFLAGS
- **Test Runner**: Runs as container service with access to proxy-server:8080

## Success Criteria

✓ All version endpoint tests pass in docker-compose environment
✓ Contract scenarios (FEATURE-A1 through FEATURE-A7) covered by test code
✓ Tests are RED before dev implementation, GREEN after
✓ No inline mocks or httptest.NewServer - real stack only
✓ Environment variables properly injected (ADMIN_SERVER_ADDR, CUSTOM_ADMIN_ADDR)

## Lifecycle

- **Created**: Contract-spec agent (this stage)
- **Owner**: Contract-spec agent (LOCKED - dev-agent cannot modify tests/contract/*)
- **Validation**: `docker compose -f tests/docker-compose.yml up --build --exit-code-from test-runner`
- **Integration**: CI runs same command as ci-int job
