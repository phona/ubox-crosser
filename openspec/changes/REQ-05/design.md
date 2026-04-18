# REQ-05: Design Decisions

## D1: Which binary gets the health endpoint?

**Decision:** Proxy server (`cmd/server`) only.

**Rationale:** The proxy server is the central hub that both clients and auth servers connect to. It is the primary deployment target for orchestrators. Adding health to client/auth_server can be done later as separate REQs if needed.

## D2: Separate HTTP listener vs. multiplexing on the TCP port

**Decision:** Separate HTTP listener on a dedicated port.

**Rationale:** The existing TCP port speaks a custom binary/JSON protocol. Multiplexing HTTP on the same port would require protocol detection (e.g., checking for "GET " prefix) which adds complexity and risk to the critical path. A separate port is clean, standard, and trivially configurable.

## D3: Default health port

**Decision:** `:8080` as the default `--health-address`.

**Rationale:** Port 8080 is a common convention for HTTP health/admin endpoints in infrastructure tooling. It avoids conflict with the main TCP port (default `:7000`). The address is fully configurable via CLI flag (`--health-address`) and config file field (`health_address`).

## D4: No new dependencies

**Decision:** Use Go stdlib `net/http` only.

**Rationale:** A single endpoint does not justify adding a router framework. `http.NewServeMux` from the standard library is sufficient and keeps the dependency footprint unchanged.

## D5: Failure behavior on bind error

**Decision:** Log error, do not crash. The main TCP server continues operating.

**Rationale:** The health endpoint is an observability aid, not core functionality. If port 8080 is unavailable, the proxy should still serve traffic. Operators will notice the missing health endpoint through their monitoring and can reconfigure.

## D6: Response format

**Decision:** JSON body `{"status":"ok"}` with `Content-Type: application/json; charset=utf-8`.

**Rationale:** JSON is the most interoperable format for health checks across tooling (Docker HEALTHCHECK with curl+jq, Kubernetes httpGet probes, load balancer health checks). The `status` field provides a simple, extensible key if future versions need to report degraded states.
