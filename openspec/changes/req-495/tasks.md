## 1. Version 包

- [ ] 1.1 新建 `version/version.go`，定义 `Version`、`Commit`、`BuildTime` 包级变量及默认值
- [ ] 1.2 新建 `version/handler.go`，实现返回版本 JSON 的 `http.HandlerFunc`

## 2. Server 集成

- [ ] 2.1 在 `models/config/config.go` 的 `ServerConfig` 中新增 `HTTPAddress` 字段
- [ ] 2.2 在 `utils/conf/` 中确保 config file 解析支持 `http_address` 字段
- [ ] 2.3 在 `cmd/server/server.go` 中添加 `--http-address` CLI flag
- [ ] 2.4 在 `cmd/server/server.go` 中当 `HTTPAddress` 非空时启动 HTTP server goroutine

## 3. 构建系统

- [ ] 3.1 更新 `Makefile`，在 ldflags 中注入 `Version`、`Commit`、`BuildTime`
- [ ] 3.2 更新 `Dockerfile`，确保构建阶段传递版本 ldflags

## 4. 测试

- [ ] 4.1 为 `version` 包编写单元测试（handler 返回正确 JSON、Content-Type）
- [ ] 4.2 为 HTTP server 启动逻辑编写集成测试（可选）
