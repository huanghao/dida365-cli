# dida MVP 命令与 Dida365 OpenAPI 映射（草案）

## 目标
在 `things-cli` 交互风格基础上，定义 `dida` 首版可交付命令，并映射到 Dida365 OpenAPI。

官方文档来源：`https://developer.dida365.com/docs/openapi.md`（访问时间：2026-02-22）

## MVP 命令清单
- `dida auth login`
- `dida auth status`
- `dida projects list`
- `dida list`（按项目聚合查看任务）
- `dida show --project <id> --id <task_id>`
- `dida add --project <id> --title <text> [...options]`
- `dida update --id <task_id> --project <id> [...options]`
- `dida done --id <task_id> --project <id>`
- `dida delete --id <task_id> --project <id>`

## API 映射表

| dida 命令 | OpenAPI 端点 | 方法 | 说明 |
|---|---|---|---|
| `auth login` | `/oauth/authorize` + `/oauth/token` | Redirect + POST | OAuth2 授权码换 token |
| `projects list` | `/open/v1/project` | GET | 列出项目 |
| `list --project` | `/open/v1/project/{projectId}/data` | GET | 取项目全量数据（含 tasks/columns） |
| `show --project --id` | `/open/v1/project/{projectId}/task/{taskId}` | GET | 取单任务详情 |
| `add` | `/open/v1/task` | POST | 创建任务 |
| `update` | `/open/v1/task/{taskId}` | POST | 更新任务 |
| `done` | `/open/v1/project/{projectId}/task/{taskId}/complete` | POST | 标记完成 |
| `delete` | `/open/v1/project/{projectId}/task/{taskId}` | DELETE | 删除任务 |

## 参数映射建议

### 1) `dida add`
建议参数：
- `--project`（必填）
- `--title`（必填）
- `--content`
- `--desc`
- `--start`
- `--due`
- `--all-day`
- `--priority`
- `--repeat`

映射字段：`projectId`, `title`, `content`, `desc`, `startDate`, `dueDate`, `isAllDay`, `priority`, `repeatFlag`

### 2) `dida update`
建议参数：
- `--id`（必填）
- `--project`（必填，兼容 API body 要求）
- 其余字段与 `add` 对齐

### 3) `dida list`
建议参数：
- `--project`（首版必填，避免跨项目汇总复杂度）
- `--status`（可选，本地过滤）
- `--limit`
- `--format table|json`

说明：官方文档未给出“全局任务列表”端点，首版用 `project/{id}/data` 提供列表能力。

## 首版实现约束
- OAuth 首版可仅支持手动登录与 access token 保存，不强依赖 refresh 自动续期。
- `done/delete/show` 需要同时具备 `taskId + projectId`（由 API path 决定）。
- 与 `things-cli` 对齐：
  - 支持 `--debug`
  - 支持 `--dry-run`（打印请求信息，不实际发送）
  - 帮助文本内给出示例

## 已识别风险
- 文档中的更新接口描述为 `POST /open/v1/task/{taskId}`（不是 PUT），实现时需按文档执行。
- 任务查询依赖项目维度，跨项目汇总可能需要多次请求并在 CLI 聚合。
- OAuth token 生命周期字段在文档页未完整展开，需在实装阶段用真实响应校验。

## 下一步实现任务
1. 先实现 `auth login/status` 与 token 持久化。
2. 先打通 `projects list` + `list --project`。
3. 完成 `add/update/done/delete/show`。
4. 增加 `--format json` 与基础集成测试。
