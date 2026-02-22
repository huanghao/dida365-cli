# things-cli 完整命令能力调研

## 能力总览
`things-cli` 是一个 macOS CLI，能力由三条技术路径组成：
- URL Scheme：新增/更新任务和项目。
- AppleScript：区域/项目/任务删除等自动化。
- SQLite 只读查询（+少量 repeat 直写）：高性能列表、搜索、过滤。

## 全局能力
- 根命令：`things`
- 全局参数：
  - `--debug`
  - `--foreground`
  - `--dry-run`
  - `-V, --version`

## 命令清单（按功能域）

### 1) 写操作
- `add`：新增 todo（支持大量字段与 repeat flags）
- `update`：更新 todo（支持按 `--id` 或按 query 批量）
- `delete`：删除/批量删除 todo（带确认机制）
- `add-area`（alias: `create-area`）
- `update-area`
- `delete-area`
- `add-project`（alias: `create-project`）
- `update-project`
- `delete-project`
- `undo`：回滚上次批量 update/trash（基于 actionlog）

### 2) 读操作（数据库）
- `tasks`（alias: `todos`）
- `search`
- `show`
- `projects`
- `areas`
- `tags`
- 快捷列表：
  - `inbox`, `today`, `upcoming`, `repeating`
  - `anytime`, `someday`, `logbook`
  - `logtoday`, `createdtoday`
  - `completed`, `canceled`, `trash`, `deadlines`
  - `all`

### 3) 认证/帮助
- `auth`：检查 token 状态并输出设置指引
- `help [command]`

## 参数体系（可迁移模式）

### 通用查询参数（大量列表命令复用）
- 状态与范围：`--status`, `--all`, `--include-trashed`
- 结构过滤：`--project`, `--area`, `--tag`
- 内容过滤：`--search`, `--query`
- 时间过滤：`--created-*`, `--modified-*`, `--due-before`, `--start-before`
- 分页排序：`--limit`, `--offset`, `--sort`
- 其它：`--has-url`, `--recursive`

### 输出参数（读命令复用）
- `--format table|json|jsonl|csv`
- `--select`
- `--json`
- `--no-header`

### 写命令通用行为
- 支持 STDIN（`-`）输入。
- 对高风险操作有确认门槛（`--confirm` / `--yes`）。
- `--dry-run` 可打印 URL/script，避免直接执行。

## 认证与权限边界
- `update` / `update-project` 等需 `THINGS_AUTH_TOKEN`。
- 数据库访问可能需要 Full Disk Access。
- AppleScript 操作可能触发系统自动化授权。

## 输出与可观测性
- 输出支持表格与机器可读格式（json/jsonl/csv）。
- 错误格式统一以 `Error:` 开头。
- `--debug` 可显示底层调用信息（open/osascript 命令等）。

## 对 dida CLI 的能力借鉴建议

### 建议保留
- 命令分组与清晰动词：`add/update/delete/list/show/search/auth`。
- 通用过滤参数和输出参数抽象。
- `--dry-run`、`--debug`、批量操作确认策略。

### 建议调整
- `things` 的 AppleScript/本地 DB 路径不适用于 Dida365。
- `dida` 应以远端 API 为主，实现一致的在线行为。

## 核心参考
- 命令注册：`~/workspace/sources/things3-cli/internal/cli/root.go`
- 通用查询参数：`~/workspace/sources/things3-cli/internal/cli/taskflags.go`
- 命令规格：`~/workspace/sources/things3-cli/doc/man/things.1.md`
- 用户概览：`~/workspace/sources/things3-cli/README.md`
