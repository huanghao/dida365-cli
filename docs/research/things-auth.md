# things-cli 认证方式调研（及 Dida365 适配）

## 结论先行
- `things-cli` 不是 OAuth 流程。
- 它依赖 Things URL Scheme 的本地授权 token：`THINGS_AUTH_TOKEN`。
- token 来源：Things App 设置页手动复制或启用 CLI 访问。
- token 使用：写操作命令在构造 URL 时附带 `auth-token` 参数。

## things-cli 的认证机制

### 1) 凭证来源
- 用户在 Things 3 中开启 URL 访问并获取 token。
- CLI 通过环境变量读取：`THINGS_AUTH_TOKEN`。
- 部分命令支持 `--auth-token` 显式传入；优先级为：
  - `--auth-token`
  - `THINGS_AUTH_TOKEN`

关键实现：`internal/cli/token.go`
- `resolveAuthToken()`：统一解析 token，缺失时返回 `ErrMissingAuthToken`。

### 2) 认证命令（诊断而非登录）
- `things auth` 不发起授权，仅检查 token 是否存在并输出配置指引。

关键实现：`internal/cli/auth.go`

### 3) 哪些命令要求 token
- 典型需要 token：`update`、`update-project`、`undo`（恢复 update 操作时）。
- 典型不需要 token：只读数据库查询类命令（`tasks`/`today`/`search` 等）。

### 4) 认证注入点
- URL scheme 构造层强制校验 token。

关键实现：`internal/things/update.go`
- `BuildUpdateURL()`：若 `AuthToken == ""` 直接报错。

### 5) 错误行为
- token 缺失时立刻失败；不会尝试交互式登录。
- 在 `--debug` 下可打印 token 来源，便于排错。

## 设计评价
- 优点：实现简单、可脚本化、无浏览器回调复杂度。
- 缺点：
  - 不适用于跨设备/云端 API 场景。
  - 凭证管理依赖环境变量，缺少 refresh token 生命周期治理。

## 对 dida365-cli 的适配建议（基于官方 OpenAPI）

## 背景差异
- Dida365 OpenAPI（`https://developer.dida365.com/docs/openapi.md`）采用 OAuth2 授权码模式。
- 需要：`client_id/client_secret` + 浏览器授权 + `code` 换取 access token。

### 推荐方案（借鉴 things-cli 交互风格，不照搬机制）
1. 保留 `dida auth` 命令，但职责升级为“完整鉴权入口”：
- `dida auth login`：发起 OAuth 授权流程
- `dida auth status`：查看 token 状态
- `dida auth logout`：清理本地凭证

2. 凭证优先级建议：
- 显式参数 > 环境变量 > 本地配置文件

3. 存储建议：
- `~/.config/dida365-cli/config.toml` 保存客户端配置与租户信息。
- `~/.config/dida365-cli/credentials.json` 保存 access/refresh token 与过期时间。

4. 请求层统一注入：
- 在 API Client middleware 统一附加 `Authorization: Bearer <token>`。
- 401 时自动尝试 refresh；失败则提示重新登录。

5. 错误体验对齐 things-cli：
- 错误信息尽量可执行（告诉用户下一步命令）。
- 提供 `--debug` 显示 token 来源与 refresh 行为。

## MVP 落地建议
- v0.1：先支持授权码换 token + 手动重登（无 refresh 自动续期也可）。
- v0.2：补 refresh token 与过期前自动续期。
- v0.3：支持多账号配置与 profile 切换。

## 参考
- `~/workspace/sources/things3-cli/internal/cli/auth.go`
- `~/workspace/sources/things3-cli/internal/cli/token.go`
- `~/workspace/sources/things3-cli/internal/things/update.go`
- `https://developer.dida365.com/docs/openapi.md`（访问时间：2026-02-22）
