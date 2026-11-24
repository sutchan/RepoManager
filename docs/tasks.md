# Repo Manager 任务清单 (docs/tasks.md)

> 最后更新：**2026-11-24 15:01**
> 维护约定：已完成任务归档至 [tasks-done.md](./tasks-done.md)；本文件仅保留进行中 / 待办任务。
> 全局版本唯一来源：`internal/buildinfo/version.go` 的 `Version`（当前 **2.3.6**）。
> 本轮处理：2026-11-24 11:00 依据 `reports/code-review-2026-11-23.md` 的 12 项发现登记 A1–A10，逐项复核并完成其中 10 项；P2-1（繁中日志误报）与 P3-4（testify 分类无需改动）直接归档至 [tasks-done.md 第十节](./tasks-done.md)。
> 收口记录：2026-11-24 11:20 全部 12 项审查发现处理完毕（10 项修复 + 2 项核验归档），A1–A10 归档至 [tasks-done.md 第十一节](./tasks-done.md)；**当前无未完成事项**。全量门禁（root-clean / gofmt / vet / test / build nogui / build ci / check:docs）实跑通过，`-race` 因本机无 gcc 纳入 CI 验证。
> 核查记录：2026-11-24 11:33 按指令复核任务进度——无已完成任务需归档、无未完成任务需调整优先级 / 状态、无剩余任务可推进；本次仅同步时间戳并修正 `tasks-done.md` 头部时间戳漂移，任务状态维持不变。
> 核查记录：2026-11-24 15:01 再次按指令复核——结论同上（无待办 / 无归档 / 无推进），仅同步时间戳；2 个文档文件的时间戳改动已提交并推送入库，`origin/main` 同步。

---

## 进行中 / 待办

| ID | 任务 | 优先级 | 状态 | 更新时间 | 备注 |
|----|------|--------|------|----------|------|
| — | 当前无待办 | — | — | 2026-11-24 11:20 | 新增任务须先登记于此并按优先级排序；完成后移入 [tasks-done.md](./tasks-done.md) |

## 本轮完成概览（v2.3.6）

| ID | 任务 | 来源 | 归档 |
|----|------|------|------|
| A1 | 配置加载失败保护（`main.go` 错误顺序 + `config.Save` 读旧错误） | P1-1 / P1-2 | [tasks-done.md 第十一节](./tasks-done.md) |
| A2 | GUI 功能按钮互斥（跨重绘 `atomic.Bool` + 执行期间禁用按钮） | P2-3 | 同上 |
| A3 | `NormalizeRoot` 文件系统根退化 + `ValidateRoot` 拒绝根目录 | P2-2 | 同上 |
| A4 | GUI 日志有界缓冲，消除 O(n²) | P2-4 | 同上 |
| A5 | 补 `cmd/` 与核心逻辑测试 + CI `-race` 门禁 | P2-5 | 同上 |
| A6 | 落地页版本展示同步 | P3-1 | 同上 |
| A7 | workflow 头注释版本同步 | P3-2 | 同上 |
| A8 | GUI `render` 闭包拆分（`runner.go` + `newLangSelect`） | P3-3 | 同上 |
| A9 | `.next` 缓存路径辅助函数 + `Run` 前置条件注释 | P3-5 | 同上 |
| A10 | 版本 bump 2.3.6 并按 §8.2 同步全部展示位 | §8.1 / §8.2 | 同上 |

> 已知限制（`windows/arm64` 无交叉工具链、`linux/arm64` 缺 arm64 架构开发库）登记于 [MAINTAINERS.md §6](./MAINTAINERS.md)，不属待办。本机无 gcc，`go test -race` 已纳入 `ci.yml`（ubuntu runner）验证，不属待办。
