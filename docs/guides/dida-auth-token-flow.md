# Dida365 认证与拿 Token 实操指南

本文针对你现在的状态：
- 已在开发者中心创建 App
- 已拿到 `client_id` 和 `client_secret`

目标：一步步拿到 `access_token`，并验证可调用 OpenAPI。

## 0. 先准备

你还需要确认 1 个关键配置：
- `redirect_uri`：必须和开发者中心里配置的回调地址完全一致（协议、域名、端口、路径都要一致）。

Q：我没有一个在外网的应用，我只是一个本地的cli，我应该怎么设计redirect_uri？
A：本地 CLI 场景建议用这两种方案之一：
- 方案 1（推荐）：本地回环地址回调，例如 `http://127.0.0.1:53682/callback`  
  优点：后续可自动化（CLI 临时起一个本地 HTTP 监听，自动拿 code）。
- 方案 2：`urn:ietf:wg:oauth:2.0:oob`（仅当平台明确支持时）  
  特点：不需要本地服务，但通常要手动复制 code。

你当前先走手动流程，建议直接在开发者中心把回调地址配置成：

```text
http://127.0.0.1:53682/callback
```

即使现在先手动复制 code，这个地址也为后续自动化登录保留了升级路径。

建议先在终端准备变量：

```bash
export DIDA_CLIENT_ID='你的_client_id'
export DIDA_CLIENT_SECRET='你的_client_secret'
export DIDA_REDIRECT_URI='你在开发者中心配置的_redirect_uri'
export DIDA_SCOPE='tasks:read tasks:write'
```

再生成一个随机 `state`（防止 CSRF）：

```bash
export DIDA_STATE="dida-$(date +%s)"
```

## 1. 获取授权码（code）

在浏览器打开下面 URL（把变量替换成真实值，注意 URL 编码）：

```text
https://dida365.com/oauth/authorize?client_id=<client_id>&scope=tasks%3Aread%20tasks%3Awrite&state=<state>&redirect_uri=<urlencoded_redirect_uri>&response_type=code
```

你可以用 Python 帮你安全编码并打印完整 URL：

```bash
python - <<'PY'
import os, urllib.parse
params = {
  'client_id': os.environ['DIDA_CLIENT_ID'],
  'scope': os.environ.get('DIDA_SCOPE', 'tasks:read tasks:write'),
  'state': os.environ['DIDA_STATE'],
  'redirect_uri': os.environ['DIDA_REDIRECT_URI'],
  'response_type': 'code',
}
print('https://dida365.com/oauth/authorize?' + urllib.parse.urlencode(params))
PY
```

授权成功后，浏览器会跳转到你的 `redirect_uri`，URL 上会带：
- `code=...`
- `state=...`

把 `code` 复制出来：

```bash
export DIDA_AUTH_CODE='从回调URL里复制的code'
```

## 2. 用 code 换 access_token

按 Dida 文档，token 端点是：
- `POST https://dida365.com/oauth/token`
- `client_id:client_secret` 通过 Basic Auth 传
- `Content-Type: application/x-www-form-urlencoded`

执行：

```bash
curl -sS -X POST 'https://dida365.com/oauth/token' \
  -u "$DIDA_CLIENT_ID:$DIDA_CLIENT_SECRET" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'grant_type=authorization_code' \
  --data-urlencode "code=$DIDA_AUTH_CODE" \
  --data-urlencode "scope=$DIDA_SCOPE" \
  --data-urlencode "redirect_uri=$DIDA_REDIRECT_URI"
```

预期会返回 JSON，至少包含：
- `access_token`
- 通常还会有 `refresh_token`、`expires_in` 等字段

你可以把 `access_token` 提取到环境变量：

```bash
export DIDA_ACCESS_TOKEN='上一步响应中的access_token'
export DIDA_REFRESH_TOKEN='上一步响应中的refresh_token'
```

Q：获得了access token已经，cli会把它保存在哪个位置了？
这个access token有过期时间吗？如果有要怎么定期刷新呢？刷新和重新获得的流程一样吗？
A：
- 当前 `dida` 实现会保存到本地配置文件：
  - 默认：`~/.config/dida365-cli/config.json`
  - 该默认路径是项目固定规则（基于 `$HOME` 拼接），不是按操作系统切换目录。
  - 可用 `--config` 指定其他路径
- 是有过期时间的。token 响应里通常会有 `expires_in`，CLI 会换算并保存 `expires_at`。
- 刷新方式：
  - 优先用 `refresh_token` 刷新（`grant_type=refresh_token`）
  - 当前命令：`dida auth refresh`
- 刷新流程和“重新授权拿 code”不是同一个流程：
  - 刷新：不需要用户再次登录授权（只要 refresh token 还有效）
  - 重新授权：需要再次走浏览器授权并拿新的 code（通常用于 refresh token 失效时）

## 3. 验证 token 可用

用 token 调一个最简单的接口：项目列表

```bash
curl -sS 'https://api.dida365.com/open/v1/project' \
  -H "Authorization: Bearer $DIDA_ACCESS_TOKEN"
```

能返回项目 JSON，即说明认证链路已打通。

## 4. access_token 过期后刷新

```bash
curl -sS -X POST 'https://dida365.com/oauth/token' \
  -u "$DIDA_CLIENT_ID:$DIDA_CLIENT_SECRET" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'grant_type=refresh_token' \
  --data-urlencode "refresh_token=$DIDA_REFRESH_TOKEN" \
  --data-urlencode "scope=$DIDA_SCOPE"
```

把新返回的 `access_token`（和可能更新的 `refresh_token`）覆盖保存。

## 5. 用本项目 dida CLI 走同一流程

如果你希望不手工拼 curl，也可以：

1. 初始化 OAuth 配置
```bash
go run ./cmd/dida auth init \
  --client-id "$DIDA_CLIENT_ID" \
  --client-secret "$DIDA_CLIENT_SECRET" \
  --redirect-uri "$DIDA_REDIRECT_URI"
```

2. 生成授权 URL
```bash
go run ./cmd/dida auth login
```

Q：重定向后我需要手动复制里面的code对吗？这个过程可能自动化吗？
现在不用自动化，但是在后续设计login流程的时候可能需要，不然用户使用这个命令的时候会很困惑。比如运行dida login，然后弹出一个用户登录，然后到一个落地页，然后就可以了。用户不需要直接细节，也不需要手动复制什么。
A：对，当前版本是手动复制 code。这个过程完全可以自动化，而且建议后续做：
- 目标体验：
  - `dida auth login`
  - CLI 自动打开浏览器
  - CLI 在本地 `127.0.0.1:<port>` 监听回调
  - 收到 code 后自动换 token 并落盘
  - 终端显示“登录成功”，用户无需手工复制
- 实现要点：
  - 回调地址固定为本地回环地址（需在开发者中心预先登记）
  - 本地监听加超时和 state 校验
  - 授权成功后返回一个简短落地页（例如“可以关闭本页面”）

这是后续很值得做的优先体验项，会明显降低使用门槛。

3. 拿到 code 后换 token
```bash
go run ./cmd/dida auth token --code "$DIDA_AUTH_CODE"
```

4. 查看状态
```bash
go run ./cmd/dida auth status
```

5. 刷新 token
```bash
go run ./cmd/dida auth refresh
```

Q：dida命令考虑自动定期刷新吗？还是不考虑？为什么？
A：建议“考虑自动刷新”，但分阶段做：
- 阶段 1（已具备基础能力）：手动刷新  
  - `dida auth refresh`，逻辑简单、可控、便于排错。
- 阶段 2（建议尽快）：请求前检查 `expires_at`，临近过期自动刷新。
- 阶段 3（更稳妥）：遇到 401 时自动刷新一次并重试一次请求。

为什么这么做：
- 全自动刷新能显著降低命令失败率与用户心智负担。
- 分阶段推进可以控制复杂度，避免一上来把认证链路做得过重。

## 6. 最容易卡住的点

- `redirect_uri` 不完全一致：最常见，必须与开发者中心配置严格一致。
- code 已被用过：授权码通常只能用一次。
- scope 不匹配：申请和换 token 时 scope 要一致。
- Basic Auth 错误：`-u client_id:client_secret` 任何一个错误都会失败。
- 时钟偏差/过期：token 过期后要走 refresh。

## 7. 安全建议

- 不要把 `client_secret`、`access_token`、`refresh_token` 提交到 Git。
- 建议只放到本地环境变量或本地配置文件（且权限最小化）。
