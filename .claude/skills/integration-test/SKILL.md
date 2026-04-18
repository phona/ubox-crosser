---
name: integration-test
description: Write integration tests for ubox-crosser components. Trigger when user asks to create, add, or write integration tests, or mentions testing the proxy tunnel, auth flow, or end-to-end connectivity. These are Docker Compose-based blackbox tests.
allowed-tools: Read, Glob, Grep, Write, Edit, TodoWrite, AskUserQuestion, Bash
---

# Integration Test Writing Skill

## Overview

This skill guides writing blackbox integration tests for the ubox-crosser reverse proxy system.
Tests are organized **by test scenario** and run inside Docker Compose with all three components (ProxyServer, Client, AuthServer) plus an echo server.

Full naming convention and rules: see [rules.md](rules.md)
Templates: [templates/tunnel.go.tmpl](templates/tunnel.go.tmpl)

---

## Step 1 — Identify Target (MANDATORY)

Before writing any test, identify the scenario:

1. **Which component(s) are being tested?**
   - ProxyServer (login, heartbeat, worker management)
   - Client (SOCKS5 tunneling, reconnection)
   - AuthServer (authentication, bridging)
   - End-to-end (full tunnel chain)

2. **What protocol interaction is being tested?**
   - LOGIN message flow
   - HEART_BEAT keep-alive
   - GEN_WORKER connection pooling
   - AUTHENTICATION gateway flow
   - SOCKS5 proxying
   - Encryption/decryption

---

## Step 2 — Choose Test Strategy

Ask the user:

```yaml
AskUserQuestion:
  Q1: Which scenario to test?
    options:
      - Full tunnel (external → auth-server → proxy → client → target)
      - Protocol level (login, heartbeat, worker generation)
      - Error handling (wrong password, invalid service, timeout)
      - Encryption (cipher mismatch, unencrypted fallback)
  Q2: Priority level?
    options:
      - P0 only (critical path)
      - P0 + P1 (core + validation)
      - All: P0 + P1 + P2 (full coverage)
```

---

## Step 3 — Generate Test File

### Convention Checklist

1. ✅ Build tag: `//go:build integration` on line 1
2. ✅ Package: `integration`
3. ✅ All addresses from environment variables (PROXY_SERVER_ADDR, AUTH_SERVER_ADDR, etc.)
4. ✅ Cipher config from env (TEST_CIPHER_KEY, TEST_CIPHER_METHOD)
5. ✅ Each test is self-contained — no shared mutable state
6. ✅ Timeouts on all network operations
7. ✅ Use `dialEncrypted()` helper for Shadowsocks connections
8. ✅ Use separate service names for protocol tests vs tunnel tests (avoid controller conflicts)

### Anti-patterns to REJECT

- ❌ Hardcoded IP addresses or ports
- ❌ Tests that share the same service name and interfere with each other's controllers
- ❌ Missing `conn.Close()` / `defer conn.Close()`
- ❌ `time.Sleep` without corresponding timeout on reads
- ❌ Ignoring `SetReadDeadline` errors

---

## Step 4 — Run and Verify

```bash
# Full integration test via Docker Compose
make test-integration

# Or run specific test
docker compose -f tests/docker-compose.yml up --build --exit-code-from test-runner
```

### Test Environment

Docker Compose provides:
- `echo-server` — HTTP echo target on port 9000
- `proxy-server` — Central proxy relay on port 7000
- `client` — NAT client connected to proxy-server
- `auth-server` — Public gateway on port 7002
- `test-runner` — Executes `go test -tags integration`

Environment variables available in test-runner:
- `PROXY_SERVER_ADDR=proxy-server:7000`
- `AUTH_SERVER_ADDR=auth-server:7002`
- `ECHO_SERVER_ADDR=echo-server:9000`
- `TEST_CIPHER_KEY`, `TEST_CIPHER_METHOD`
- `TEST_LOGIN_PASSWORD`, `TEST_AUTH_PASSWORD`
- `TEST_SERVE_NAME`
