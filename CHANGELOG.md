# Changelog

## v0.0.1 - 2026-02-23

首个可用版本，完成 Dida365 CLI 的 OAuth、任务管理与 Homebrew 发布链路。

### Added
- 初始化 CLI 工程结构（Cobra）：`cmd/` + `internal/`。
- OAuth 命令：`auth init/login/token/status/refresh/logout`。
- 核心任务命令：`projects list`、`projects create`、`list`、`show`、`add`、`update`、`done`、`delete`。
- 输出能力：`list/show` 支持 `--format table|json`，并兼容 `--json`。
- 写操作支持 `--dry-run` 预检输出。

### Improved
- 默认展示字段可读性优化：状态与优先级展示为枚举标签。
- 输入校验：`add` 对 `title/content/desc` 做 `< 500` 字符限制（按 rune）。
- 本地防抖与缓存：写操作防抖，读操作短时缓存。

### Docs & Quality
- 补充 README、man 文档与多份 guides/research 文档。
- 增加单元测试（CLI 输出、错误分支、refresh 路径等），`go test ./...` 通过。

### Distribution
- 发布 GitHub Release：`v0.0.1`。
- 发布 Homebrew tap formula：`huanghao/homebrew-tap`（`dida`）。
