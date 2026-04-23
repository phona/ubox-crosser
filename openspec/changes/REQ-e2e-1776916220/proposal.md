# Proposal: Add /healthz Endpoint for Service Uptime Monitoring

## Problem Statement

The ProxyServer lacks observability for service uptime and health status. Operational monitoring requires a simple, standard endpoint to check if the service is running and how long it has been up.

## Solution

Add an HTTP `/healthz` endpoint to the ProxyServer that returns:
- Service health status
- Uptime in seconds
- Human-readable duration string

This endpoint runs on a dedicated HTTP port and does not interfere with existing message-based TCP communication.

## Scope

### In Scope
- Add HTTP listener to ProxyServer (port configurable via `PROXY_HTTP_PORT`, default 8080)
- Implement `/healthz` endpoint returning JSON with status and uptime
- Track service startup time and calculate uptime
- Contract tests verifying endpoint behavior
- OpenAPI specification of the endpoint

### Out of Scope
- Metrics/Prometheus integration
- Other health checks (connection counts, latency, etc.)
- Authentication/authorization for /healthz

## Success Criteria

1. ProxyServer exposes `/healthz` endpoint on HTTP port 8080
2. Response includes JSON with `status`, `uptime.seconds`, `uptime.duration`
3. Contract tests pass verifying response format and uptime tracking
4. HTTP listener runs non-blocking and doesn't interfere with core proxy functions
5. Specification locked in as contract for future implementations

## Timeline

- contract-spec-agent: Write contract tests and OpenAPI spec (this stage)
- dev-agent: Implement HTTP endpoint and uptime tracking
- accept-agent: Verify implementation matches contract

## References

- REQ: REQ-e2e-1776916220
- Service: ProxyServer
- Related Standards: Kubernetes /healthz convention
