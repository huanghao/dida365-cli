# 待实现的功能

## 文档使用规则
- 我会手动往这个文件里添加任务
- agent从里面按顺序来考虑实现
- 完成后在对应的任务下标记DONE（不要删除已完成的任务），以及完成的时间，commit。然后更新CHANGELOG.md

## 任务

### 实现文档目录结构，完成dida.1.md

查看调研文件`doc/things-agent-docs.md`，里面提到了文档层次

> ### 建议 2：建立三层文档
> - `README.md`：用户快速上手
> - `doc/man/dida.1.md`：命令规范（持续对齐代码）
> - `docs/research/*.md`：调研与设计决策沉淀

DONE
- 完成时间：2026-02-23 00:32:50 +0800
- 完成内容：
  - 新建命令规范文档：`doc/man/dida.1.md`
  - 文档覆盖全局参数、auth 子命令、任务命令、示例与配置路径
  - 已体现当前实装状态：`auth init/login/token/status` 可直接使用（你已实操验证）
- commit：`c7e230c`

### 所有命令都实现--json参数，输出的结果给agent读

DONE
- 完成时间：2026-02-23 01:19:33 +0800
- 完成内容：
  - 补齐 `done`、`delete`、`version` 的 `--json` 输出
  - 补齐 `auth` 子命令 `init/login/token/status/refresh/logout` 的 `--json` 输出
  - `--dry-run` 在 JSON 模式下返回结构化预检对象
  - 同步更新 `README.md`、`doc/man/dida.1.md`、`docs/guides/agent-cli-quick-manual.md`
- commit：`65e8048`


### 展示字段的枚举值，而不是裸数字
比如下面的status、priority
```
❯ go run ./cmd/dida show --project 699b3180e4b03ce3a5fa07cf --id 699b328de4b03ce3a5fa13cf
Field          Value
-------------  ------------------------
ID             699b328de4b03ce3a5fa13cf
Project        699b3180e4b03ce3a5fa07cf
Title          some-task
Status         0
Completed      no
CompletedTime
Start
Due
Priority       0
Content        这是内容
Description
```

DONE
- 完成时间：2026-02-23 01:31:00 +0800
- 完成内容：
  - `list/show` 表格输出将 `status` 改为枚举标签（如 `incomplete/completed`）而非裸数字
  - `list/show` 表格输出将 `priority` 改为枚举标签（如 `none/low/medium/high`）而非裸数字
  - 新增展示层映射与测试：`internal/cli/task_view.go`, `internal/cli/task_view_test.go`
- commit：`ebabd94`

### 实现token refresh
- API文档上说grant_type目前只支持授权，但OAuth标准是支持fresh的。也许只需要验证一下当前的实现是否已经有效了。
- 如果无效，还要设计一下怎么才能实现。
- 检查一下目前每次获得的access token的默认过期时间有多长？（记录到文档里： `docs/guides/dida-auth-token-flow.md`）

DONE
- 完成时间：2026-02-23 01:44:30 +0800
- 完成内容：
  - 验证 refresh_token 路径：当前 app 返回 `Unauthorized grant type: refresh_token`，不可用
  - CLI 增加明确降级提示：`auth refresh` 失败时引导重新授权（`auth login` + `auth token --code ...`）
  - 文档补充验证结论与观测时长（`expires_in=15551999`，约 180 天）
  - 更新文档：`docs/guides/dida-auth-token-flow.md`、`docs/research/dida-api-overview.md`
- commit：待本次提交后回填

### 实现缓存，减少抖动
- 比如短时间调用多次查看命令，不需要都调用API
- 实现短时间缓存，先做设计，缓存到哪里，存多长时间
- 为了避免agent操作错误，比如短期快速调用多次创建，你要怎么设计一下防抖

### 整理展示的字段和API的对应关系
- 是否存在重要的字段没有被展示出来
- 是否存在代码里无用的字段

### 限制创建task时输入的内容长度（<500字）
- 限制创建时输入的长度，但不限制展示的长度。因为我还可以通过其他App在界面去修改，而不仅仅是这个工具

### 考虑其他适配agent的功能、限制、统计等功能

### 项目增加license信息
- 给我几个选项，应该用什么

### 打包和发布
- 命令行是不是用Homebrew来发布最好。我应该要先看什么准备信息？

### 面向agent来阅读这个项目，会需要考虑增加什么文件和导引？
