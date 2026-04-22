# 设计方案：/version 端点 git SHA 注入

## 高层架构

```
构建流程                        运行时                    HTTP 请求
───────────                     ────────                  ────────
git describe                                    GET /version
  ↓                                              ↓
获取 SHA ──→ go build -ldflags ──→ 二进制 ──→  Handler ──→ JSON 响应
                                    (嵌入)
```

## 技术选型

### SHA 获取方式：编译时注入（-ldflags）

**why：**
- 不依赖运行时获取 git 信息（容器环境可能无 .git 目录）
- 构建时一次性注入，性能最优
- 与业界标准做法对齐

**how：**
```bash
GIT_SHA=$(git rev-parse HEAD)
go build -ldflags "-X main.GitSHA=${GIT_SHA}"
```

### 响应格式：JSON

```json
{
  "sha": "abc1234567..."
}
```

**why：**
- 标准 REST API 响应格式
- 易于扩展（将来可加 `version` / `build_time` 等字段）
- 与现有 ubox-crosser API 风格对齐

## 关键设计点

### 1. 全 SHA vs 简短 SHA
- **决策**：使用全 SHA（40 字符）
- **理由**：唯一性强，便于生产环境溯源，不需要解析 git 历史

### 2. 变量声明位置
- **决策**：`main.go` 中声明 package-level 变量 `var GitSHA string`
- **理由**：单一职责，便于 -ldflags 注入
- **示例**：
  ```go
  package main
  
  var GitSHA = "unknown"
  
  func versionHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{
      "sha": GitSHA,
    })
  }
  ```

### 3. 路由注册
- **路由模式**：`/version`（不需要前缀，直接在根路由注册）
- **HTTP 方法**：GET（幂等，无副作用）
- **状态码**：200（总是成功）

### 4. CI/CD 集成
- **Makefile**：增加 `build` target，注入 GIT_SHA
- **容器镜像构建**：Dockerfile 里执行 `make build`，确保 SHA 正确注入

## 风险与缓解

| 风险 | 缓解方案 |
|-----|---------|
| 分离构建导致 SHA 不匹配 | 在 makefile 中写死 GIT_SHA 获取，保证一致性 |
| 开发环境 SHA 为 unknown | 合法行为，仅在生产构建时保证注入 |
| 响应格式变更导致兼容性 | 通过 Contract spec 锁定接口，acceptance 覆盖 |

## 后续扩展点

- 可为 /version 加 query param 如 `?detail=true` 返回更多信息（版本号、构建时间等）
- 可与监控/日志系统集成，自动上报部署版本
