# Changelog

## 2026-02-22

### 0. 当前进度（持续更新）
- 更新时间：2026-02-22
- 阶段状态：
  - 阶段 A（调研与规格固化）：已完成
  - 阶段 B（工程骨架搭建）：已完成
  - 阶段 C（认证与 API 客户端）：进行中（核心链路已打通，测试待补强）
  - 阶段 D（MVP 命令实现）：进行中（核心命令已实现，参数与交互待增强）
  - 阶段 E（质量与发布准备）：进行中（README/Makefile/基础测试已补充）
- 本次已完成：
  - 完成 4 份 `things-cli` 调研文档与 1 份 `dida` MVP 映射文档
  - 初始化 Go 工程（`cmd/` + `internal/`）并接入 `cobra`
  - 实现 OAuth 相关命令：`auth init/login/token/status/logout`
  - 实现核心任务命令：`projects list`、`list`、`show`、`add`、`update`、`done`、`delete`
  - 补充 `README.md`、`Makefile`、`version` 命令与 `internal/dida` 基础单元测试

### 记录 #1 / commit `fbe1d2c`
- 本次进展结果：
  - 完成阶段 A 调研产出：
    - `docs/research/things-agent-docs.md`
    - `docs/research/things-auth.md`
    - `docs/research/things-cli-capabilities.md`
    - `docs/research/things-storage.md`
    - `docs/research/dida-mvp-command-api-mapping.md`
  - 初始化 Go 工程结构并接入 `cobra`。
  - 实现 OAuth 基础命令与任务核心命令：`projects/list/show/add/update/done/delete`。
- 下一次计划：
  - 补充 README、构建脚本和版本命令。
  - 补充 API client 基础单元测试。
  - 同步计划执行状态。

### 记录 #2 / commit `dfd66e2`
- 本次进展结果：
  - 补充 `README.md` 和 `Makefile`。
  - 增加 `dida version` 命令。
  - 补充 `internal/dida/client_test.go` 基础测试。
  - 通过 `go test ./...` 与 `go build ./cmd/dida`。
- 下一次计划：
  - 补充 `internal/cli` 命令层测试与 API 错误分支测试（401/403/404/5xx）。
  - 统一 `list/show` 输出参数风格（`--format` / `--json`）。
  - 推进 token 刷新策略设计与实现。

### 记录 #3 / in-progress
- 本次进展结果：
  - 按要求把进展日志迁移到 `CHANGELOG.md`，并停止继续改动 `PLAN.md`。
  - `list` 与 `show` 命令统一输出参数：新增 `--format table|json`，保留 `--json` 兼容。
  - 新增命令层单测：`internal/cli/output_format_test.go`。
  - 新增 API 错误分支测试：`internal/dida/client_error_test.go`（401/404/非法 JSON）。
  - 已通过验证：`go test ./...`、`go build ./cmd/dida`。
- 下一次计划：
  - 继续补充 `internal/cli` 命令行为测试（参数校验、dry-run 输出）。
  - 设计并实现 token 刷新（先 `auth refresh` 手动命令，再评估自动刷新）。
  - 增强文档中的认证故障排查与实际回调示例。

### 记录 #4 / in-progress
- 本次进展结果：
  - 继续统一输出能力：`list/show` 支持 `--format table|json`，并保持 `--json` 兼容。
  - 新增输出格式解析模块：`internal/cli/output_format.go`。
  - 新增命令层测试：`internal/cli/output_format_test.go`。
  - 新增 API 错误分支测试：`internal/dida/client_error_test.go`（401/404/响应解码异常）。
  - 新增手动 token 刷新能力：`dida auth refresh`。
  - 新增刷新流程测试：`internal/dida/client_refresh_test.go`。
  - 更新 README 命令列表，包含 `auth refresh`。
  - 已通过验证：`go test ./...`、`dida auth --help`。
- 下一次计划：
  - 增加 `internal/cli` 命令行为测试（参数必填校验、dry-run 输出断言）。
  - 评估并实现自动 refresh 策略（401 自动尝试刷新并重试一次）。
  - 补充认证排障文档（code 获取、redirect URI 不匹配、scope/授权失败案例）。

### 记录 #5 / in-progress
- 本次进展结果：
  - 新增认证实操文档：`docs/guides/dida-auth-token-flow.md`。
  - 文档包含完整“拿 token”流程：
    - 浏览器获取 `code`
    - `curl` 兑换 `access_token`
    - `curl` 验证 API 调用
    - `refresh_token` 刷新
    - 对应 `dida auth` 命令流程
  - 增加可直接复制的命令模板和常见排错项。
- 下一次计划：
  - 把该文档在 `README.md` 增加入口链接。
  - 继续补 `internal/cli` 命令行为测试（必填参数与 dry-run 输出）。

### 记录 #6 / in-progress
- 本次进展结果：
  - 默认配置路径策略改为固定 HOME 规则：`$HOME/.config/dida365-cli/config.json`。
  - 代码不再使用 `os.UserConfigDir()` 推导默认路径。
  - 同步更新文档说明（README 与认证实操文档）。
- 下一次计划：
  - 补充命令行为测试（含默认配置路径断言）。
  - 继续推进自动化 login（本地回调监听）设计。

### 记录 #7 / commit `c7e230c`
- 本次进展结果：
  - 完成 TODO 文档任务：建立 `doc/man` 层并新增 `doc/man/dida.1.md`。
  - `dida.1.md` 已覆盖当前命令规范（全局参数、auth、projects、任务命令、示例）。
  - 已把你实操成功的信息同步到文档语义：`auth init/login/token/status` 可直接使用。
- 下一次计划：
  - 回填 TODO 任务中的 commit hash。
  - 继续完善认证文档（刷新能力的文档表述与服务端能力边界）。

### 记录 #8 / in-progress
- 本次进展结果：
  - 参考 Dida OpenAPI（`POST /open/v1/project`）实现 `dida projects create`。
  - 新增项目创建参数：`--name`(required), `--color`, `--sort-order`, `--view-mode`, `--kind`, `--json`。
  - API client 新增 `CreateProject`，并补充对应单元测试。
  - 同步更新 README、man 文档与 MVP 映射文档。
- 下一次计划：
  - 补 CLI 命令级行为测试（create 参数校验与 dry-run）。
  - 继续完善认证流程与服务端能力边界说明。

### 记录 #9 / in-progress
- 本次进展结果：
  - 根据 OpenAPI Task 字段增强 `list/show` 默认展示。
  - `dida list` 现在默认展示：完成状态、due date、priority、content 预览。
  - `dida show` 现在默认展示：content、completed、completedTime、due date、priority。
  - 新增 `internal/cli/task_view.go` 与对应测试 `task_view_test.go`。
- 下一次计划：
  - 继续补命令级测试（`list/show` 输出与参数行为）。
  - 基于真实 API 响应进一步校准状态字段语义。

### 记录 #10 / in-progress
- 本次进展结果：
  - 完成代码/文档一致性检查（命令帮助与现有文档对齐）。
  - 新增 agent 高频手册：`docs/guides/agent-cli-quick-manual.md`。
  - 新增手册维护规则：`docs/guides/agent-manual-rules.md`。
  - README 增加 agent 文档入口。
- 下一次计划：
  - 按规则继续补命令行为测试。
  - 若命令行为再变化，优先更新 agent 手册再更新 README/man。

### 记录 #11 / in-progress
- 本次进展结果：
  - 补齐命令 `--json` 覆盖：`done`、`delete`、`version`、`auth` 各子命令（init/login/token/status/refresh/logout）。
  - 写操作在 `--dry-run --json` 下返回结构化预检输出。
  - `auth status --json` 返回统一状态对象，便于 agent 读取。
  - 同步更新 README、man 文档与 agent 快速手册。
- 下一次计划：
  - 回填 TODO 任务 commit hash。
  - 继续推进“展示字段枚举值而不是裸数字”。
