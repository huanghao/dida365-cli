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
