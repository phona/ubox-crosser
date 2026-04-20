# QA Report — REQ-528: README Version History Section

**Date:** 2026-04-20
**Agent:** accept-agent (issue #547 / yetji54a)
**Result:** FAIL — ephemeral environment not ready

---

## Scenario Matrix

| # | Scenario | Expected | Actual | Status |
|---|----------|----------|--------|--------|
| 1 | Spin up ttpos-arch-lab ephemeral env | Playwright + Android emulator + DB + services running | Environment not available | FAIL |
| 2 | Run tests/acceptance/* | All acceptance tests pass in e2e env | Not executed — no environment | BLOCKED |
| 3 | Run tests/ui/* | UI tests pass | No UI tests exist; not applicable | N/A |

---

## Environment Check Evidence

### Server: vm-node04 (5b25f0cd)
```
command: which ttpos-arch-lab; docker ps --filter name=ttpos-arch-lab
result: empty output — no ttpos-arch-lab tooling or containers found
```

### Server: vm-node01 (3c178ff1)
```
command: ls /opt/ttpos-arch-lab/
result: Makefile, README.md, build, charts, docs, scripts, seed-data, values, vm01-nginx-domains.yaml

command: cat /opt/ttpos-arch-lab/Makefile (head)
result: Helm/Kubernetes-based deployment (helm upgrade --install).
        Targets: deps, up, down, status, test-deploy, test-seed, etc.
        Requires: Kubernetes cluster, Helm, ghcr auth token
        This is infrastructure deployment tooling, NOT an ephemeral
        acceptance test runner with Playwright/Android emulator.
```

### Conclusion
`/opt/ttpos-arch-lab` contains Helm chart deployment scripts for a Kubernetes-based lab environment. It does **not** provide:
- Playwright browser automation
- Android emulator
- Ephemeral per-PR test isolation

The acceptance test infrastructure required by the ACCEPT workflow is not yet implemented.

---

## Code Review (informational, not a substitute for e2e)

The branch adds:
- `README.md`: Version History section with markdown table (1 entry: v0.0.1)
- `tests/acceptance/readme_version_history_test.go`: 7 test functions validating section existence, table format, ordering, etc. (build tag: `acceptance`)
- OpenSpec artifacts: proposal, design, dev spec, contract spec, tasks

The VERIFY agent (#544, result:pass) confirmed these tests pass locally with `go test -tags acceptance`. However, local verification is not equivalent to e2e acceptance in an ephemeral environment.

---

## Prior ACCEPT Attempts

| Issue # | Result | Reason |
|---------|--------|--------|
| #535 | fail | Environment not ready |
| #538 | fail | Environment not ready |
| #541 | fail | Environment not ready |
| #543 | fail | Environment not ready |
| #547 (this) | fail | Environment not ready |

---

## Recommendation

Once ttpos-arch-lab gains ephemeral environment support (Playwright + test runner), re-run ACCEPT for REQ-528.
