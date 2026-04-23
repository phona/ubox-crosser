# REQ-m15-1776940432: /buildinfo endpoint

## One-liner

在 `cmd/server` 暴露 HTTP 端点 `/buildinfo`,返回 JSON `{"git_sha": "<7-char short SHA>", "build_id": "<env BUILD_ID>", "go_version": "go1.23"}`,用作 sisyphus M15 verifier rename(`dev`→`dev_cross_check`、`spec`→`spec_lint`)的端到端验证抓手。

## Why

这是一个 **TEST 类需求**,目的不是用户功能,而是给 sisyphus orchestrator 一个最小、闭合、可机械验证的 REQ,跑通 analyze → spec(contract + acceptance)→ dev → staging-test → pr-ci → accept 全链路,确认 M15 的 stage 改名没把流水线打断。

业务侧顺带得到一个轻量构建溯源端点(回答"线上跑的是哪个 commit、哪个构建批次"),与已有 /version 互补但更全。

## Scope (in)

- 新增 HTTP 路由 `/buildinfo`(server 子命令,与 /version 共享同一监听端口,不开新 listener)
- 编译期通过 `-X main.GitSHA=$(git rev-parse --short HEAD)` 注入 git_sha 缩短为 7 字符
- 运行期读 `BUILD_ID` 环境变量,缺省 "dev"
- `go_version` 字面量 "go1.23"
- 端点未鉴权,返回 200 + JSON
- 单测覆盖 handler(纯函数测,httptest)
- 集成测走真实 docker stack(`tests/integration/`),GET 实际拿 200 + 校验三字段

## Scope (out)

- 不动 /version 既有行为
- 不引入 net/http mux 抽象 —— 跟 /version 一样,在现有 `handleHTTPRequest` 里多加一条路由分支即可
- 不为 client / auth_server 加同名端点(只 server)
- 不做版本协商、不返回 commit message、不返回 build time

## 影响范围

| 文件 | 改动 |
|---|---|
| `cmd/server/server.go` | 新增 `var BuildID = "dev"` 包变量(ldflags 不注入,运行期 env 覆盖)、把 BuildID 传给 ProxyServer |
| `server/server.go` | ProxyServer 增 buildID 字段;`handleHTTPRequest` 增 `/buildinfo` 分支;新增 `isBuildinfoRequest` |
| `server/server_test.go`(新建)| handler 单测 |
| `Makefile` | `-X main.GitSHA=$(GIT_SHA_SHORT)` 改用 `git rev-parse --short HEAD`(注:与 REQ-final2 已注入的 full SHA 不同 —— 见 design.md "git_sha 长度冲突") |
| `tests/integration/buildinfo_test.go`(新建)| 集成测 |
| `tests/docker-compose.yml` 或 `tests/integration/docker-compose.yml` | 注入 BUILD_ID env 给 server 服务 |

## 支持性判定

- **可拆分?** 否。改动量 <100 LOC,共享 ProxyServer 结构体改造,强耦合。单 dev agent。
- **跨 repo?** 否。phona/ubox-crosser 单仓库。
- **依赖在飞 REQ?** 是 —— REQ-e2e-1776916220 + REQ-final2-1776868985(共享 stage/REQ-e2e-1776916220-dev 分支)已立 HTTP 嗅探脚手架 + GitSHA ldflags + /version 路由。本 REQ 期望那批先 merge,我们再 rebase 上去做小增量。详见 design.md 的 "并发 REQ 协调" 段。
- **集成 repo?** 否,纯单仓库。
