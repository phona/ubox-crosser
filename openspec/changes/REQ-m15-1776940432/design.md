# Design: /buildinfo endpoint

## 现状(master HEAD = 233e8fa)

- `cmd/server` 是 cobra 入口,启动 `server.NewProxyServer` 跑 TCP 代理(shadowsocks)
- **没有 HTTP 监听** —— 所有 HTTP 端点(`/version`、`/healthz`、待加的 `/buildinfo`)都依赖在飞分支 `stage/REQ-e2e-1776916220-dev` 引入的"TCP listener 上嗅探 HTTP 请求"机制

### 在飞分支 stage/REQ-e2e-1776916220-dev 已有

1. `cmd/server/server.go`:`var GitSHA = "unknown"`(ldflags 注入,full 40 字符)
2. `server/server.go`:
   - `ProxyServer` 多了 `gitSHA string` 字段
   - `NewProxyServer(configs, gitSHA)` 多一个参数
   - `handleConnection` 改成先 peek 4 字节,根据是否 HTTP 动词分流到 `handleHTTPRequest`
   - `handleHTTPRequest` 用裸 string 拼 HTTP 响应(不走 `net/http`),识别 `GET /version ` 走 `/version` 分支,否则 404
   - `bufferedConn` 包裹 net.Conn,把已 peek 的 4 字节回灌给 TCP 协议解析路径
3. `Makefile`:`go build -ldflags="-s -w -X main.GitSHA=$(GIT_SHA)" ...`,GIT_SHA 取 `git rev-parse HEAD`(full)
4. `tests/integration/version_test.go`:build tag `integration`,起容器后 `http.Get(http://server:port/version)`

## 本 REQ 增量

### 代码改动(假设 stage/REQ-e2e-1776916220-dev 已 merge 到 master)

1. **`cmd/server/server.go`** —— 加包变量 + 透传:
   ```go
   var GitSHA = "unknown"   // 已存在
   var BuildID = "dev"      // 新加,但 ldflags 不注入,看 BUILD_ID env

   func main() {
       // ...
       if v := os.Getenv("BUILD_ID"); v != "" {
           BuildID = v
       }
       proxy := server.NewProxyServer(configs, GitSHA, BuildID)
       // ...
   }
   ```

2. **`server/server.go`** —— ProxyServer 加字段 + 路由:
   ```go
   type ProxyServer struct {
       // ...
       gitSHA  string
       buildID string  // 新加
   }

   func NewProxyServer(configs map[string]config.ServerConfig, gitSHA, buildID string) *ProxyServer {
       // ...
       buildID: buildID,
   }

   func (p *ProxyServer) handleHTTPRequest(conn net.Conn, peek []byte) {
       // ...
       if p.isVersionRequest(string(totalData)) {
           // 既有
       } else if p.isBuildinfoRequest(string(totalData)) {
           body := fmt.Sprintf(`{"git_sha":"%s","build_id":"%s","go_version":"go1.23"}`+"\n",
               shortSHA(p.gitSHA), p.buildID)
           response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s",
               len(body), body)
           conn.Write([]byte(response))
       } else {
           // 既有 404
       }
   }

   func (p *ProxyServer) isBuildinfoRequest(req string) bool {
       return len(req) > 14 && req[:15] == "GET /buildinfo "
   }

   func shortSHA(s string) string {
       if len(s) >= 7 { return s[:7] }
       return s
   }
   ```

3. **`Makefile`** —— 不动既有 `-X main.GitSHA=...`(它注入 full SHA,我们用 `shortSHA` helper 截到 7 字符);**不**为 BUILD_ID 加 ldflags 注入,走 env。

   > **取舍**:文档要求"通过 ldflags 注入 7 字符 SHA"。但既有分支已注入 full SHA,改 Makefile 会触发跟 /version 的回归(REQ-final2-1776868985 spec FEATURE-A1 明确写"sha 字段是 40 字符 hex")。**采取折中**:Makefile 不动,代码层 `shortSHA` 截断只用于 `/buildinfo`;`/version` 继续返回 full SHA。

4. **单测 `server/server_test.go`** —— 直接构造 `ProxyServer{gitSHA: "abc1234deadbeef", buildID: "ci-42"}`,用 `httptest` 拿不动,因为我们不走 `net/http` 而是裸 conn。改方案:`net.Pipe()` 起一对 conn,goroutine 写 `GET /buildinfo HTTP/1.1\r\n\r\n`,主 goroutine 读响应字节、字符串切出 body、`json.Unmarshal` 校验三字段。

5. **集成测 `tests/integration/buildinfo_test.go`** —— 跟 `version_test.go` 同形,build tag `integration`,共用 `proxyAddr`。

6. **`tests/docker-compose.yml`** —— server 服务加 `environment: BUILD_ID=ci-test`。

## 并发 REQ 协调

| REQ | 状态 | 影响 |
|---|---|---|
| REQ-e2e-1776916220(/healthz)| stage 分支已有 dev 在跑 | 引入 HTTP 嗅探脚手架(我们依赖) |
| REQ-final2-1776868985(/version)| 同分支 stage/REQ-e2e-1776916220-dev | 引入 GitSHA ldflags + /version 路由(我们扩展) |
| REQ-m15-1776940432(本)| analyze | 在上述基础上加 /buildinfo |

**dev-agent 干活时的两种情况**:

- **情况 A** —— stage 分支已 merge 到 master:rebase feat/REQ-m15-1776940432 上 master,直接加上面的小增量。
- **情况 B** —— stage 分支还没 merge:dev-agent 自己把 stage 的脚手架手抄一份(HTTP 嗅探 + GitSHA 包变量),再在上面加 /buildinfo。merge 时由 PR-CI 解决冲突。

dev-agent 决策时机:开工时 `git fetch origin && git log master --grep "REQ-e2e-1776916220"` 检查是否已 merge,选 A 或 B。

## 风险

1. **裸字符串拼 HTTP 响应**:既有 /version 这么写,但脆弱(没 chunked、没 keep-alive、单个 4096 buf 限死)。本 REQ 跟随既有风格,不重构。
2. **HTTP 嗅探只看 4 字节前缀**:对 `GET /buildinfo` 没问题(开头是 `GET `),但若客户端发 GET 后带超长 query 把 4096 缓冲挤爆,/buildinfo 也会跟着 /version 一起出问题 —— 不在本 REQ 解决。
3. **net.Pipe 单测对 timeout 敏感**:写完响应必须 `conn.Close()`,否则读端阻塞。单测要 `time.AfterFunc` 兜底。
4. **集成测对 docker-compose 的 BUILD_ID env 敏感**:CI 跑时如果忘了注入,集成测会断言 `build_id == "dev"`(default),要 spec 明确兜底分支。

## 不在本 REQ 范围

- 不重构 HTTP 处理为 `net/http` mux
- 不为 /buildinfo 加 metrics、不打 prometheus
- 不暴露在 client / auth_server
- 不做 /buildinfo 的并发压测(scope 给 /healthz REQ-e2e 的 FEATURE-A6)
