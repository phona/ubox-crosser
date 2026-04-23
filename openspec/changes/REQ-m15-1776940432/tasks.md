# REQ-m15-1776940432 Tasks

## Stage: contract-tests (owner: contract-spec-agent)

- [ ] 在 `openspec/changes/REQ-m15-1776940432/specs/buildinfo-endpoint/spec.md` 写 acceptance scenarios(Gherkin 风格),覆盖:
  - FEATURE-A1: GET /buildinfo 返回 200 + JSON,字段 git_sha / build_id / go_version 都存在
  - FEATURE-A2: git_sha 是 7 字符(`^[a-f0-9]{7}$`),"unknown" 兜底时长度可不等于 7(单独场景)
  - FEATURE-A3: build_id 默认 "dev",BUILD_ID env 设了就读 env 值
  - FEATURE-A4: go_version 字面量是 "go1.23"(不是动态读 runtime)
  - FEATURE-A5: 端点未鉴权(无 Authorization header 也 200)
- [ ] 跑 `openspec validate REQ-m15-1776940432`(如有该工具)

## Stage: acceptance-tests (owner: acceptance-spec-agent)

- [ ] 新建 `tests/integration/buildinfo_test.go`,build tag `//go:build integration`:
  - [ ] `TestBuildinfoEndpointReturns200`:GET 拿 200
  - [ ] `TestBuildinfoEndpointJSONFields`:三字段都存在且非空
  - [ ] `TestBuildinfoEndpointGitSHAFormat`:`git_sha` 匹配 `^[a-f0-9]{7}$` 或 == "unknown"
  - [ ] `TestBuildinfoEndpointBuildID`:`build_id == os.Getenv("EXPECTED_BUILD_ID")`(从测试侧 env 读期望值,docker-compose 注入两边)
  - [ ] `TestBuildinfoEndpointGoVersion`:`go_version == "go1.23"`
- [ ] 修改 `tests/docker-compose.yml`(或 `tests/integration/docker-compose.yml`,看 dev-agent 在哪起 stack):server 服务加 `environment: - BUILD_ID=ci-buildinfo-test`
- [ ] 集成测共用 `tests/integration/main_test.go` 的 `proxyAddr`(若已存在)

## Stage: implementation (owner: dev-agent)

### Step 1:决定走情况 A 还是 B(见 design.md "并发 REQ 协调")

- [ ] `git fetch origin master && git log master --oneline | grep "REQ-e2e-1776916220"`
  - 找到 → 情况 A,rebase 即可
  - 没找到 → 情况 B,把 `stage/REQ-e2e-1776916220-dev` 的 HTTP 嗅探脚手架(`handleConnection` peek、`isHTTPRequest`、`handleHTTPRequest`、`bufferedConn`、`var GitSHA`、Makefile ldflags)cherry-pick / 手抄到 feat 分支

### Step 2:加 /buildinfo

- [ ] `cmd/server/server.go`:
  - 新加 `var BuildID = "dev"`
  - main 起手读 `os.Getenv("BUILD_ID")` 覆盖 BuildID
  - `server.NewProxyServer(configs, GitSHA, BuildID)`(签名扩成三参)
- [ ] `server/server.go`:
  - `ProxyServer` 加 `buildID string`
  - `NewProxyServer` 第三参 buildID,赋给字段
  - `handleHTTPRequest`:在 `isVersionRequest` 后加 `else if isBuildinfoRequest`,返回 `{"git_sha":"...","build_id":"...","go_version":"go1.23"}\n`,Content-Type application/json
  - 新加 `isBuildinfoRequest(req string) bool`:`len(req) > 14 && req[:15] == "GET /buildinfo "`
  - 新加 `shortSHA(s string) string`:返回前 7 字符,len<7 时返回原串(防御 "unknown" 等 fallback)

### Step 3:单测

- [ ] 新建 `server/server_test.go`(若已有,追加 `TestHandleBuildinfo`):
  - 用 `net.Pipe()` 起一对 conn
  - 构造 `ProxyServer{gitSHA: "abc1234deadbeef0000000000000000000000000", buildID: "ci-42"}`
  - goroutine 写 `GET /buildinfo HTTP/1.1\r\nHost: localhost\r\n\r\n` 到一端
  - 主 goroutine 调 `p.handleHTTPRequest(otherEnd, []byte("GET "))` —— 实际签名要看实现,可能是 `handleConnection` 入口
  - 读响应,字符串切出 body,`json.Unmarshal` 校验 `git_sha == "abc1234"`、`build_id == "ci-42"`、`go_version == "go1.23"`
  - `TestShortSHA` 边界:输入 ""、"abc"、"abc1234"、"abc1234extra" 各自返回 ""、"abc"、"abc1234"、"abc1234"
  - `TestIsBuildinfoRequest` 边界:`GET /buildinfo `(true)、`GET /buildinfo`(无尾空格,false)、`GET /buildinfox `(false)、`POST /buildinfo `(false)
- [ ] `make unit-test` 必须过

### Step 4:集成测准备

- [ ] 改 `tests/docker-compose.yml`(或 integration/docker-compose.yml):server 服务环境变量增 `BUILD_ID=ci-buildinfo-test`(同步改测试侧期望值,或测试侧也读相同 env 名)
- [ ] `make build`(确认 Makefile ldflags 注入 GitSHA 没回归)
- [ ] 本地跑 `docker compose -f tests/docker-compose.yml up --build`,手 curl `http://server:7000/buildinfo`(端口看 server 监听)确认 200 + JSON

### Step 5:质量门

- [ ] `make fmt vet lint`
- [ ] `make unit-test`
- [ ] `make test`(若集成测被 build tag gate,单独 `go test -tags=integration ./tests/integration/...`)

## 完成判定

- proposal/design/tasks/spec 文件齐
- unit test 通过,handler 三字段都校验
- 集成测在 docker stack 里 GET /buildinfo 拿 200,JSON 字段对齐 spec
- /version 既有行为不回归(回归 test 走 REQ-final2 既有用例)
