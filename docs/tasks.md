# Repo Manager 任务清单 (docs/tasks.md)

> 最后更新：**2026-11-25**
> 登记记录：2026-11-25 依据 [ARCHITECTURE.md §7](./ARCHITECTURE.md) 评估快照登记 S1-S9 优化项（文档系统性完善的配套评估结论）；责任人默认为维护者 + AI 协作会话，计划完成时间见下表，完成后归档至 [tasks-done.md](./tasks-done.md)。
> 维护约定：已完成任务归档至 [tasks-done.md](./tasks-done.md)；本文件仅保留进行中 / 待办任务。
> 全局版本唯一来源：`internal/buildinfo/version.go` 的 `Version`（当前 **2.3.7**）。
> 本轮处理：2026-11-24 11:00 依据 `reports/code-review-2026-11-23.md` 的 12 项发现登记 A1–A10，逐项复核并完成其中 10 项；P2-1（繁中日志误报）与 P3-4（testify 分类无需改动）直接归档至 [tasks-done.md 第十节](./tasks-done.md)。
> 收口记录：2026-11-24 11:20 全部 12 项审查发现处理完毕（10 项修复 + 2 项核验归档），A1–A10 归档至 [tasks-done.md 第十一节](./tasks-done.md)；**当前无未完成事项**。全量门禁（root-clean / gofmt / vet / test / build nogui / build ci / check:docs）实跑通过，`-race` 因本机无 gcc 纳入 CI 验证。
> 核查记录：2026-11-24 11:33 按指令复核任务进度——无已完成任务需归档、无未完成任务需调整优先级 / 状态、无剩余任务可推进；本次仅同步时间戳并修正 `tasks-done.md` 头部时间戳漂移，任务状态维持不变。
> 核查记录：2026-11-24 15:01 再次按指令复核——结论同上（无待办 / 无归档 / 无推进），仅同步时间戳；2 个文档文件的时间戳改动已提交并推送入库，`origin/main` 同步。

---

## 进行中 / 待办

| ID | 任务 | 优先级 | 责任人 | 计划完成 | 状态 | 来源与备注 |
|----|------|:------:|--------|----------|------|------------|
| S1 | 日志轮转：`repo-manager.log` 超 5 MB 归档为 `.log.1`（仅留 1 份） | P2 | 维护者 + AI 协作会话 | 2026-12-15 | 待办 | [ARCHITECTURE.md §7.1](./ARCHITECTURE.md)；日志占用收敛为约 10 MB 上限 |
| S2 | Windows 环境变量改注册表直写（`golang.org/x/sys/windows/registry`），消除 `setx` 截断风险 | P2 | 维护者 + AI 协作会话 | 2026-12-15 | 待办 | 同上 §7.1；保留 `setx` 降级路径 |
| S3 | `EXTRA_CLEAN_DIRS` 支持受限相对子路径（逐段校验，默认关闭） | P3 | 维护者 | 择期（随需求） | 待办 | 同上 §7.1 |
| S4 | GUI 组件级自动化测试（`fyne.io/fyne/v2/test` 软件驱动）纳入 CI | P2 | 维护者 + AI 协作会话 | 2026-12-31 | 待办 | 同上 §7.1；替代部分 §9.3 人工冒烟 |
| S5 | GUI 执行中可取消（`context.Context` 贯通 runner → engine） | P2 | 维护者 + AI 协作会话 | 2026-12-31 | 待办 | 同上 §7.2 |
| S6 | 项目扫描深度可配置（`SCAN_DEPTH`，默认 3、上限 8） | P3 | 维护者 | 择期（随需求） | 待办 | 同上 §7.2 |
| S7 | CLI `-json` 结构化输出（退出码不变），便于脚本 / CI 集成 | P3 | 维护者 | 择期（随需求） | 待办 | 同上 §7.2；为未来 API 化铺垫 |
| S8 | 系统级安装包评估（`.msi` / `.dmg`，`fyne-cross`），CI 产物追加 | P3 | 维护者 | 择期（随需求） | 待办 | 同上 §7.2 |
| S9 | 源码文件头注释版本随改动渐进同步（禁止全仓刷写） | P3 | 维护者 + AI 协作会话 | 随改动持续 | 待办 | 同上 §7.3；按 [CONTRIBUTING.md §8.2](./CONTRIBUTING.md) 执行 |

> 已知限制（`windows/arm64` 无交叉工具链、`linux/arm64` 缺 arm64 架构开发库）登记于 [MAINTAINERS.md §6](./MAINTAINERS.md)，不属待办。本机无 gcc，`go test -race` 已纳入 `ci.yml`（ubuntu runner）验证，不属待办。

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
