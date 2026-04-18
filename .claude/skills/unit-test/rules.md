# Unit Test Rules & Best Practices

## Test Categories

### Category 1: Pure Tests (no build tag)

Tests that need **zero external services**. These run in CI, locally, and everywhere.

```go
package utils_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
    // pure logic — no network, no external service
}
```

**No build tag required.** These tests MUST pass with just `go test ./...`.

### Category 2: Service-Dependent Tests (build tag required)

Tests that require running proxy server, network connections, or other services.

```go
//go:build integration

package connector

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    os.Exit(m.Run())
}

func TestCoordinator_SendRecv(t *testing.T) {
    addr := os.Getenv("PROXY_SERVER_ADDR")
    if addr == "" {
        t.Skip("PROXY_SERVER_ADDR not set — skipping")
    }
    // ...
}
```

**Rules:**
1. MUST have `//go:build integration` on line 1
2. MUST read connection info from environment variables
3. MUST `t.Skip()` when the required service env var is empty

---

## Environment Variable Convention

### Proxy Server
```go
addr := os.Getenv("PROXY_SERVER_ADDR")   // empty = skip
```

### Auth Server
```go
addr := os.Getenv("AUTH_SERVER_ADDR")     // empty = skip
```

### Cipher
```go
key    := os.Getenv("TEST_CIPHER_KEY")
method := os.Getenv("TEST_CIPHER_METHOD")
```

### Helper function
```go
func envOrDefault(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
```

---

## Table-Driven Test Pattern

Use for any function with multiple input/output combinations:

```go
func TestMessageMarshal(t *testing.T) {
    tests := []struct {
        name     string
        msg      message.Message
        expected string
    }{
        {"login", message.Message{Type: message.LOGIN, ServeName: "svc"}, `{"type":0,"serve_name":"svc","password":""}`},
        {"heartbeat", message.Message{Type: message.HEART_BEAT}, `{"type":1,"serve_name":"","password":""}`},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            buf, err := json.Marshal(tt.msg)
            assert.NoError(t, err)
            assert.JSONEq(t, tt.expected, string(buf))
        })
    }
}
```

---

## Assertion Style

### Prefer testify/assert
```go
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.Nil(t, result)
require.NoError(t, err)  // fatal on failure
```

### Forbidden
```go
// ❌ NEVER use panic in tests
if err != nil {
    panic(err)  // WRONG — use t.Fatal(err)
}
```

---

## File Organization

```
models/config/
    config.go              # source code
    config_test.go         # pure tests (no build tag)

utils/connector/
    coordinator.go         # source code
    coordinator_test.go    # pure tests
    coordinator_integration_test.go  # network tests (//go:build integration)
```

**Rule:** Separate pure tests from service-dependent tests into different files.

---

## Common Anti-Patterns (REJECT these)

### 1. Hardcoded service addresses
```go
// ❌ NEVER
conn, _ := net.Dial("tcp", "192.168.1.100:7000")

// ✅ CORRECT
addr := os.Getenv("PROXY_SERVER_ADDR")
if addr == "" {
    t.Skip("PROXY_SERVER_ADDR not set")
}
conn, err := net.Dial("tcp", addr)
require.NoError(t, err)
```

### 2. Using panic instead of test assertions
```go
// ❌ NEVER
if err != nil { panic(err) }

// ✅ CORRECT
require.NoError(t, err)
```

### 3. Shared mutable state across tests
```go
// ❌ NEVER
var globalConn net.Conn

// ✅ CORRECT — each test creates its own
func TestSomething(t *testing.T) {
    conn, err := net.Dial(...)
    defer conn.Close()
}
```
