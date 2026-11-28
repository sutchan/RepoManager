# Repo Manager 任务清单 (docs/tasks.md)

> 最后更新：**2026-11-28**（v2.3.14 安全加固与文档对齐：符号链接越界防护、功能 1 删除边界、预演零副作用、GUI 执行期交互修复；活动表不变）
> 登记记录：2026-11-25 依据 [ARCHITECTURE.md §7](./ARCHITECTURE.md) 快照登记 S1–S9；2026-11-26 依据同文 §7.4 补充评估登记 S10–S15，并新增「待验证技术细节」T1–T5。责任人默认为维护者（@sutchan）+ AI 协作会话，计划完成时间为明确日期；完成后归档至 [tasks-done.md](./tasks-done.md)。
> 维护约定：已完成任务归档至 [tasks-done.md](./tasks-done.md)；本文件仅保留进行中 / 待办任务。
> 全局版本唯一来源：`internal/buildinfo/version.go` 的 `Version`（当前 **2.3.14**）。
> 本轮处理：2026-11-24 11:00 依据 `reports/code-review-2026-11-23.md` 的 12 项发现登记 A1–A10，逐项复核并完成其中 10 项；P2-1（繁中日志误报）与 P3-4（testify 分类无需改动）直接归档至 [tasks-done.md 第十节](./tasks-done.md)。
> 收口记录：2026-11-24 11:20 全部 12 项审查发现处理完毕（10 项修复 + 2 项核验归档），A1–A10 归档至 [tasks-done.md 第十一节](./tasks-done.md)；该时点无未完成事项（**当前待办以下表为准**）。全量门禁（root-clean / gofmt / vet / test / build nogui / build ci / check:docs）实跑通过，`-race` 因本机无 gcc 纳入 CI 验证。
> 核查记录：2026-11-24 11:33 按指令复核任务进度——无已完成任务需归档、无未完成任务需调整优先级 / 状态、无剩余任务可推进；本次仅同步时间戳并修正 `tasks-done.md` 头部时间戳漂移，任务状态维持不变。
> 核查记录：2026-11-24 15:01 再次按指令复核——结论同上（无待办 / 无归档 / 无推进），仅同步时间戳；2 个文档文件的时间戳改动已提交并推送入库，`origin/main` 同步。
> 核查记录：2026-11-26 全仓文档一致性审计——实算 `*_test.go` 共 14 个（含 1 个辅助文件）、`Test*` 函数 42 个、覆盖 8 个包，同步修正 [OPENSPEC.md §9.3](./OPENSPEC.md) 与 [PRODUCTION-READINESS.md §2.1](./PRODUCTION-READINESS.md) 的过时计数（10/31/6 包）；核实 `-tags ci` 为 Fyne 内置标签（`app/app_software.go` 于 `ci` 下改用 headless 测试驱动 + software painter，均纯 Go 无 cgo），`TROUBLESHOOTING.md` / `ARCHITECTURE.md` / `CONTRIBUTING.md` / `OPENSPEC.md` / `docs/README.md` / `MAINTAINERS.md` 相关描述属实无误；Release 资产命名与 `scripts/package-release.sh` 实况一致（amd64→x64），无冲突。任务状态不变，本轮无新增待办。
> 本轮记录：2026-11-26 系统性文档完善——版本 bump 至 v2.3.8、修复配置读取失败退出码、补齐 OpenAPI 3.1 CLI 契约映射（`CLI-REFERENCE.md` 附录 A）、FAQ 快速索引（`TROUBLESHOOTING.md` §0）与角色权限矩阵（`USER-GUIDE.md` §1.1）；登记 S10–S15 与 T1–T5，并将 S3 / S6–S9 的「择期」改为明确日期。历史「当前无未完成事项」仅指 2026-11-24 时点。
> 核验记录：2026-11-26 v2.3.9 系统性文档标准化核验——6 大内容模块与质量标准逐项对照通过；修正 3 处把平台仓库名误写为 Git 远程的表述（`docs/README.md` §1.2、`CONTRIBUTING.md` §3.2 规则 8、`landing/README.md` §2）与 1 处「当前无未完成项」过时表述；版本同步至 2.3.9。无新增待办，S1–S15 / T1–T5 状态不变。
> 本轮记录（v2.3.11 / 2026-11-26）：实施完成度审查的 P0–P3 改进——S10 局部失败汇总（返回非零码）、S11 npm 安装超时（`REPO_MANAGER_NPM_TIMEOUT`）、7 处文档漂移修复、S18 CI 未引用符号检测（强制门禁）+ 资源引用扫描 + 覆盖率采集；同步版本与 CHANGELOG。实算测试 18 文件 / 55 函数 / 8 包；覆盖率基线总体 49.9%（核心包 `engine` 74.8% / `config` 81.5% / `i18n` 91.7%），距 80% 门槛仍有差距，已登记为分阶段目标。
> 本轮记录（v2.3.12 / 2026-11-26）：S4 GUI 组件测试、S5 GUI 取消上下文、S17 平台错误本地化完成并归档；S7 活动表状态与已归档实现对齐；T4 / T5 验证结论归档。已完成批次全文见 [tasks-done.md 第十五节](./tasks-done.md)。
> 本轮记录（v2.3.13 / 2026-11-28）：完善 `landing/` 展示页——补齐目标受众、演示、技术栈、更新记录、联系方式板块，三语词条拆到 `dict.js`，补响应式 / 可访问性 / 加载性能约束并扩写 `landing/README.md` 维护说明；无新增待办，S9 与 T1–T3 状态不变。
> 本轮记录（v2.3.14 / 2026-11-28）：全仓代码审查修复——符号链接越界防护（根目录 / 清理目标 / `.next` / `.npmrc`）、功能 1 只删 Node 项目依赖、预演不回写配置、配置原子保存、GUI 执行期禁用语言切换与进度误报、递归深度日志、原始字典校验；补 4 个回归测试文件（现 24 文件 / 72 函数 / 8 包），拆分 `clean.go` / `config.go` 至 ≤200 行；前端补语义化 ID、ARIA 与令牌对齐，移除落地页非页脚版本号；同步 CLI 契约、OPENSPEC、DESIGN-SYSTEM、CHANGELOG 与全部版本展示位。任务状态不变，S9 与 T1–T3 维持原优先级。

---

## 进行中 / 待办

| ID | 任务 | 优先级 | 责任人 | 计划完成 | 状态 | 来源与备注 |
|----|------|:------:|--------|----------|------|------------|
| S9 | 源码文件头注释版本随改动渐进同步（禁止全仓刷写） | P3 | 维护者 + AI 协作会话 | 2026-12-31 | 进行中（持续任务） | 本轮已改动文件的头注释同步至 v2.3.12；按 [CONTRIBUTING.md §8.2](./CONTRIBUTING.md) 禁止全仓刷写，剩余未改动文件随后续改动渐进同步 |

> 已知限制（`windows/arm64` 无交叉工具链、`linux/arm64` 缺 arm64 架构开发库）登记于 [MAINTAINERS.md §6](./MAINTAINERS.md)，不属待办。本机无 gcc，`go test -race` 已纳入 `ci.yml`（ubuntu runner）验证，不属待办。

### 待验证技术细节

| ID | 待验证项 | 优先级 | 验证方法 | 责任人 | 计划完成 | 状态 |
|----|----------|:------:|----------|--------|----------|------|
| T1 | Windows runner `CGO_ENABLED=1` 是否自带可用 gcc | P2 | `workflow_dispatch` 触发发布工作流打包，检查 `windows/amd64` 构建日志与归档产物 | 维护者 + AI 协作会话 | 2026-12-10 | 待验证（受阻：需 workflow_dispatch 权限 / tag 授权） |
| T3 | tag push 后 `release` job 的 `gh api` 轮询能否等到同 SHA 的 `ci.yml` 结论 | P2 | 首次 `vX.Y.Z` tag 实跑，比对 release / ci 运行时序与结论 | 发布者 @sutchan | 2026-12-20 | 待验证（受阻：需打 tag 触发发版，待用户授权） |
| T2 | `macos-latest` 上 darwin/amd64 交叉编译的 SDK 可用性 | P3 | 同 T1，核对两项 macOS 归档是否生成 | 维护者 + AI 协作会话 | 2026-12-10 | 待验证（受阻：需 workflow_dispatch 权限 / tag 授权） |

> 待验证项与待办的区别：待办是明确的代码 / 文档改动任务；待验证项是对既有 CI 与运行环境事实的确认。T4 / T5 已完成验证并归档至 [tasks-done.md 第十五节](./tasks-done.md)。验证失败即转化为新的 `S` 编号任务并按优先级排期。

## 已完成归档索引

| 批次 | ID | 归档位置 | 迁移时间 |
|------|----|----------|----------|
| v2.3.12 | S4 / S5 / S17、T4 / T5（S7 状态修正） | [tasks-done.md 第十五节](./tasks-done.md) | 2026-11-26 23:14 |
| v2.3.11 | S1 / S2 / S3 / S6 / S7 / S8 / S10 / S11 / S12 / S13 / S14 / S15 / S16 / S18 | [tasks-done.md 第十四节](./tasks-done.md) | 2026-11-26 17:12 |
| v2.3.10 | C1–C8 | [tasks-done.md 第十三节](./tasks-done.md) | 2026-11-26 16:23 |
| v2.3.8 | D1–D5 | [tasks-done.md 第十二节](./tasks-done.md) | 2026-11-26 16:23 |
| v2.3.6 | A1–A10 | [tasks-done.md 第十一节](./tasks-done.md) | 历史归档 |
| v2.3.5 | RV1–RV3 | [tasks-done.md 第十节](./tasks-done.md) | 历史归档 |

> 已完成任务的完整描述、状态与上下文统一在 [tasks-done.md](./tasks-done.md) 维护；本文件只保留索引，避免双份正文漂移。
