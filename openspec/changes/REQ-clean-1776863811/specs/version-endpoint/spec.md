# Acceptance Spec: /version 端点

## Overview
验收规范定义了 `/version` HTTP 端点在部署后的端到端行为验证。这些场景从用户和运维视角验证版本信息的可访问性、准确性和可靠性。

---

## ADDED

### FEATURE-A1: 部署后端点立即可访问
**用户故事**: 作为运维人员，我希望在服务启动后立即能够访问 `/version` 端点来验证部署状态。

**Given** 服务已成功启动并就绪  
**When** 向 `GET /version` 发送 HTTP 请求  
**Then** 返回 HTTP 200 状态码  
**And** 响应体是有效的 JSON 格式  
**And** 响应包含 `sha` 字段（非空字符串）  

**验证前置条件**:
- 服务已完全启动（readiness probe 通过）
- 网络连接正常

**预期响应结构**:
```json
{
  "sha": "<git-commit-sha>",
  "build_time": "2026-04-22T12:34:56Z"
}
```

---

### FEATURE-A2: 返回的 git SHA 与部署代码一致
**用户故事**: 作为发布工程师，我希望验证返回的 git SHA 确实对应当前部署的代码版本，确保没有版本不匹配。

**Given** 服务已部署，对应 git commit SHA 为 `<expected-sha>`  
**When** 向 `/version` 发送 GET 请求  
**Then** 响应中的 `sha` 字段等于 `<expected-sha>`  
**And** 可以通过 `git show <sha>` 验证该提交确实对应当前部署代码  

**验证方法**:
```bash
# 获取当前部署的 SHA
DEPLOYED_SHA=$(curl -s http://service/version | jq -r '.sha')

# 验证 SHA 有效性（是否在 git 历史中）
git cat-file -t $DEPLOYED_SHA  # 应返回 "commit"
```

**风险缓解**:
- 确保编译时正确注入了 `-ldflags` 中的 SHA
- CI/CD 流程获取 SHA 时使用 `git rev-parse HEAD`

---

### FEATURE-A3: 端点支持幂等调用（多次调用结果一致）
**用户故事**: 作为监控系统管理员，我希望多次轮询 `/version` 端点进行健康检查时，同一服务实例返回的版本信息始终一致，避免版本漂移告警。

**Given** 服务已启动，版本信息已注入  
**When** 连续调用 `/version` 端点 10 次  
**Then** 所有响应都返回 HTTP 200  
**And** 所有响应中的 `sha` 字段完全相同  
**And** 所有响应中的 `build_time` 字段相同  

**验证步骤**:
```bash
# 多次调用并收集 SHA
for i in {1..10}; do
  curl -s http://service/version | jq -r '.sha'
done | sort -u | wc -l  # 应输出 1（所有 SHA 相同）
```

**性质**: 关键验收条件，确保版本信息的稳定性

---

### FEATURE-A4: 端点响应时间满足性能要求（<100ms）
**用户故事**: 作为站点可靠性工程师，我希望版本检查端点足够快速，不成为监控或负载均衡器的性能瓶颈。

**Given** 服务在正常运行  
**When** 向 `/version` 发送 GET 请求  
**Then** 响应时间 < 100ms（包括网络往返）  
**And** 端点在高并发下（100+ req/s）保持 <100ms 的响应时间  

**验证工具**:
```bash
# 单次调用测试
time curl -s http://service/version > /dev/null

# 并发压力测试（可选）
ab -n 1000 -c 100 http://service/version
# 检查 Requests per second、Time per request
```

**性能基线**:
- 平均响应时间: < 50ms
- P99 响应时间: < 100ms
- 吞吐量: > 1000 req/s

---

### FEATURE-A5: 错误处理和边界条件
**用户故事**: 作为系统集成测试工程师，我希望验证 `/version` 端点在各种调用方式下都能正确响应，避免意外的 HTTP 错误。

**Given** 服务已启动

**Scenario A5a: 不支持的 HTTP 方法被正确拒绝**  
**When** 向 `/version` 发送 POST 请求  
**Then** 返回 HTTP 405 (Method Not Allowed) 或 404  
**And** 响应包含适当的错误信息  

**Scenario A5b: GET 请求无查询参数**  
**When** 向 `/version?foo=bar` 发送 GET 请求  
**Then** 返回 HTTP 200（忽略查询参数）  
**And** 响应内容与无参数请求相同  

**Scenario A5c: 无效的请求头处理**  
**When** 向 `/version` 发送包含非标准 Accept-Language 的 GET 请求  
**Then** 返回 HTTP 200 和有效的 JSON  

---

## 验收标准总结

| 编号 | 功能 | 验收条件 | 优先级 |
|------|------|--------|------|
| A1 | 端点可访问性 | 返回 HTTP 200，有效 JSON，包含 sha 字段 | 🔴 关键 |
| A2 | SHA 准确性 | 返回的 SHA 与部署代码一致 | 🔴 关键 |
| A3 | 幂等性 | 多次调用返回相同结果 | 🟡 重要 |
| A4 | 性能 | 响应时间 < 100ms | 🟡 重要 |
| A5 | 错误处理 | 正确处理不支持的方法和边界条件 | 🟢 低 |

---

## 执行环境要求

- **部署环境**: Docker 或 K3s 环境，服务已运行
- **测试工具**: `curl`、`jq`、Git 命令行
- **访问权限**: 能访问服务的 HTTP 端口
- **时间同步**: 服务器时间同步（验证 build_time 时需要）

---

## 自动化执行集成

这些场景可通过以下方式自动化验证:

1. **Docker Compose 黑盒测试**: 启动服务容器，通过 HTTP 客户端验证各场景
2. **Kubernetes Pod 测试**: 在 K3s 环境中部署服务，通过 exec 或端口转发验证
3. **CI/CD 集成**: 在 GitHub Actions 中运行验收测试，作为部署前置条件

示例测试脚本将在 `tests/acceptance/` 目录下提供。
