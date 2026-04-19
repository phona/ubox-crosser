## 1. Version 包

- [ ] 1.1 创建 `internal/version/version.go`，定义 `Version`、`GitCommit`、`BuildTime` 变量和 `Info()` 返回结构体的方法
- [ ] 1.2 为 version 包编写单元测试，验证默认值和 JSON 序列化

## 2. HTTP 端点

- [ ] 2.1 在 `server/` 包中新增 HTTP handler，注册 `GET /version` 路由，返回 JSON 响应
- [ ] 2.2 在 `ProxyServer` 启动流程中根据配置启动 HTTP listener goroutine

## 3. 配置变更

- [ ] 3.1 在 `models/config/config.go` 的 `ServerConfig` 中新增 `HTTPAddress` 字段（json tag: `http_address`）
- [ ] 3.2 更新示例配置文件，添加 `http_address` 字段说明

## 4. 构建系统

- [ ] 4.1 修改 Makefile，在 `-ldflags` 中通过 `-X` 注入 version、git commit 和 build time
- [ ] 4.2 同步更新 Dockerfile 中的构建命令（如适用）

## 5. 测试

- [ ] 5.1 编写 HTTP handler 的单元测试，验证 /version 端点响应格式和状态码
- [ ] 5.2 验证 HTTP listener 在 `http_address` 为空时不启动
