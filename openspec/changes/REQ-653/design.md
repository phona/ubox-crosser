---
id: REQ-653
title: GET /version endpoint – Design
---

## Decision: Separate HTTP admin server

The existing ubox-crosser uses raw TCP with a custom binary protocol — no HTTP framework. Rather than mixing HTTP into the TCP message handler, we add a lightweight `net/http` server on a separate port (`--admin-addr`, default `:8080`).

### Trade-offs

| Option | Pros | Cons |
|--------|------|------|
| Embed in TCP protocol | No new port | Breaks protocol, clients must understand HTTP |
| Separate HTTP server | Clean separation, standard tooling | Extra port to expose |

**Chosen:** Separate HTTP server. Clean, zero coupling to the tunnel protocol.

## Build-time injection

Use `go build -ldflags -X` to set `version.Commit` and `version.BuildTime` at compile time. The `version.Version` constant is hardcoded as `0.1.0` — it changes with code, not builds.
