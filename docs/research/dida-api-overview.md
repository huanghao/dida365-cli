# Dida365 OpenAPI 使用与认证说明（含核心业务概念）

本文用于快速回答三个问题：
- Dida API 怎么用？
- Dida API 怎么认证？
- Dida 的核心业务概念有哪些？

参考来源：`https://developer.dida365.com/docs/openapi.md`

## 1. API 怎么使用

## 1.1 基本调用方式
- 协议：HTTPS + JSON
- 基础域名：`https://api.dida365.com`
- OpenAPI 路径前缀：`/open/v1`
- 认证方式：`Authorization: Bearer <access_token>`

示例：
```http
GET /open/v1/project HTTP/1.1
Host: api.dida365.com
Authorization: Bearer <access_token>
```

## 1.2 常用核心接口（v1）
- 项目：
  - `GET /open/v1/project`：项目列表
  - `GET /open/v1/project/{projectId}`：项目详情
  - `GET /open/v1/project/{projectId}/data`：项目数据（含任务等）
  - `POST /open/v1/project`：创建项目
  - `POST /open/v1/project/{projectId}`：更新项目
  - `DELETE /open/v1/project/{projectId}`：删除项目
- 任务：
  - `POST /open/v1/task`：创建任务
  - `POST /open/v1/task/{taskId}`：更新任务
  - `GET /open/v1/project/{projectId}/task/{taskId}`：任务详情
  - `POST /open/v1/project/{projectId}/task/{taskId}/complete`：完成任务
  - `DELETE /open/v1/project/{projectId}/task/{taskId}`：删除任务

## 1.3 请求与响应要点
- 时间字段常见格式：`yyyy-MM-dd'T'HH:mm:ssZ`
- 常见任务字段：
  - `title`, `projectId`, `content`, `desc`
  - `startDate`, `dueDate`, `timeZone`
  - `isAllDay`, `priority`, `repeatFlag`
  - `items`（子任务/checklist）
- 常见错误码：
  - `401` 未认证（token 无效/缺失）
  - `403` 无权限（scope 不足或应用权限问题）
  - `404` 资源不存在

## 2. API 怎么认证（OAuth2）

Dida OpenAPI 使用 OAuth2 授权码流程。

## 2.1 前置条件
- 在开发者中心创建应用，拿到：
  - `client_id`
  - `client_secret`
  - `redirect_uri`

## 2.2 授权码流程
1. 引导用户打开授权地址：
   - `https://dida365.com/oauth/authorize`
   - 参数：
     - `client_id`
     - `redirect_uri`
     - `response_type=code`
     - `scope`（例如 `tasks:read tasks:write`）
     - `state`（防 CSRF）
2. 用户授权后，Dida 回调 `redirect_uri`，并带上 `code`
3. 服务端换 token：
   - `POST https://dida365.com/oauth/token`
   - `Content-Type: application/x-www-form-urlencoded`
   - 使用 Basic Auth 传 `client_id:client_secret`
   - 参数：
     - `grant_type=authorization_code`
     - `code`
     - `redirect_uri`
     - `scope`
4. 拿到 `access_token` 后，后续请求都带 `Bearer access_token`

## 2.3 刷新 token
- 使用 `grant_type=refresh_token` 到同一 token 端点换新 token
- 失败时常见处理：
  - refresh token 失效 -> 重新走授权码流程
  - scope 变化 -> 重新授权

## 3. 核心业务概念

## 3.1 项目（Project）
- 任务容器，可理解为任务分组/清单
- 通过 `projectId` 关联任务
- CLI 里常见是先列项目，再在项目下列任务

## 3.2 任务（Task）
- 主要业务实体
- 关键能力：创建、更新、完成、删除、查询详情
- 典型字段：标题、正文、开始时间、截止时间、重复规则、优先级

## 3.3 子任务 / 检查项（items）
- 任务内的子项列表
- 每个 item 有自己的状态、标题、排序、时间字段

## 3.4 状态与完成语义
- 任务一般包含状态字段（如未完成/完成等）
- “完成”操作由专门接口表达：`.../complete`
- 不建议仅本地改字段模拟完成，应调用完成接口保持服务端语义一致

## 3.5 时间与时区
- `startDate` / `dueDate` 和 `timeZone` 联动
- `isAllDay=true` 时，时间语义与普通定时任务不同
- CLI 层建议统一校验输入格式，避免调用时才报错

## 4. 对 dida CLI 的实现建议

- 命令层建议映射：
  - `projects list` -> `GET /project`
  - `list --project` -> `GET /project/{id}/data`
  - `show --project --id` -> `GET /project/{id}/task/{id}`
  - `add` -> `POST /task`
  - `update` -> `POST /task/{id}`
  - `done` -> `POST /project/{id}/task/{id}/complete`
  - `delete` -> `DELETE /project/{id}/task/{id}`
- 认证层建议：
  - 保留 `auth login/token/refresh/status/logout`
  - access token + refresh token 本地持久化
  - 后续可补自动 refresh（401 时尝试一次）

## 5. 常见排错清单

- `401 Unauthorized`：
  - 检查 `Authorization` 头是否为 `Bearer <token>`
  - 检查 token 是否过期
- `403 Forbidden`：
  - 检查应用 scope 是否包含所需权限（如 `tasks:write`）
- 授权后拿不到 `code`：
  - 检查 `redirect_uri` 是否和开发者后台配置完全一致
- 换 token 失败：
  - 检查 Basic Auth 的 `client_id/client_secret`
  - 检查 `grant_type` 和表单参数名
