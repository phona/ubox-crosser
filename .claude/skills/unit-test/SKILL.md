---
name: unit-test
description: Write unit tests for Go packages. Trigger when user asks to create, add, or write unit tests, or mentions testing a function, method, or package. These are package-level tests — NOT integration tests.
allowed-tools: Read, Glob, Grep, Write, Edit, TodoWrite, AskUserQuestion, Bash
---

# Unit Test Writing Skill

## Overview

This skill guides writing unit tests for Go packages following project conventions.
Tests live alongside the code they test and are organized by test category (pure vs service-dependent).

Full rules and anti-patterns: see [rules.md](rules.md)
Templates: [templates/pure.go.tmpl](templates/pure.go.tmpl), [templates/service.go.tmpl](templates/service.go.tmpl)

---

## Step 1 — Identify Target (MANDATORY)

Before writing any test, identify:

1. **Read the source file** being tested — understand function signatures, dependencies, and return types
2. **Classify the test category**:
   - **Pure**: No external services needed (logic, parsing, validation, data transformation)
   - **Service-dependent**: Requires network connections, external services

From the source code, derive:
- **Package**: The Go package being tested
- **Target functions**: Which functions/methods to test
- **Dependencies**: What external services are needed (if any)
- **Test style**: Internal (`package foo`) or external (`package foo_test`)

---

## Step 2 — Choose Test Strategy

Ask the user:

```yaml
AskUserQuestion:
  Q1: Which functions/methods should be tested?
    (free text or "all exported")
  Q2: Test style?
    options:
      - Internal (same package, access private fields)
      - External / blackbox (package_test, only public API)
  Q3: Does this need external services?
    options:
      - No — pure logic test
      - Yes — TCP/network connections
```

**Smart defaults** (apply if user doesn't specify):
- Pure utility functions → External blackbox, no services
- Coordinator/Dispatcher → Internal, may need network
- Config parsing → External blackbox, no services

---

## Step 3 — Generate Test File

Check if a test file already exists for the package:
- **Exists**: append new test functions, respect existing structure
- **New file**: create from the appropriate template

### Convention Checklist (enforce before writing)

1. ✅ **No hardcoded addresses** — all service connections via environment variables
2. ✅ **t.Skip when service unavailable** — never fail because external service is down
3. ✅ **Table-driven tests** for functions with multiple input/output combinations
4. ✅ **testify/assert** for assertions (not raw `if` + `t.Errorf` for simple checks)
5. ✅ **t.Fatal / t.Error** — never use `panic` in tests
6. ✅ **t.Parallel()** for independent tests when safe
7. ✅ **Subtests** (`t.Run`) for logical grouping
8. ✅ **TestMain** when package needs shared setup (logger, connections)
9. ✅ Build tag `//go:build integration` for service-dependent tests
10. ✅ No build tag for pure tests (they run everywhere)

### Anti-patterns to REJECT

- ❌ Hardcoded IP addresses
- ❌ Hardcoded passwords or credentials
- ❌ `panic()` instead of `t.Fatal()` or `t.Error()`
- ❌ Package-level mutable `var` shared across tests
- ❌ Tests that fail when a service is simply not running
- ❌ `time.Sleep` for synchronization (use channels, WaitGroups, or `assert.Eventually`)

---

## Step 4 — Run and Verify

After writing, run the tests:

**Pure tests (no build tag):**
```bash
go test -v -run "TestFunctionName" ./path/to/package/...
```

**Service-dependent tests (with build tag):**
```bash
go test -v -tags=integration -run "TestFunctionName" ./path/to/package/...
```

Verify:
- All tests pass (or skip gracefully if service unavailable)
- No hardcoded addresses in the diff
- `go vet` passes on the test file
