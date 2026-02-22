# things-cli 面向 Agent 的文档调研

## 调研范围
- 仓库：`~/workspace/sources/things3-cli`
- 重点：对 agent/协作流程有直接指导价值的文档

## 文档清单与作用

### 1) `AGENTS.md`
- 作用：仓库级执行规范，明确项目结构、测试命令、代码风格、提交/PR 约定、权限注意事项。
- 对 agent 的价值：
  - 明确“改哪里”：`cmd/`、`internal/cli`、`internal/db`、`internal/things`。
  - 明确“怎么验”：`make test`、`go test ./...`。
  - 明确“风险点”：
    - 更新命令依赖 `THINGS_AUTH_TOKEN`
    - DB 读取可能需要 Full Disk Access
    - AppleScript 可能触发系统权限弹窗
  - 明确“联动要求”：CLI 变更后需要同步更新 agent scripts/skills。

### 2) `README.md`
- 作用：用户视角的能力总览与操作入口。
- 对 agent 的价值：
  - 给出完整命令集合（可作为回归清单）。
  - 给出认证、数据库、重复任务的关键行为与约束。
  - 给出命令行为差异点（例如 `--foreground`、`--dry-run`、删除确认策略）。

### 3) `doc/man/things.1.md` + `share/man/man1/things.1`
- 作用：结构化 CLI 规格文档（命令说明、参数、示例）。
- 对 agent 的价值：
  - 命令行为“准规范”，尤其适合生成帮助文本、E2E 校验脚本。
  - 对参数兼容（别名、必填、优先级）提供权威描述。

### 4) `docs/RELEASING.md`、`scripts/*.sh`、`Formula/*.rb`
- 作用：发布流程和分发脚本。
- 对 agent 的价值：
  - 约束版本发布动作，不属于核心实现逻辑，但影响交付闭环。

## 文档组织方式观察
- 分层清晰：
  - `AGENTS.md`：贡献/协作约束
  - `README.md`：用户入口
  - `man` 文档：命令规格
  - `docs/` + `scripts/`：工程流程
- 与代码一一对应：命令文件名与子命令名对齐（如 `add.go`、`update.go`）。
- 文档强调“可执行性”：每项能力都能映射到命令、参数、环境变量或测试。

## 对 dida365-cli 的可迁移建议

### 建议 1：保留仓库级 `AGENTS.md`
- 最小应包含：
  - 目录结构说明
  - 构建/测试命令
  - OAuth 凭证与本地配置约定
  - 调试与错误排查入口

### 建议 2：建立三层文档
- `README.md`：用户快速上手
- `doc/man/dida.1.md`：命令规范（持续对齐代码）
- `docs/research/*.md`：调研与设计决策沉淀

### 建议 3：让“调研文档 -> 实现任务”可追踪
- 每份调研文档末尾追加：
  - 结论
  - 待实现项
  - 风险项
  - 验证方式

## 结论
`things-cli` 的可复用价值不只是命令设计，而是“文档即执行约束”的组织方式。`dida365-cli` 建议从 Day 1 就采用同样分层，避免后续行为漂移。
