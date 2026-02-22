# things-cli 存储设计调研（SQLite）

## 结论先行
- `things-cli` 采用“混合写入策略”：
  - 主路径：通过 Things URL Scheme / AppleScript 写入应用数据。
  - 查询路径：直接读取本地 Things SQLite（`main.sqlite`）。
  - 特例写入：重复规则通过 SQLite 直接更新（`OpenWritable`）。
- 这是一种“能力优先”设计：当 URL Scheme 能力不足时，用 DB 补齐。

## 存储层结构

### 1) DB 驱动与连接
- 驱动：`modernc.org/sqlite`（纯 Go）。
- 连接封装：`internal/db/db.go` 的 `Store`。
- 支持两种打开方式：
  - `Open`（`mode=ro`）只读
  - `OpenWritable`（`mode=rw`）读写

### 2) 数据库路径解析
- `internal/db/path.go`:
  - 优先级：`--db` > `THINGSDB` > 默认 `ThingsData-*` 路径 > legacy 路径。
  - 自动挑选最新数据库文件。
- 失败返回 `ErrDatabaseNotFound`，CLI 层统一格式化为可执行错误提示。

### 3) 数据模型
- `internal/db/models.go` 定义：
  - `Task`、`Project`、`Area`、`Tag`
  - 查询过滤结构：`TaskFilter`、`ProjectFilter`
- 查询核心表（从 SQL 可见）：
  - `TMTask`（任务/项目/heading）
  - `TMArea`
  - `TMTag` / `TMTaskTag`
  - `TMChecklistItem`

### 4) 查询策略
- `internal/db/queries.go` 统一构造 SQL：
  - 支持状态、项目/区域/标签、全文检索、时间区间、排序、分页等过滤。
  - Today/Upcoming/Deadlines 等视图通过专用 where 条件组合而成。
- 输出层与查询层解耦：命令层先获取结构化结果，再按 table/json/jsonl/csv 输出。

### 5) 写入策略（关键）
- 常规增删改：
  - 使用 URL Scheme 或 AppleScript，不直接改 DB。
- 重复任务：
  - 使用 `OpenDefaultWritable` 直接更新 `TMTask` 的 recurrence 字段。
  - 典型字段：`rt1_recurrenceRule`、`rt1_nextInstanceStartDate` 等。

关键实现：`internal/db/repeat.go`
- `ApplyRepeatRule()`
- `ClearRepeatRule()`

### 6) 辅助存储
- `internal/cli/actionlog.go` 维护 `~/.config/things3-cli/actions.jsonl`。
- 作用：记录批量 update/trash 以支持 `undo`。
- 说明：这不是主数据存储，而是 CLI 恢复日志。

## 设计优缺点
- 优点：
  - 读取性能高，支持复杂查询和过滤。
  - 能绕过 URL Scheme 能力边界（如 repeat）。
- 风险：
  - 直接写 DB 有版本兼容风险。
  - 需要系统权限（Full Disk Access）。
  - 需要谨慎处理 WAL 与并发一致性。

## 对 dida365-cli 的建议

### 建议 1：默认不引入 SQLite 主存储
- Dida365 是云 API，首版应以远端 API 为主事实源。
- 本地 SQLite 可选，仅用于缓存和离线查询，不作为权威写入源。

### 建议 2：缓存策略（可选）
- 若需要高性能列表：
  - 本地缓存可用 SQLite/BoltDB。
  - 只缓存只读投影（task/project 摘要）并记录同步时间戳。

### 建议 3：写入路径统一走 API
- 避免 `things-cli` 式“URL 写 + DB 写”双路径复杂度。
- 所有写操作统一通过 Dida OpenAPI，错误可观测性更好。

### 建议 4：本地元数据与凭证分离
- 凭证：单独安全存储（凭证文件/系统 keychain）。
- 操作日志：可单独 JSONL（类似 `actionlog`）用于撤销/审计。

## 参考
- `~/workspace/sources/things3-cli/internal/db/db.go`
- `~/workspace/sources/things3-cli/internal/db/path.go`
- `~/workspace/sources/things3-cli/internal/db/models.go`
- `~/workspace/sources/things3-cli/internal/db/queries.go`
- `~/workspace/sources/things3-cli/internal/db/repeat.go`
- `~/workspace/sources/things3-cli/internal/cli/actionlog.go`
