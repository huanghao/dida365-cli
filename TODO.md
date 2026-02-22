# 待实现的功能

## 文档使用规则
- 我会手动往这个文件里添加任务
- agent从里面按顺序来考虑实现
- 完成后在对应的任务下标记DONE（不要删除已完成的任务），以及完成的时间，commit

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
- commit：`5dbbc31`

### 实现缓存，减少抖动
- 比如短时间调用多次查看命令，不需要都调用API
- 实现短时间缓存，先做设计，缓存到哪里，存多长时间
- 为了避免agent操作错误，比如短期快速调用多次创建，你要怎么设计一下防抖

DONE
- 完成时间：2026-02-23 01:48:55 +0800
- 完成内容：
  - 按 `DG-0001 Decision: B` 落地：实现读缓存 + 写操作防抖
  - 为 `projects list`、`list`、`show` 增加 10 秒本地读缓存（`~/.config/dida365-cli/cache.json`）
  - 为 `add/update/done/delete/projects create` 增加 3 秒短窗口防抖
  - 防抖状态持久化到 `~/.config/dida365-cli/debounce.json`
  - 任意写操作成功后清空读缓存，避免读到陈旧数据
  - 新增全局禁用开关：`--no-cache` / `DIDA_NO_CACHE=1`
  - 写入成功后记录签名，窗口内重复请求直接拦截，避免 agent 误重复写
  - 新增测试：`internal/cli/write_debounce_test.go`、`internal/cli/read_cache_test.go`
- commit：`e0c4f09`

### 整理展示的字段和API的对应关系
- 是否存在重要的字段没有被展示出来
- 是否存在代码里无用的字段

DONE
- 完成时间：2026-02-23 01:43:40 +0800
- 完成内容：
  - 新增字段映射文档：`docs/research/dida-task-field-mapping.md`（API 字段 -> list/show/json 对照）
  - `show` 表格补充高价值字段：`AllDay`、`TimeZone`、`Repeat`
  - 结论：高频字段已覆盖；当前无应直接删除的“无用字段”，未展示字段保留给后续能力
- commit：`07d1bf8`

### 限制创建task时输入的内容长度（<500字）
- 限制创建时输入的长度，但不限制展示的长度。因为我还可以通过其他App在界面去修改，而不仅仅是这个工具

DONE
- 完成时间：2026-02-23 01:45:47 +0800
- 完成内容：
  - `dida add` 增加输入长度校验：`--title`、`--content`、`--desc` 必须 `< 500` 字符
  - 长度按字符数（rune）计算，支持中文
  - 新增测试：`internal/cli/task_input_test.go`（499/500 边界）
  - 更新文档：`README.md`、`doc/man/dida.1.md`、`docs/guides/agent-cli-quick-manual.md`
- commit：`c0ff3c7`

### 考虑其他适配agent的功能、限制、统计等功能
- 设计闸门补充：已在 `DG-0002` 明确 stats 使用场景与指标草案，待你最终选择 A/B/C

### 项目增加license信息
- 给我几个选项，应该用什么

DONE
- 完成时间：2026-02-23 02:03:28 +0800
- 完成内容：
  - 根据 `DG-0003 Decision: A` 落地 MIT license
  - 新增 `LICENSE`
  - `README.md` 增加 License 段落
- commit：`a648f91`

### 打包和发布
- 命令行是不是用Homebrew来发布最好。我应该要先看什么准备信息？

DONE
- 完成时间：2026-02-23 02:03:28 +0800
- 完成内容：
  - 新增发布指南：`docs/guides/release-homebrew.md`
  - 覆盖 Homebrew 发布准备清单、最小步骤、formula 示例、当前缺口与建议顺序
- commit：`a648f91`

### 面向agent来阅读这个项目，会需要考虑增加什么文件和导引？

DONE
- 完成时间：2026-02-23 02:03:28 +0800
- 完成内容：
  - 新增 `docs/guides/agent-project-onboarding.md`
  - 明确新 agent 进入项目后的最小必读文件、执行规则与常用命令
- commit：`a648f91`

### 实现due date
- 创建任务的时候支持截止时间
- 支持更新一个任务的due date

### 实现重复任务
- 检查OPEN API是否支持创建重复任务，比如每周一早上8点之类的。如果支持，就做一下接口设计，然后支持它