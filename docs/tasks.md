# Repo Manager 任务清单 (docs/tasks.md)

> 最后更新：**2026-11-26 16:23**（已完成任务已迁移至 [tasks-done.md](./tasks-done.md) 第十二、十三节，本文件仅保留归档索引）
> 登记记录：2026-11-25 依据 [ARCHITECTURE.md §7](./ARCHITECTURE.md) 快照登记 S1–S9；2026-11-26 依据同文 §7.4 补充评估登记 S10–S15，并新增「待验证技术细节」T1–T5。责任人默认为维护者（@sutchan）+ AI 协作会话，计划完成时间为明确日期；完成后归档至 [tasks-done.md](./tasks-done.md)。
> 维护约定：已完成任务归档至 [tasks-done.md](./tasks-done.md)；本文件仅保留进行中 / 待办任务。
> 全局版本唯一来源：`internal/buildinfo/version.go` 的 `Version`（当前 **2.3.11**）。
> 本轮处理：2026-11-24 11:00 依据 `reports/code-review-2026-11-23.md` 的 12 项发现登记 A1–A10，逐项复核并完成其中 10 项；P2-1（繁中日志误报）与 P3-4（testify 分类无需改动）直接归档至 [tasks-done.md 第十节](./tasks-done.md)。
> 收口记录：2026-11-24 11:20 全部 12 项审查发现处理完毕（10 项修复 + 2 项核验归档），A1–A10 归档至 [tasks-done.md 第十一节](./tasks-done.md)；该时点无未完成事项（**当前待办以下表 S1–S15 为准**）。全量门禁（root-clean / gofmt / vet / test / build nogui / build ci / check:docs）实跑通过，`-race` 因本机无 gcc 纳入 CI 验证。
> 核查记录：2026-11-24 11:33 按指令复核任务进度——无已完成任务需归档、无未完成任务需调整优先级 / 状态、无剩余任务可推进；本次仅同步时间戳并修正 `tasks-done.md` 头部时间戳漂移，任务状态维持不变。
> 核查记录：2026-11-24 15:01 再次按指令复核——结论同上（无待办 / 无归档 / 无推进），仅同步时间戳；2 个文档文件的时间戳改动已提交并推送入库，`origin/main` 同步。
> 核查记录：2026-11-26 全仓文档一致性审计——实算 `*_test.go` 共 14 个（含 1 个辅助文件）、`Test*` 函数 42 个、覆盖 8 个包，同步修正 [OPENSPEC.md §9.3](./OPENSPEC.md) 与 [PRODUCTION-READINESS.md §2.1](./PRODUCTION-READINESS.md) 的过时计数（10/31/6 包）；核实 `-tags ci` 为 Fyne 内置标签（`app/app_software.go` 于 `ci` 下改用 headless 测试驱动 + software painter，均纯 Go 无 cgo），`TROUBLESHOOTING.md` / `ARCHITECTURE.md` / `CONTRIBUTING.md` / `OPENSPEC.md` / `docs/README.md` / `MAINTAINERS.md` 相关描述属实无误；Release 资产命名与 `scripts/package-release.sh` 实况一致（amd64→x64），无冲突。任务状态不变，本轮无新增待办。
> 本轮记录：2026-11-26 系统性文档完善——版本 bump 至 v2.3.8、修复配置读取失败退出码、补齐 OpenAPI 3.1 CLI 契约映射（`CLI-REFERENCE.md` 附录 A）、FAQ 快速索引（`TROUBLESHOOTING.md` §0）与角色权限矩阵（`USER-GUIDE.md` §1.1）；登记 S10–S15 与 T1–T5，并将 S3 / S6–S9 的「择期」改为明确日期。历史「当前无未完成事项」仅指 2026-11-24 时点。
> 核验记录：2026-11-26 v2.3.9 系统性文档标准化核验——6 大内容模块与质量标准逐项对照通过；修正 3 处把平台仓库名误写为 Git 远程的表述（`docs/README.md` §1.2、`CONTRIBUTING.md` §3.2 规则 8、`landing/README.md` §2）与 1 处「当前无未完成项」过时表述；版本同步至 2.3.9。无新增待办，S1–S15 / T1–T5 状态不变。
> 本轮记录（v2.3.11 / 2026-11-26）：实施完成度审查的 P0–P3 改进——S10 局部失败汇总（返回非零码）、S11 npm 安装超时（`REPO_MANAGER_NPM_TIMEOUT`）、7 处文档漂移修复、S18 CI 未引用符号检测（告警模式）+ 资源引用扫描 + 覆盖率采集；同步版本与 CHANGELOG。实算测试 18 文件 / 55 函数 / 8 包；覆盖率基线总体 49.9%（核心包 `engine` 74.8% / `config` 81.5% / `i18n` 91.7%），距 80% 门槛仍有差距，已登记为分阶段目标。

---

## 进行中 / 待办

| ID | 任务 | 优先级 | 责任人 | 计划完成 | 状态 | 来源与备注 |
|----|------|:------:|--------|----------|------|------------|
| S1 | 日志轮转：`repo-manager.log` 超 5 MB 归档为 `.log.1`（仅留 1 份） | P2 | 维护者 + AI 协作会话 | 2026-12-15 | **已完成（v2.3.11）** | `internal/logs/logs.go` 写入前按 5 MB 阈值轮转并保留 1 份归档；`logs_test.go` 覆盖轮转边界与下游转发 |
| S2 | Windows 环境变量改注册表直写（`golang.org/x/sys/windows/registry`），消除 `setx` 截断风险 | P2 | 维护者 + AI 协作会话 | 2026-12-15 | **已完成（v2.3.11）** | `platform_windows.go` 直写 `HKCU\Environment`（`REG_SZ`），删除值不存在视为成功；失败回退 `setx` / `reg`；`go mod tidy` 将 `x/sys` 提升为直接依赖 |
| S4 | GUI 组件级自动化测试（`fyne.io/fyne/v2/test` 软件驱动）纳入 CI | P2 | 维护者 + AI 协作会话 | 2026-12-31 | 待办 | 同上 §7.1；替代部分 §9.3 人工冒烟 |
| S5 | GUI 执行中可取消（`context.Context` 贯通 runner → engine） | P2 | 维护者 + AI 协作会话 | 2026-12-31 | 待办 | 同上 §7.2 |
| S7 | CLI `-json` 结构化输出（退出码不变），便于脚本 / CI 集成 | **P2（重评↑）** | 维护者 @sutchan | 2027-01-31 | 待办 | 同上 §7.2；S10 修复退出码后脚本集成价值上升 |
| S10 | 危险操作局部失败汇总并返回非零码（对齐功能 5 契约） | **P1（重评↑）** | 维护者 + AI 协作会话 | 2026-12-31 | **已完成（v2.3.11）** | [ARCHITECTURE.md §7.4](./ARCHITECTURE.md)；`failures.go` 汇总器 + `engine.go` / `clean.go` 全功能覆盖 + 专项测试已落地 |
| S11 | `npm install -g` 增加 `context.WithTimeout`（默认 10 min）与可配置超时 | P2 | 维护者 + AI 协作会话 | 2026-12-31 | **已完成（v2.3.11）** | 同上 §7.4；`install.go` 用 `exec.CommandContext`，`REPO_MANAGER_NPM_TIMEOUT`（秒）可覆盖，超时显式报错 |
| S13 | i18n 字段级回退 + 三语键 parity 测试纳入 CI | **P1（重评↑）** | 维护者 + AI 协作会话 | 2027-01-15 | **进行中（并行会话实现中）** | 同上 §7.4；`i18n.go` / `log.go`；缺字段回退会产生空串 |
| S17 | 平台层错误文案硬编码中文，未走 i18n 字典 | **P2（重评↑）** | 维护者 + AI 协作会话 | 2027-01-31 | 待办 | 冗余清理实查；`internal/platform/platform.go` 的 `ValidateRoot` / `globalNodeModulesPath`；与 S13 合并推进 |
| S3 | `EXTRA_CLEAN_DIRS` 支持受限相对子路径（逐段校验，默认关闭） | P3 | 维护者 @sutchan | 2027-01-31 | **已完成（v2.3.11）** | `platform.IsSafeCleanDirPath` 逐段校验；config / engine 双层防线；`dirs_test.go` 与 `keys_test.go` 覆盖穿越输入 |
| S6 | 项目扫描深度可配置（`SCAN_DEPTH`，默认 3、上限 8） | P3 | 维护者 @sutchan | 2027-01-31 | **已完成（v2.3.11）** | `config.ScanDepth` + `-scan-depth` flag；`listDirs` 边界钳制；`TestScanDepthConfig` 覆盖解析/回落/往返 |
| S8 | 系统级安装包评估（`.msi` / `.dmg`，`fyne-cross`），CI 产物追加 | P3 | 维护者 @sutchan | 2027-02-28 | **已完成（v2.3.11）** | 评估结论落地 [MAINTAINERS.md §6](./MAINTAINERS.md)：签名证书 / 公证成本高于收益，维持仅归档分发；触发条件与后续 S 编号机制已记录 |
| S9 | 源码文件头注释版本随改动渐进同步（禁止全仓刷写） | P3 | 维护者 + AI 协作会话 | 2026-12-31 | 待办 | 同上 §7.3；按 [CONTRIBUTING.md §8.2](./CONTRIBUTING.md) 执行；每月末复核 |
| S12 | 非法功能号 / GUI 路径参数可见化，补 `-version` flag | **P3（重评↓）** | 维护者 + AI 协作会话 | 2027-01-15 | **已完成（v2.3.11）** | `main.go` 新增 `-version`、非法功能号报错与 GUI 位置参数提示；帮助文本同步 |
| S14 | `src/lib` i18n 与数据层增加 Vitest 逻辑测试 | P3 | 维护者 @sutchan | 2027-02-28 | **已完成（v2.3.11）** | 技术选型改为零依赖 `node:test`：`scripts/i18n-parity.test.mjs` 校验三语键集合、`SUPPORTED_LANGS` 顺序与 `HTML_LANG` 映射，`npm run check:i18n` 已绿；避免为弃用资产引入 Vitest 依赖 |
| S15 | 统一本地 6 目标 / CI 4 目标构建矩阵口径 | P3 | 维护者 @sutchan | 2027-02-28 | **已完成（v2.3.11）** | `ARCHITECTURE.md` §2 与 [INSTALL.md §4](./INSTALL.md) 已加口径说明 |
| S16 | `config.Load` / `config.Save` 的 `key=value` 解析逻辑重复，抽取共用解析器 | P3 | 维护者 + AI 协作会话 | 2027-01-31 | **已完成（v2.3.11）** | `parseKeyValue` 统一跳过注释/空行并归一化 key，Load / Save 均改用它 |
| S18 | 冗余预防门禁：CI 增加未引用符号检测与静态资源引用扫描 | **P3（重评↓）** | 维护者 + AI 协作会话 | 2027-01-15 | **已完成（v2.3.11）** | `ci.yml` 新增 `staticcheck -checks U1000 ./cmd/... ./internal/...`（固定版本 2025.1.1）与 `npm run check:asset-refs`；清单见 [reports/cleanup-report-2026-11-26.md](../reports/cleanup-report-2026-11-26.md) |

> 已知限制（`windows/arm64` 无交叉工具链、`linux/arm64` 缺 arm64 架构开发库）登记于 [MAINTAINERS.md §6](./MAINTAINERS.md)，不属待办。本机无 gcc，`go test -race` 已纳入 `ci.yml`（ubuntu runner）验证，不属待办。

### 本轮重评记录（2026-11-26 16:23）

| ID | 原优先级 / 状态 | 新优先级 / 状态 | 变更理由 |
|----|----------------|----------------|----------|
| S7 | P3 / 待办 | P2 / 待办 | D1 修复退出码后，脚本依赖可靠退出状态；结构化输出价值上升 |
| S10 | P2 / 待办 | P1 / 进行中 | 局部失败静默返回 0 属正确性缺陷；并行会话已落地主体实现，本轮补测试与验证 |
| S12 | P2 / 待办 | P3 / 待办 | 参数可见化与 `-version` 属体验优化，不影响数据正确性 |
| S13 | P2 / 待办 | P1 / 进行中 | 整字典回退会让漏译字段显示空串，属用户可见缺陷；并行会话实现中 |
| S15 | P3 / 待办 | P3 / 已完成（待归档） | 口径说明已于 v2.3.8 落在 `ARCHITECTURE.md` §2 与 `INSTALL.md` §4，待统一收口归档 |
| S17 | P3 / 待办 | P2 / 待办 | 硬编码中文错误与 S13 字段回退同源，合并推进可复用同一测试 |
| S18 | P2 / 待办 | P3 / 待办 | 预防性门禁，排在正确性修复之后 |

### 待验证技术细节

| ID | 待验证项 | 优先级 | 验证方法 | 责任人 | 计划完成 | 状态 |
|----|----------|:------:|----------|--------|----------|------|
| T1 | Windows runner `CGO_ENABLED=1` 是否自带可用 gcc | P2 | `workflow_dispatch` 触发发布工作流打包，检查 `windows/amd64` 构建日志与归档产物 | 维护者 + AI 协作会话 | 2026-12-10 | 待验证 |
| T4 | `go test -race` 在 CI ubuntu runner 的实际结果（本机无 gcc 无法复现） | P2 | 查看 `ci.yml` 的 `check` job 日志 | 维护者 + AI 协作会话 | 2026-12-10 | 待验证 |
| T3 | tag push 后 `release` job 的 `gh api` 轮询能否等到同 SHA 的 `ci.yml` 结论 | P2 | 首次 `vX.Y.Z` tag 实跑，比对 release / ci 运行时序与结论 | 发布者 @sutchan | 2026-12-20 | 待验证 |
| T5 | 三语键集合全量一致性（现有测试仅抽查部分字段） | P2 | 随 S13 增加 parity 测试并跑 `go test -tags nogui ./...` | 维护者 + AI 协作会话 | 2027-01-15 | 待验证 |
| T2 | `macos-latest` 上 darwin/amd64 交叉编译的 SDK 可用性 | P3 | 同 T1，核对两项 macOS 归档是否生成 | 维护者 + AI 协作会话 | 2026-12-10 | 待验证 |

> 待验证项与待办的区别：待办是明确的代码 / 文档改动任务；待验证项是对既有 CI 与运行环境事实的确认。验证失败即转化为新的 `S` 编号任务并按优先级排期。

## 已完成归档索引

| 批次 | ID | 归档位置 | 迁移时间 |
|------|----|----------|----------|
| v2.3.10 | C1–C8 | [tasks-done.md 第十三节](./tasks-done.md) | 2026-11-26 16:23 |
| v2.3.8 | D1–D5 | [tasks-done.md 第十二节](./tasks-done.md) | 2026-11-26 16:23 |
| v2.3.6 | A1–A10 | [tasks-done.md 第十一节](./tasks-done.md) | 历史归档 |
| v2.3.5 | RV1–RV3 | [tasks-done.md 第十节](./tasks-done.md) | 历史归档 |

> 已完成任务的完整描述、状态与上下文统一在 [tasks-done.md](./tasks-done.md) 维护；本文件只保留索引，避免双份正文漂移。
