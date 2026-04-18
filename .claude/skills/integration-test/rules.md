# Integration Test Naming Convention & Rules

## Test ID Structure

```
Test{Component}{Action}{Scenario}
```

### Component

| Component | Description | Connection |
|-----------|------------|------------|
| `ProxyServer` | Central relay | Encrypted TCP :7000 |
| `Client` | NAT-traversal client | Connects to proxy |
| `AuthServer` | Public gateway | Plain TCP :7002 → encrypted to proxy |
| `Tunnel` | Full end-to-end chain | auth → proxy → client → target |

### Action

`Login` · `Heartbeat` · `Worker` · `Auth` · `Bridge` · `Tunnel` · `Reconnect`

### Scenario

| Category | Scenario suffixes | Priority |
|----------|------------------|----------|
| Success | `HappyPath`, `Success` | P0 |
| Auth failure | `WrongPassword`, `InvalidServeName` | P0 |
| Network | `Timeout`, `Disconnect`, `Reconnect` | P1 |
| Encryption | `CipherMismatch`, `NoCipher` | P1 |
| Edge cases | `ConcurrentWorkers`, `LargePayload` | P2 |

---

## File Organization

```
tests/integration/
    tunnel_test.go       # full tunnel + protocol tests
    reconnect_test.go    # reconnection and recovery tests
    encryption_test.go   # cipher-related tests
    concurrent_test.go   # concurrent connection tests
```

---

## Mandatory Rules

1. `//go:build integration` must be the first line
2. All network addresses from environment variables
3. Use separate service names for protocol-level tests (`protocol_test`) vs tunnel tests (`test_service`) to avoid controller conflicts
4. Every `net.Conn` must have `defer conn.Close()`
5. Every read operation must have `SetReadDeadline`
6. Use `dialEncrypted()` helper for Shadowsocks connections to proxy

---

## Key Helpers (defined in tunnel_test.go)

```go
// Get environment variable with fallback
func getEnv(key, fallback string) string

// Create Shadowsocks cipher from test config
func newCipher() *ss.Cipher

// Dial proxy-server with Shadowsocks encryption
func dialEncrypted(addr string) (net.Conn, error)
```

---

## Known Traps

### Controller conflict
The ProxyServer's `controllers` map is keyed by `serveName`. If a test logs in as `test_service`, it **overwrites** the real client's controller. Use `protocol_test` service name for login/auth protocol tests.

### Healthcheck interference
The proxy-server healthcheck creates TCP connections that cause EOF errors in logs. This is normal and doesn't affect tests.

### SOCKS5 through AuthServer
When connecting through auth-server:7002, the connection is bridged to a client SOCKS5 worker. You must speak SOCKS5 protocol (handshake → connect → data) to reach the target service.
