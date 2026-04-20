# QA Acceptance Report: REQ-567 GET /healthz

- **Date**: 2026-04-20
- **Agent**: accept-agent (BKD issue 8tb300ph)
- **Result**: FAIL
- **Reason**: Ephemeral environment not ready

## Environment Attempt

| Step | Command | Result |
|------|---------|--------|
| Check namespace | `kubectl get ns ttpos-arch-lab` | Active (101m) |
| Check pods | `kubectl get pods -n ttpos-arch-lab` | No resources found |
| Spin up env | `make up` (vm-node04:/opt/ttpos-arch-lab) | **FAIL** exit 2 |

### Error Output

```
Secret 'ghcr-pull' does not exist in namespace 'ttpos-arch-lab'
Please set GH_TOKEN env var or run in interactive terminal
make: *** [Makefile:29: auth] Error 1
```

## Scenario Matrix

| Scenario | Status | Notes |
|----------|--------|-------|
| GET /healthz returns 200 + JSON `{"status":"ok"}` | NOT RUN | env unavailable |
| Content-Type is application/json | NOT RUN | env unavailable |
| Response time < 200ms | NOT RUN | env unavailable |
| Endpoint accessible without auth | NOT RUN | env unavailable |

## Conclusion

The ttpos-arch-lab ephemeral environment on vm-node04 cannot be provisioned because:

1. The `ghcr-pull` Kubernetes secret is missing in the `ttpos-arch-lab` namespace
2. No `GH_TOKEN` environment variable is configured for non-interactive use
3. Without container image pull credentials, Helm cannot deploy the services

**No acceptance tests were executed.** The feature branch `feat/REQ-567` exists and VERIFY passed (result:pass), but e2e acceptance in the target environment is blocked on infrastructure setup.

### Recommendation

Before re-running acceptance:
1. Configure `GH_TOKEN` on vm-node04 (or create the `ghcr-pull` secret manually)
2. Ensure `make up` completes successfully in the `ttpos-arch-lab` namespace
3. Re-trigger the ACCEPT pipeline stage
