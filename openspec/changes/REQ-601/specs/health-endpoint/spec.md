# Capability: health-endpoint

## ADDED Requirements

### Requirement: Health endpoint returns OK on GET

代理服务进程 SHALL 在配置的 health 监听地址上暴露 `GET /health`，返回 HTTP 200，header `Content-Type: application/json`，body 严格等于 `{"status":"ok"}`。

#### Scenario: FEATURE-S1 Successful health check

- **GIVEN** 代理服务以 `health_address=":8080"` 启动
- **WHEN** 客户端发送 `GET /health` 到 health HTTP 监听器
- **THEN** 服务端返回 HTTP 200
- **AND** 响应头 `Content-Type` 等于 `application/json`
- **AND** 响应体严格等于 `{"status":"ok"}`

---

### Requirement: Health endpoint rejects non-GET methods

`/health` 路径上的非 GET HTTP 方法 SHALL 返回 HTTP 405 Method Not Allowed，并设置 `Allow: GET` 响应头。

#### Scenario: FEATURE-S2 POST /health returns 405

- **GIVEN** health 监听器已启动
- **WHEN** 客户端发送 `POST /health`
- **THEN** 服务端返回 HTTP 405
- **AND** 响应头 `Allow` 等于 `GET`

#### Scenario: FEATURE-S3 PUT /health returns 405

- **GIVEN** health 监听器已启动
- **WHEN** 客户端发送 `PUT /health`
- **THEN** 服务端返回 HTTP 405
- **AND** 响应头 `Allow` 等于 `GET`

#### Scenario: FEATURE-S4 DELETE /health returns 405

- **GIVEN** health 监听器已启动
- **WHEN** 客户端发送 `DELETE /health`
- **THEN** 服务端返回 HTTP 405
- **AND** 响应头 `Allow` 等于 `GET`

---

### Requirement: Unknown paths return 404

health HTTP 监听器 SHALL 对所有非 `/health` 路径返回 HTTP 404 Not Found。

#### Scenario: FEATURE-S5 Request to root path returns 404

- **GIVEN** health 监听器已启动
- **WHEN** 客户端发送 `GET /`
- **THEN** 服务端返回 HTTP 404

#### Scenario: FEATURE-S6 Request to unknown path returns 404

- **GIVEN** health 监听器已启动
- **WHEN** 客户端发送 `GET /metrics`
- **THEN** 服务端返回 HTTP 404

#### Scenario: FEATURE-S7 Trailing slash on /health/ returns 404

- **GIVEN** health 监听器已启动
- **WHEN** 客户端发送 `GET /health/`
- **THEN** 服务端返回 HTTP 404

---

### Requirement: Health address is configurable

代理服务 SHALL 支持通过 config 文件字段 `health_address`（JSON key `"health_address"`）和 CLI flag `--health-address` 指定 health HTTP 监听器的 `host:port`。

#### Scenario: FEATURE-S8 Health address set via config file

- **GIVEN** server config 文件的 common section 包含 `"health_address": ":9090"`
- **WHEN** 启动代理服务
- **THEN** health HTTP 监听器绑定到 `:9090`

#### Scenario: FEATURE-S9 Health address set via CLI flag

- **GIVEN** 未通过 config 文件指定 `health_address`
- **WHEN** 启动命令为 `ubox-crosser-server --health-address :9090 ...`
- **THEN** health HTTP 监听器绑定到 `:9090`

---

### Requirement: Health server disabled when address is empty

当 `health_address` 为空或未设置时，代理服务 SHALL NOT 启动 health HTTP 监听器，且 TCP 代理功能照常运行。

#### Scenario: FEATURE-S10 No health address configured

- **GIVEN** 启动参数与配置中均未提供 `health_address`
- **WHEN** 代理服务启动
- **THEN** 进程内不启动任何 HTTP 监听器
- **AND** TCP 代理监听器与原有行为一致

---

### Requirement: Health server startup errors are reported

health HTTP 监听器的启动失败（端口占用、地址非法等）SHALL 通过 `ProxyServer.errs` channel 上报，并由现有错误循环写入日志。

#### Scenario: FEATURE-S11 Port conflict on startup

- **GIVEN** 配置的 `health_address` 端口已被占用
- **WHEN** 代理服务启动
- **THEN** health server 通过 `ProxyServer.Err()` 返回 bind error
- **AND** 错误经由现有 logrus 错误循环输出

#### Scenario: FEATURE-S12 ErrServerClosed is not reported as failure

- **GIVEN** health server 正常运行后被显式关闭
- **WHEN** `ListenAndServe` 返回 `http.ErrServerClosed`
- **THEN** 不向 errs channel 写入错误
