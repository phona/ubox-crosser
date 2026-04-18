## Why

The ubox-crosser project currently only exposes a TCP/SOCKS5-based protocol with no HTTP API surface. Adding a simple HTTP endpoint enables health-check tooling, debugging, and lays the groundwork for future management APIs. Starting with `/api/echo` provides a minimal, verifiable first endpoint.

## What Changes

- Add an HTTP server capability to the project
- Expose `GET /api/echo?msg=X` that returns a JSON response echoing the provided message
- The HTTP server runs alongside the existing TCP/SOCKS5 services

## Capabilities

### New Capabilities
- `echo-api`: HTTP endpoint that accepts a query parameter and returns it in a JSON response

### Modified Capabilities
<!-- None — this is a new, additive capability with no changes to existing behavior. -->

## Impact

- **New dependency**: Go stdlib `net/http` (no external dependencies needed)
- **New code**: HTTP handler, server setup, and routing
- **Configuration**: HTTP listen address/port configuration
- **Existing systems**: No impact on existing SOCKS5 proxy functionality
