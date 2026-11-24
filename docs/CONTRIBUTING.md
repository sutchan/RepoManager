# Repo Manager 开发规范与贡献指南 v2.3.7

> 本文是 Repo Manager 仓库的**唯一开发规范来源**，规定「如何改代码、如何写文档、如何提交与评审」。
> 本文存放于 `docs/`；所有贡献者（含 AI 协作会话）在动手前必须阅读本文。
> 文档导航见 [docs/README.md](./README.md)；若本文与代码实现冲突，以代码为准并及时修正本文。

---

## 1. 文档体系与职责边界

仓库文档按「一个主题一份权威文档、统一归档 `docs/`」组织，禁止同一事实在多处重复维护。下表路径均为**仓库根相对路径**；完整文档地图与按角色阅读路径见 [`docs/README.md`](./README.md) §1 / §2，本文只定义职责边界。

### 1.1 适用范围与优先级

**规则优先级（高 → 低）**：用户明确指示 > 项目规范（本文 / `OPENSPEC.md` / `DESIGN-SYSTEM.md`）> 全局用户规则 > 工具默认行为。同级冲突以「更具体、更新近」为准；**代码与文档冲突时以代码实况为准**，并立即修正文档。

| 规范类别 | 权威位置 | 适用范围 | 约束强度 | 冲突裁决 |
|----------|----------|----------|----------|----------|
| 功能行为契约 | `OPENSPEC.md` §3 / §4 / §7 | 7 项功能、交互、兼容性、已知限制 | 强制 | 功能事实以规格为准；规格与代码冲突以代码为准 |
| 工程与协作规范 | 本文 §3 / §4 / §5 / §6 / §7 / §9 / §10 / §11 | 全部源码、脚本、文档、提交与审查 | 强制 | 流程冲突以本文为准 |
| UI 视觉规范 | `DESIGN-SYSTEM.md` + `tokens.json` | 产品 GUI、参考界面、落地页 | 强制（遗留代码按现状） | 视觉冲突以设计系统为准 |
| 版本与发布规范 | 本文 §8 + `MAINTAINERS.md` §3 | bump、CHANGELOG、Release | 强制 | 版本权威源为 `internal/buildinfo.Version` |
| 历史归档 | `CHANGELOG.MD` 旧版本段、`tasks-done.md` | 仅追溯，不作当前事实 | 不追溯改写 | 不得作为现行规范引用 |
| 生产就绪门槛 | [`PRODUCTION-READINESS.md`](./PRODUCTION-READINESS.md) | 代码 / 测试 / 安全 / 性能 / 日志 / 部署 / 版本 / 文档八类非功能要求 | 强制（含 `[待实现]` 过渡项） | 非功能指标以该文为准，功能冲突以 `OPENSPEC.md` 为准 |
| 本地记忆 | `.codebuddy/memory/`（见 §10.11） | AI 协作背景与决策记录 | 参考，非规范 | 不覆盖上述任何规范 |

**格式统一约定**：标题层级为 `# 文档名 vX.Y.Z` → `## N. 章节` → `### N.M 小节`；术语以 §5 为唯一来源；链接按实际大小写书写（§5）；页脚版本格式见 §8.2。

| 文档 | 职责（写什么） | 不写什么 |
|------|---------------|----------|
| `README.MD`（根） | 面向使用者：功能、快速开始、参数、FAQ、文档索引 | 不写开发规范、不写版本历史明细 |
| `docs/README.md` | 文档中心：文档地图、职责边界、按角色阅读路径、维护约定 | 不复述各规范正文 |
| `docs/OPENSPEC.md` | 产品功能规格：7 项功能的行为契约、交互、环境、兼容性、已知限制 | 不写代码风格、不写提交/分支规则、不重复变更日志 |
| `docs/CONTRIBUTING.md`（本文） | 开发规范：编码风格、命名、目录、提交、分支、审查、版本、术语、安全 | 不写功能介绍、不写设计令牌 |
| `docs/DESIGN-SYSTEM.md` | UI 设计规范：设计令牌、组件库、交互、响应式、a11y、差距清单 | 不写业务功能、不写开发流程 |
| `docs/tokens.json` | 机器可读设计令牌（与 `DESIGN-SYSTEM.md` §2 对齐） | 不作为产品版本号来源 |
| `prototype/README.md` | `prototype/` 目录说明：原型现状与目录用途 | 不复述设计规范正文、不存放代码 |
| `CHANGELOG.MD`（根） | **唯一版本历史来源**，按版本倒序记录变更 | 不重复功能规格 |
| `LICENSE`（根） | MIT 许可证正文 | — |
| `.github/CODE_OF_CONDUCT.md` | 社区行为标准与举报方式 | 不复述贡献流程 |
| `.github/SECURITY.md` | 安全策略：支持版本、私密漏洞报告渠道、删除类工具的安全边界 | 不写实现细节 |
| `.github/SUPPORT.md` | 提问前自查、有效 Issue 要素、不适用渠道 | 不复述功能文档 |
| `.github/FUNDING.yml` | 赞助配置（当前为全注释占位） | 不得填写未核实的账号 |
| `.github/ISSUE_TEMPLATE/` | `config.yml` + 缺陷报告 / 功能建议 Issue Form | — |
| `.github/CODEOWNERS` | 默认审查人 | — |
| `.github/PULL_REQUEST_TEMPLATE.md` | PR 审查清单模板 | — |

**引用规则**：需要版本历史 → `CHANGELOG.MD`；需要功能行为 → `docs/OPENSPEC.md`；需要协作流程 → 本文；需要视觉/组件规范 → `docs/DESIGN-SYSTEM.md`。禁止把 `CHANGELOG.MD` 的内容复制到 `docs/OPENSPEC.md` 或其它文档。

**文件名大小写**：根文档历史形成为大写扩展名 `README.MD` / `CHANGELOG.MD`；`docs/` 下全部为小写 `.md`；`docs/tokens.json` 为小写 `.json`。**所有链接必须按实际大小写与相对层级书写**（GitHub / Linux 区分大小写）：
- 从 `docs/` 内的文档引用根文件用 `../README.MD`、`../CHANGELOG.MD`、`../LICENSE`；
- 引用同目录文档用 `./OPENSPEC.md`、`./CONTRIBUTING.md`、`./DESIGN-SYSTEM.md`、`./tokens.json`；
- 禁止把链接目标写成小写扩展名的 `CHANGELOG.md` 或错误大小写的 `Openspec.md` 等不存在的路径。

**Community Health Files 归属**：GitHub 会自动识别仓库根、`docs/` 与 `.github/` 下的特定文件。根 `README.MD` / `CHANGELOG.MD` / `LICENSE` 与 `docs/CONTRIBUTING.md` 已分别承担 README、变更历史、许可证、贡献指南角色，因此 `.github/` **不再复制** `CONTRIBUTING.md` 正文（避免两处维护、且 `.github/` 优先级会盖住 `docs/` 的完整规范）。`.github/` 只放行为准则、安全策略、支持指南、赞助配置、Issue/PR 模板与 CODEOWNERS。

**禁止**：在 `.github/` 下创建 `readme.md`（与根 README 重复且易误导）；Community Health 文件与 PR 模板不受此限。

贡献前请阅读[行为准则](../.github/CODE_OF_CONDUCT.md)；报告安全漏洞走[安全策略](../.github/SECURITY.md)的私密渠道，不得公开提 Issue 或 PR；使用疑问先看[支持指南](../.github/SUPPORT.md)。

---

## 2. 环境与可用命令

| 项 | 要求 |
|----|------|
| 操作系统 | 产品支持 Windows / macOS / Linux；本仓库在三个平台上均可开发 |
| Go | **≥ 1.25**（`go.mod` 声明 1.25；产品本体 `repo-manager`） |
| GUI 框架 | Fyne **v2.8.2**（唯一 UI 依赖；默认标签构建需 CGO + OpenGL / X11 / Wayland） |
| Node.js | 文档最低 **≥ 18**；CI 固定 **20**（`ci.yml` / `release.yml`）；本机 `.nvmdrc` 为 **26.13.1**。贡献环境 ≥18 即可，CI 以 20 为准，本机标记仅供参考 |
| 参考界面 | Next.js 16.3.0 / React 19.2.8 / TypeScript 5（strict）/ Tailwind 4 / ESLint 9；已弃用，仅遗留预览 |

在**仓库根目录**执行（与 [`ci.yml`](../.github/workflows/ci.yml) 的 `check` job 一致；Windows 本机用对应命令观察输出）：

```bash
npm run check:root-clean                    # 根目录清洁度（CI 首步）
gofmt -l ./cmd ./internal                   # 输出应为空（CI 用 test -z "$(...)" 判定）
go vet -tags nogui ./...                    # 静态检查
go test -tags nogui -count=1 ./...          # 单元 / 临时目录集成测试
go test -tags nogui -race -count=1 ./...    # race 检测（需 CGO + gcc，无 gcc 时留待 CI）
go build -tags nogui ./...                  # 编译（nogui 排除 GUI 包）
go build -tags ci ./...                     # 编译 GUI 代码层（software driver，无需 C 工具链）
npm ci                                      # 安装参考界面依赖（CI）
npm run check:docs                          # 文档相对链接（含大小写）

go run ./cmd/repo-manager -h                # 命令行帮助
./scripts/build.ps1                         # 6 目标交叉编译（Windows / PowerShell）
./scripts/build.sh                          # 6 目标交叉编译（macOS / Linux）
```

> `npm run dev` / `type-check` / `lint` / `build` / `landing:version` 服务于参考界面与落地页，**不属于产品核心门禁**；产品质量门禁为上述 8 步，外加 §9.3 的 GUI 手动冒烟。

### 2.1 技术栈与依赖矩阵

| 层次 | 技术 / 版本 | 权威来源 | 状态 |
|------|-------------|----------|------|
| 产品本体 | Go 1.25 + Fyne v2.8.2（CGO / 系统 OpenGL） | `go.mod` | 当前 |
| 产品入口 | `cmd/repo-manager/`（main + gui/headless build tag） | `cmd/` | 当前 |
| 落地页 | 原生 HTML / CSS / JS + `landing/assets/version.js` | `landing/` | 当前 |
| 参考界面 | Next.js 16.3.0 / React 19.2.8 / TS 5 strict / Tailwind 4 / ESLint 9 | `package.json` / `tsconfig.json` | 已弃用，仅遗留预览 |
| 门禁脚本 | Node 18+（CI 20）+ `scripts/check-*.mjs` | `.github/workflows/ci.yml` | 当前 |
| 发布打包 | `scripts/package-release.sh`（zip / tar.gz + SHA256SUMS） | `.github/workflows/release.yml` | 当前 |

### 2.2 Build tag 与运行环境矩阵

| 标签 | 用途 | CGO / OpenGL | 产物 | 使用范围 |
|------|------|--------------|------|----------|
| 默认（无 tag） | 完整 Fyne GUI | 需要 C 工具链 + X11 / Wayland 开发库 | 可执行文件 | 发布构建、`gui-build` job |
| `nogui` | 排除 `internal/ui`，验证业务与平台层 | 不需要 | 无 GUI | 本地无 C 工具链、CI 静态门禁 |
| `ci` | Fyne software driver 编译 GUI 代码层 | 不需要 | 无（仅验证 API） | CI `check` / `windows-smoke` |
| `x11` | 仅 X11 后端，绕过 Wayland 编译 | 需要 X11 开发库 | 不支持 Wayland，**不用于发布** | 本地排障（`TROUBLESHOOTING.md` §2） |

---

## 3. 目录组织

### 3.1 根目录布局

```
RepoManager/
├── .github/                     # Community Health Files（禁止创建 .github/readme.md）
│   ├── workflows/               # ci.yml（门禁）+ release.yml（version / build / package / release，见 §8.4）
│   ├── CODE_OF_CONDUCT.md       # 行为准则
│   ├── SECURITY.md              # 安全策略与私密漏洞报告
│   ├── SUPPORT.md               # 支持指南
│   ├── FUNDING.yml              # 赞助配置（全注释占位）
│   ├── CODEOWNERS               # 默认审查人
│   ├── PULL_REQUEST_TEMPLATE.md # PR 审查模板
│   └── ISSUE_TEMPLATE/
│       ├── config.yml           # 模板配置与联系入口
│       ├── bug_report.yml       # 缺陷报告表单
│       └── feature_request.yml  # 功能建议表单
├── docs/                        # 所有规范与规格文档归档
│   ├── README.md                # 文档中心（文档地图、职责边界、阅读路径）
│   ├── OPENSPEC.md              # 产品功能规格
│   ├── CONTRIBUTING.md          # 本文：开发规范
│   ├── DESIGN-SYSTEM.md         # UI 设计规范
│   ├── TROUBLESHOOTING.md       # 使用与构建排障
│   ├── MAINTAINERS.md           # 维护者手册（命名 / 仓库设置 / 发布 / 文档治理）
│   ├── PRODUCTION-READINESS.md  # 八类非功能要求与发布前检查
│   ├── tokens.json              # 设计令牌
│   ├── config.example.ini       # 运行配置示例模板（4 键）
│   └── tasks.md / tasks-done.md # 任务单一来源与完成归档
├── public/                      # 参考界面应用静态资源（SVG）
├── landing/                     # Web 展示项目落地页（index.html + assets/；README.md 为目录说明）
│   ├── README.md
│   └── assets/                  # 样式 / 三语脚本 / 构建期注入的 version.js
├── scripts/                     # 工程脚本
│   ├── check-doc-links.mjs      # 文档相对链接校验（npm run check:docs）
│   ├── check-root-clean.mjs    # 根目录清洁度门禁（npm run check:root-clean）
│   ├── gen-landing-version.mjs # 落地页版本与 build 编号构建期注入（npm run landing:version）
│   ├── package-release.sh      # Release 资产归档与 SHA256SUMS 生成
│   ├── build.ps1                # Go 交叉编译（Windows / PowerShell）
│   └── build.sh                 # Go 交叉编译（macOS / Linux）
├── .tmp/                        # 会话级临时文件唯一落点（内容不入库，用完即清）
├── reports/                     # 工具正式产物目录（报告 / 覆盖率 / 测试结果，内容不入库）
├── go.mod / go.sum              # Go 模块定义与依赖锁定（Fyne）
├── cmd/repo-manager/            # 产品本体入口：main.go + gui.go（默认）/ headless.go（nogui）
├── internal/                    # Go 源码：buildinfo / config / engine / i18n / logs / platform / ui（7 个包）
├── src/                         # 已弃用参考界面应用（Next.js App Router，v1.x 遗留）
│   ├── app/                     # 路由与全局样式：layout.tsx / page.tsx / globals.css / favicon.ico
│   ├── components/
│   │   ├── icons.tsx            # 内联 SVG 图标（零外部图标依赖）
│   │   └── ui/                  # 无状态基础组件：alert/badge/button/card/code-block/table/tabs
│   └── lib/
│       ├── i18n.ts              # 唯一 UI 文案字典（zh / zh-Hant / en）
│       ├── utils.ts             # 纯工具函数（cn 等）
│       └── data/                # 单一数据源：features/project-types/steps/repos + index.ts barrel
├── prototype/                   # 原型目录，仅保留 README.md 说明，无代码、无规范正文
│   └── README.md
├── RepoCleaner.bat              # v1.x 旧版 Windows 脚本（已弃用，保留兼容）
├── config.ini                   # 运行配置：4 键（本机生成，v2.2.0 起不入库，示例见 docs/config.example.ini）
├── lang.ini                     # 语言文案参考节 [en] / [zh]（v1.x 历史保留，Go 不读取）
├── README.MD / CHANGELOG.MD / LICENSE     # 入口、变更历史、许可证
├── RepoManager.code-workspace   # VS Code 工作区配置，禁止删除（命名现状见 MAINTAINERS.md §1）
├── .nvmdrc                      # 本机 Node 版本标记
├── package.json / package-lock.json / tsconfig.json / next.config.ts
├── postcss.config.mjs / eslint.config.mjs / next-env.d.ts
└── .gitignore
```

### 3.2 组织规则

1. **所有规范 / 规格文档放 `docs/`**；根只保留入口 `README.MD`、`CHANGELOG.MD`、`LICENSE`。**产品代码一律放根 `cmd/` 与 `internal/`（Go）**；`src/` 为已弃用参考界面，禁止新增功能代码；`prototype/` 只放原型说明，禁止再放回可运行代码或规范正文；**Web 展示项目的落地页面放 `landing/`**，其内容约定见 [`landing/README.md`](../landing/README.md)。
2. **单一数据源**：产品 UI 与执行日志文案只放 `internal/i18n`（界面用 `T()`、日志用 `Logs()`）；`engine` 的业务常量（清理目录清单、依赖包列表、`.npmrc` 模板）集中在其包内定义，禁止散落在 `ui` 层。已弃用的 `src/lib/` 不再新增数据；落地页文案须与 `OPENSPEC.md` 事实一致，禁止虚构数据。
3. **Go 依赖**：只经 `go.mod` 管理（`go get` / `go mod tidy`），提交 `go.sum`；禁止 vendor 目录与多套依赖清单并存。UI 依赖仅限 `fyne.io/fyne/v2`，不得再引入第二套 GUI 框架。
4. **新增模块按职责落位**：入口 → `cmd/repo-manager/`；OS 差异 → `internal/platform/`；业务功能 → `internal/engine/`；配置读写 → `internal/config/`；文案 → `internal/i18n/`；界面 → `internal/ui/`；落地页 → `landing/`；规范 → `docs/`；工程脚本 → `scripts/`。
5. **文件行数**：所有源代码文件（`.ts/.tsx/.js/.jsx/.vue/.py/.dart/.go/.java` 等）单文件 **≤ 200 行**。超出时按职责拆分为更小模块（组件拆子组件、逻辑抽纯函数/常量/类型），且**不得改变对外导出契约**。
   - **文档文件（`.md/.markdown/.mdx/.rst`）不拆分**，即使超过 200 行也保持单文件完整连贯。
   - `.bat` 为批处理脚本，不纳入源码 200 行规则，但仍须按标签分节、保持可读。
6. **临时文件统一落固定子目录**：一次性脚本、探针、调试日志、中间产物一律写入 `.tmp/`；工具正式产物（审计报告、覆盖率、`*-result.json` 等）写入 `reports/` / `coverage/` / `test-results/`。写入前先 `mkdir -p`，本轮产生的 `.tmp/` 内容在任务结束前清理；需长期留存的证据改写进 `docs/` 或测试断言。目录忽略规则以 `.gitignore` 实况为准，并在 ESLint / `tsconfig.json` / `check-doc-links.mjs` 中排除，避免临时脚本让门禁假红。**禁止**在仓库根或顶层散落 `.tmp-*.mjs`、`tmp-*.js`、`*.log`、`audit-*.json` 等；根目录白名单由 `scripts/check-root-clean.mjs` 强制校验（CI 首步），新增根目录工程文件须同时核对 `.gitignore` 与白名单。产品运行时的 `config.ini` / `repo-manager.log` 属 exe 同目录便携配置，不在此约束内。
7. **新增 / 删除 / 重命名文档**须同步三处：[`docs/README.md`](./README.md) §1 文档地图与 §3 目录树、本文 §3.1、[`MAINTAINERS.md`](./MAINTAINERS.md) §4 变更对照表；根 `README.MD` 面向使用者与构建者，只保留少量入口链接，**不再维护文档索引与目录树**。新增文档标题带当前版本号，并在 `CHANGELOG.MD` 对应版本的 `Docs` 分组记录。**新增根目录子目录**（如 `landing/`）还须登记 `scripts/check-root-clean.mjs` 的 `ALLOWED_DIRS` 白名单并核对 `.gitignore`。同一事实只在一处维护，其余用链接引用。
8. **命名层次**：产品名 `Repo Manager`、Go 模块与二进制 `repo-manager`、Git 远程 `sutchan/RepoManager`；本机目录名 `RepoManager` 与工作区文件 `RepoManager.code-workspace` 与产品名对应，v2.3.3 已完成技术标识同步。对照与决策见 [`MAINTAINERS.md` §1](./MAINTAINERS.md)；当前无未完成项，见 [`tasks.md`](./tasks.md)。

---

## 4. 编码风格

### 4.1 通用

- **编码**：所有文本文件使用 **UTF-8**；`RepoCleaner.bat` 通过 `chcp 65001` 将控制台切到 UTF-8（不写「GBK/UTF-8」这类模棱两可的表述）。
- **缩进**：Go 用 `gofmt` 默认（Tab）；TypeScript / JS / CSS / JSON 用 2 空格；`RepoCleaner.bat` 用 4 空格。
- **引号**：TS/JS 统一双引号；JSX 属性用双引号；CSS 用双引号（字体名等）。
- **分号**：TS/JS 语句必须写分号。
- **注释**：关键逻辑、非直觉判断（尤其批处理的延迟展开与安全边界）必须写**中文注释**；禁止用注释复述显而易见的代码。
- **文件头注释**：源代码文件首行标注 `// <相对路径> vX.Y.Z`；CSS 用 `/* ... */`；`.bat` 用 `::`；`.ini` 用 `;`。示例（`vX.Y.Z` 为占位，实际写该文件本次改动后的产品版本号）：
  ```ts
  // src/components/ui/button.tsx vX.Y.Z
  ```
  头注释中的版本号是**该文件自身的版本**，随「本次实际修改该文件」更新；**未改动的文件不得批量刷新头注释**，因此仓库中不同文件的头注释版本并存属正常现象。

### 4.2 TypeScript / React（已弃用 `src/` 参考界面，规则仅适用遗留代码）

- 启用 `strict`，**禁止 `any`**；需要宽松类型时用 `unknown` + 类型收窄。
- 组件用**函数组件 + 命名导出**；仅页面/布局用 `export default`。
- 需要 hooks、事件、`navigator` 等浏览器 API 的组件，文件首行必须是 `"use client"`（放在 import 之前）。App Router 下服务端组件使用 `useState` 会导致构建失败。
- 组件 props 继承原生属性（如 `React.ButtonHTMLAttributes<HTMLButtonElement>`）并透传 `...props`，保证 `id`、`aria-*` 可传入。
- 样式用 `cn()`（`src/lib/utils.ts`）合并 className，禁止字符串拼接。
- 状态尽量局部化；无状态组件不引入 hooks。
- `Promise` 异常必须捕获处理（如 `CodeBlock` 的 clipboard 写入用 `try/catch`）。

### 4.3 样式与设计令牌

- 颜色/字体/圆角/动效**一律走 `src/app/globals.css` 的 `@theme` 令牌生成的工具类**：`bg-background`、`text-foreground`、`text-muted-foreground`、`bg-card`、`border-border`、`text-primary` 等。
- **禁止散落硬编码颜色**（`bg-[#0d1117]`、`text-gray-400` 等任意值）；新令牌先在 `docs/tokens.json` 与 `docs/DESIGN-SYSTEM.md` 登记，再写入 `globals.css`。
- 不引入 `next/font/google` 等**构建期联网**资源：离线环境会导致 `next build` 失败。使用 `font-sans` 系统字体栈。
- 图标使用 `src/components/icons.tsx` 内联 SVG，**不引入第三方图标库**。
- 当前仅实现**暗色主题**（`color-scheme: dark`）；浅色值属规划，不得在代码中假定已实现。

### 4.4 语义化 id

为页面级容器、主要区块、关键交互元素添加语义化 `id`，便于调试定位、ARIA 关联与端到端测试选择器。

- 命名用 **kebab-case**，表达用途：`app-header`、`overview`、`main-tabs`、`security-note`、`app-footer`。
- 列表项用稳定前缀 + 业务键：`repo-card-<repo>`、`tweak-card-<id>`、`btn-clean-<repo>`。
- 禁止 `div1`、`box`、`a` 等无意义命名；主要容器 `id` 全仓库唯一。
- 新增页面/组件时同步补齐；重构超 200 行文件时一并检查。

### 4.5 Windows 批处理（`RepoCleaner.bat`，已弃用）

- 保持 `setlocal enabledelayedexpansion`；块内（`( )` 中）读取块内赋值的变量必须用 `!VAR!` 而非 `%VAR%`。
- 所有路径变量**必须加双引号**，支持含空格路径。
- 删除目标只能来自 `for /d` 枚举出的目录变量，禁止把用户输入直接拼接进 `rd` / `del`。
- 功能 3 生成的 `.next` 文件带 `+h +s +r` 属性；功能 4 必须先 `attrib -h -s -r` 再 `del /f /q`（`del` 的多 `/a:` 只有末项生效，且默认不删隐藏/系统文件）。
- 危险操作（删除 `node_modules`）前必须有 `Y` 二次确认。
- 写 `config.ini` 用「读旧文件 → 过滤目标键 → 补写 → 回写」方式，**禁止整文件覆盖 `lang.ini`**（会清空 `[en]/[zh]` 文案节）。

### 4.6 Go（`repo-manager` 产品本体）

- 一律通过 `gofmt` 格式化、`go vet ./...` 无告警；导出符号必须有中文 doc 注释（`// FuncName ...`）。
- 错误处理：必须显式处理 `error`，禁止 `_ =` 吞错误；包装错误用 `fmt.Errorf("...: %w", err)`。
- 路径一律使用 `path/filepath`，禁止手工拼接 `/` 或 `\`；判定目录用 `os.Stat` + `IsDir`，禁止只靠「路径存在即目录」。
- 删除类操作只能针对 `platform` 枚举出的路径，禁止把用户输入直接拼接后删除。
- `ui` 层禁止直接调用 `os.Remove` 等系统调用，必须经 `engine` → `platform`。
- 单文件 ≤ 200 行（§3.2 规则 5 对 `.go` 同样生效）；超出按职责拆分，不得改变对外导出契约。

---

## 5. 命名约定

| 类别 | 规范 | 示例 |
|------|------|------|
| Go 源文件 / 包 | `snake_case.go`；包名小写单词；平台实现用 `build tag` 后缀 | `platform_windows.go`、`engine.go` |
| React 组件文件 | `PascalCase.tsx`（本仓库 `src/components/ui/` 沿用既有小写文件名，新增文件优先 PascalCase） | `CodeBlock.tsx` |
| 工具 / hooks / 普通模块 | `camelCase.ts` | `useLang.ts`、`utils.ts` |
| 数据模块 | 语义复数或主题名，`kebab-case` 可用 | `features.ts`、`project-types.ts` |
| 类型 / 接口 | `PascalCase`，不加 `I` 前缀 | `ScanResult`、`ButtonProps` |
| 常量 | `SCREAMING_SNAKE_CASE` | `APP_VERSION`、`DEFAULT_REPO_ROOT` |
| CSS 类名 / DOM id | `kebab-case` | `security-note`、`app-footer` |
| 文档文件 | 根目录既有 `README.MD` / `CHANGELOG.MD` 保持历史大写；`docs/` 下统一小写 `.md`，新增文档一律小写 | `docs/OPENSPEC.md`、`docs/README.md` |
| 分支 | 见 §7 | `fix/next-restore`、`docs/spec-sync` |
| 图标组件 | `PascalCase`，语义化 | `ScanLine`、`ShieldCheck` |

**术语统一**（全仓库文档与代码注释一致，禁止混用）：

| 统一写法 | 禁止写法 |
|----------|----------|
| `Next.js` | NextJS、Nextjs、next.js（句首除外） |
| `Nuxt` | NuxtJS |
| `node_modules`、`.next`、`dist`、`build` | 中文意译变体；代码块/行内代码用反引号 |
| 仓库根目录（`REPO_ROOT`，含多个项目的父目录） | 项目根目录（二者含义不同，禁止混用） |
| 项目根目录（单个项目，`REPO_ROOT` 的直接子目录） | — |
| 参考界面应用（根 `src/` 的 Next.js 应用，已弃用） | prototype 应用、前端工程 |
| 产品名 `Repo Manager` | `RepoCleaner` / `项目清洁助手`（旧名，仅 CHANGELOG 与归档文档的历史条目保留） |
| 产品本体（`repo-manager`，Go 编译的跨平台可执行文件） | Go 版、新版本 |
| 规范文档目录（`docs/`） | 原型目录（后者专指 `prototype/`） |
| 原型目录（`prototype/`，仅说明文件） | 设计规范目录（规范在 `docs/`） |
| 批处理脚本（`RepoCleaner.bat`，已弃用旧版） | bat 脚本、主程序 |
| 功能 1 ～ 功能 7（名称见 `OPENSPEC.md` §3.0） | 功能一、F1（仅代码变量名可用 F1_） |

---

## 6. 提交信息格式

采用 Conventional Commits 精简格式：

```
<type>: <description>

[可选正文]

[可选页脚]
```

### 6.1 type（必填，限以下取值）

| type | 用途 |
|------|------|
| `feat` | 新增功能 |
| `fix` | 修复缺陷 |
| `docs` | 文档更新 |
| `style` | 代码风格调整（不改行为） |
| `refactor` | 重构（不改行为） |
| `test` | 测试相关 |
| `chore` | 构建 / 依赖 / 配置变更 |
| `perf` | 性能优化 |
| `ci` | CI/CD 配置变更 |
| `revert` | 回滚提交 |

### 6.2 书写要求

- `description`：**≤ 50 字符，小写字母或中文开头、以动词开头、不加句号**，如 `fix: 修复功能4恢复 .next 缓存失效`。
- 正文（可选）：每行 ≤ 72 字符，解释「为什么」与影响范围，用 `-` 列表。
- 页脚（可选）：引用 issue / PR、标注 breaking changes；**提交信息正文或页脚须包含本次新版本号**（见 §8）。
- 一次提交只做一件事；禁止把无关改动混入同一提交。

示例：

```
docs: 建立 docs 文档体系并同步规范至 v1.2.5

- 新建 docs/ 归档全部规格与规范文档，补充文档中心索引
- 对齐 OPENSPEC 与实现，修正过时描述并统一术语
- 版本号改由 package.json 单一来源在构建时注入

v1.2.5
```

### 6.3 提交前检查

按 §2 的 8 步门禁命令逐项执行（**唯一命令清单**，与 `ci.yml` 对齐，勿在此维护第二套），并额外确认：

```bash
git status --porcelain             # 无 node_modules / .next / *.log / *.tsbuildinfo / bin/ / .tmp/ 被跟踪
git diff --check                   # 无空白错误
```

> 无 OpenGL / C 工具链时用 `-tags nogui` 跑静态门禁；`-tags ci` 验证 GUI 代码层。两者已列在 §2，`-tags x11` 仅用于本地排障、不用于发布。

---

## 7. 分支管理

### 7.1 分支模型

| 分支 | 类型 | 用途 |
|------|------|------|
| `main` | 常驻 | 生产/发布分支，始终保持可构建；**当前仓库实际仅启用此分支** |
| `dev` | 常驻（按需启用） | 开发集成分支，稳定后合并到 `main` |
| `feature/<短描述>` | 短期 | 新功能，从 `dev` 切出 |
| `fix/<短描述>` | 短期 | 非紧急缺陷修复 |
| `docs/<短描述>` | 短期 | 纯文档变更 |
| `chore/<短描述>` | 短期 | 依赖 / 配置 / 构建变更 |
| `hotfix/<短描述>` | 短期 | 生产紧急修复，从 `main` 切出，修完同时回合 `dev` |
| `release/<版本>` | 短期 | 发布准备（版本号、CHANGELOG 定稿） |

短描述用 `kebab-case`，如 `feature/dry-run-mode`、`fix/next-cache-restore`。

### 7.2 规则

1. **禁止直接向 `main` 强推**；禁止 `--force` 覆盖 `main` / `dev` 历史。
2. 短期分支合并后删除；保留提交历史，**不使用 squash 合并**，便于回滚。
3. 分支名与提交 type 保持一致语义。
4. 推送目标以 `git remote -v` 实况为准（当前 `https://github.com/sutchan/RepoCleaner.git`，远端名 `RepoCleaner` / `origin`）；平台仓库名为 `sutchan/RepoManager`（GitHub 改名后旧 URL 仍可访问）。两者并存属历史状态，详见 `MAINTAINERS.md` §1；禁止仅按文档硬编码远程地址。
5. 单人小改动（≤ 3 个文件、纯文档/拼写）可直接提交到 `main`；涉及 ≥ 5 个文件、新增依赖或架构调整，须走 PR 流程（§9）。

---

## 8. 版本号与变更记录

### 8.1 版本语义（SemVer）

| 级别 | 适用 | 升级方式 |
|------|------|----------|
| patch | 修复、文档、重构、样式、配置、依赖等任意修改 | 第三位 +1（`1.2.4 → 1.2.5`） |
| minor | 新功能、向后兼容变更 | 第二位 +1，第三位归 0 |
| major | 破坏性变更 | 第一位 +1，其余归 0 |

**任何文件修改都必须至少 bump 一次 patch**；minor/major 变更不再叠加 patch。

### 8.2 版本号单一来源

**唯一权威源：`internal/buildinfo/version.go` 的 `Version` 常量**（Go 产品本体）。`package.json` 的 `version` 保持同步（服务已弃用参考界面与文档引用），禁止在任何页面 / 组件中硬编码版本号。

每次 bump 需同步以下**全局展示位**（各一处）：

- `internal/buildinfo/version.go` → `Version` 常量（产品唯一权威源）
- `package.json` / `package-lock.json` → `version`（顶层与 `packages[""]` 两处）
- `README.MD`（根） → 版本徽章 `Version-X.Y.Z` 与页脚版本号
- `CHANGELOG.MD`（根） → 新增对应版本小节（置于最上方）
- `docs/tokens.json` → `version`
- `docs/DESIGN-SYSTEM.md`、`docs/README.md`、`docs/OPENSPEC.md`、`docs/CONTRIBUTING.md`、`docs/TROUBLESHOOTING.md`、`docs/MAINTAINERS.md`、`prototype/README.md` → 标题版本
- `docs/tasks.md` → 头部「当前版本」说明
- `.github/` 下新增的 Community Health 文件 → 标题或首行注释版本
- 被改动**代码文件**的头注释版本号

**禁止**：全仓库批量刷写所有源码文件的头注释；未改动文件头注释保持不变（头注释是文件自版本，不等于产品版本，看到旧版本与新版本并存属正常，不是版本漂移）。

**页脚版本与 build 编号（强制）**：应用页面的**页脚区**是版本号与 build 编号的唯一展示位置；页眉 / 顶部栏 / 其他区块一律不得出现版本号，页脚若原本展示「更新日期」一律改为 build 编号。

- **统一格式（全项目唯一）**：`vX.Y.Z · build yymmddhhmm`（示例 `v2.3.4 · build 2611221000`）——小写 `v`、`·` 两侧各一个空格、10 位构建时间（2 位年份 + 月 + 日 + 时 + 分，24 小时制）。禁止 `V2.3.4· build ...`（大写 V / 缺空格）与 12 位 `yyyymmddhhmm`；该格式对 Go GUI 页脚、CLI `-h`、参考界面页脚与落地页页脚一律适用。
- **Go 产品本体**：`internal/buildinfo/version.go` 的 `Version` 与 `BuildID`，由 `Footer()` 生成页脚文本；`internal/ui/ui.go` 用 `container.NewBorder` 的 bottom 区渲染（`NewBorder` 底部即页脚，页脚容器在 Web 侧须带 `id="app-footer"`）。`scripts/build.ps1`（`(Get-Date).ToString("yyMMddHHmm")`）/ `scripts/build.sh`（`date +%y%m%d%H%M`）在构建期以 `-ldflags "-X repo-manager/internal/buildinfo.BuildID=<10位>"` 注入；未注入时从 `debug.ReadBuildInfo` 的 `vcs.time` 推导本地时间（Go 布局 `0601021504`），仍缺失则回退为仅版本号。CLI `-h` 的 `BuildInfo()` 输出不算页面页眉。
- **已弃用 Web 参考界面**：`next.config.ts` 从 `package.json` 读取版本并注入 `NEXT_PUBLIC_APP_VERSION` 与 10 位 `NEXT_PUBLIC_BUILD_ID`，`src/app/page.tsx` 的 `<footer id="app-footer">` 展示，注入缺失时回退仅显示版本号。
- **落地页**：`scripts/gen-landing-version.mjs` 从版本权威源读取 `Version` 并生成 10 位 build 编号，写入 `landing/assets/version.js` 的 `label`（`npm run landing:version`），页脚渲染该 label；生成文件禁止手工修改。
- **禁止**：在页面 / 组件硬编码版本号或 build 编号；禁止运行时用 `new Date()` / `time.Now()` 现场生成（必须构建期注入，保证同一构建内稳定）。
- **验证**：`go test ./internal/buildinfo`（`Footer` 格式、回退、10 位布局）+ `npm run build` + `node scripts/gen-landing-version.mjs`（各页脚均输出 `vX.Y.Z · build <10位数字>`）。

### 8.3 CHANGELOG 写法

- 按版本倒序，标题 `## [X.Y.Z] - YYYY-MM-DD`。
- 分组用 `Added` / `Changed` / `Fixed` / `Docs` / `Security`，条目以 `- **要点** - 说明` 书写。
- 历史版本号与历史条目的原始措辞是**历史记录，不得改写**；批量替换版本文本或统一术语时必须排除 `CHANGELOG.MD` 的历史小节（其中仍保留 `NextJS` 等旧写法属正常，不算术语违规）。

### 8.4 发布流程与 Release 自动化

**规则（目标）**：GitHub Actions **所有工作流全部通过**才算完整发布闭环；禁止在门禁未全绿时手动上传或跳过 CI 发版。

**当前实现边界（workflow 实况，v2.3.3）**：`release.yml` 的 `release` job 以 `needs: [version, build, package]` 阻断，并在发布前通过 `gh api` 校验同 SHA 的其它工作流（含 `ci.yml` 与今后新增的工作流）全部 `completed|success`；运行中轮询等待，失败或超时即退出。`release` job 单独提权 `contents: write` / `actions: read`，顶层与其余 job 仅 `contents: read`。

| 项 | 约定 |
|----|------|
| 触发 | push 符合 `vX.Y.Z` 的 tag（必须与 §8.2 权威源 `Version` 一致）；两条工作流另支持 `workflow_dispatch` 与 `concurrency` 防重复，手动触发只产出压缩包 artifact、不创建 Release |
| 全绿门槛（v2.3.3） | `release.yml` 的 `version` + `build` + `package` 全部成功；门禁由 `ci.yml` 在同 SHA 运行，`release` job 自动等待并校验其结论，缺失或失败即阻断 |
| 发布动作 | `version` job 从 `internal/buildinfo/version.go` 读取版本；`build` job 按 4 目标矩阵产出二进制；`package` job 用 `scripts/package-release.sh` 打包为 `RepoManager-v<版本>-<os>-<arch>.zip / .tar.gz` 归档并生成 `SHA256SUMS.txt`；`release` job 用该 tag 创建 Release，说明取 `CHANGELOG.MD` 对应小节，附件只含归档与校验文件（build 号按 §8.2 在构建期 `-ldflags` 注入） |
| 一致性校验 | 发布前校验 tag == `internal/buildinfo/version.go` 的 `Version` == `package.json`；不一致即失败，不得发布 |
| 失败处理 | 发布工作流自身任一步失败则不发布；修复后删除并重新打同一 tag 重跑，禁止在失败状态下手工发版 |

**当前状态：已落地（v2.3.3）**——实现为两个工作流：

- [`.github/workflows/ci.yml`](../.github/workflows/ci.yml)：`push`（`main` 与 `v*` tag）、`pull_request` 与 `workflow_dispatch` 触发；`check` job 执行 §9.1 全部仓库级门禁（nogui，含 `go test -tags nogui -race`，ubuntu runner 具备 CGO/gcc），`windows-smoke` job 在 Windows runner 补跑 `go test -tags nogui` 与 `go build -tags ci`，`gui-build` job 在 ubuntu 安装 OpenGL / X11 / Wayland 开发库后以**默认构建标签**完整编译 Fyne GUI。`go-gl/glfw v3.4` 在 Linux 默认双后端，缺 `libwayland-dev` / `libxkbcommon-dev` 会在 `wl_platform.h` 处编译失败，排障见 [TROUBLESHOOTING.md §2](./TROUBLESHOOTING.md)。
- [`.github/workflows/release.yml`](../.github/workflows/release.yml)：tag push 或 `workflow_dispatch` 触发；`version` job 读取版本权威源，`build` job 按 4 目标矩阵产出二进制，`package` job 调用 `scripts/package-release.sh` 打包归档并把 `release-archives` 作为 artifact 保留，`release` job（仅 tag push）校验同 SHA 其它工作流全绿、tag == 版本权威源 == `package.json`，从 `CHANGELOG.MD` 提取对应小节作 Release 说明；`gh release create` 只上传 `*.zip` / `*.tar.gz` 归档与校验文件，预发布以 `--prerelease` 标记区分。门禁不在本工作流重复执行。
- **全绿语义**：`release` job 以 `needs: [version, build, package]` 保证发布工作流自身全绿，并通过 `gh api` 校验同 SHA 其它工作流全部 `completed|success`；新增工作流会被该轮询自动纳入阻断范围。
- **已知缺口**：GitHub 托管 runner 无 `aarch64-w64-mingw32-gcc`，不产出 `windows/arm64`；`linux/arm64` 缺 arm64 架构的 X11 / Wayland / OpenGL 开发库，为避免该目标失败拖垮整个 `build` job，v2.3.3 起从发布矩阵移除（本地 `scripts/build.ps1` 仍可补齐）；`SHA256SUMS.txt` 仅覆盖实际产出的归档。

---

## 9. 代码审查流程

### 9.1 流程

1. 先用 [Issue 模板](https://github.com/sutchan/RepoManager/issues/new/choose) 描述问题或建议并确认未被记录；从目标分支切短期分支（§7.1）。
2. 完成改动，按 §6.3 跑完全部门禁。
3. `git add` → `git commit`（格式见 §6）→ `git push origin <分支>`。
4. 在 GitHub 创建 PR 到 `main`（或 `dev`），**PR 标题格式同提交信息**，描述按 `.github/PULL_REQUEST_TEMPLATE.md` 填写。
5. 审查人按 §9.2 清单逐项核对；有异议在对应代码行评论，修改后重新推送。
6. 全部通过后以 **merge commit** 合并（不 squash），删除短期分支。

### 9.2 审查清单

- [ ] 功能行为与 `docs/OPENSPEC.md` 对应条目一致；变更功能时已同步规格文档
- [ ] 新增/修改的 UI 对齐 `docs/DESIGN-SYSTEM.md` 与 `docs/tokens.json`，无硬编码颜色
- [ ] 主要容器与关键元素已加语义化 `id`（§4.4）
- [ ] 无 `fmt.Println` 调试输出、无吞错（`_ = err`）、无未使用的 import / 变量；Go 代码已 `gofmt` 且导出符号有 doc 注释
- [ ] §2 的 8 步门禁全部通过：`gofmt -l ./cmd ./internal` 无输出；`go vet/test/build` 用 `-tags nogui`、`go build -tags ci` 验证 GUI 层；`npm ci`、`check:docs`、`check:root-clean` 通过
- [ ] 单文件 ≤ 200 行（代码文件）；超出已按职责拆分且不破坏导出契约
- [ ] 平台层改动：路径经 `path/filepath` 处理、删除目标仅来自枚举、危险操作有二次确认
- [ ] 提交信息格式合规，且含新版本号
- [ ] 版本号已 bump，`package.json` / README 徽章 / CHANGELOG / 展示位一致
- [ ] `CHANGELOG.MD` 已新增对应小节；`README.MD` / `docs/OPENSPEC.md` 使用说明随功能变化更新
- [ ] 文档链接按实际大小写与相对层级书写，新增文档归入 `docs/`
- [ ] 新增 / 删除 / 重命名文档已同步 `docs/README.md` §1 / §3、本文 §3.1、`docs/MAINTAINERS.md` §4（根 `README.MD` 只保留使用者入口链接）；排障类变更同步 `docs/TROUBLESHOOTING.md`
- [ ] 未完成事项只记入 `docs/tasks.md`，完成后归档到 `docs/tasks-done.md`；无「规划写成既成事实」的表述
- [ ] `npm run check:docs` 与 `npm run check:root-clean` 均通过；GitHub Actions 各工作流为绿；发布类变更已按 §8.4 校验 tag 与版本权威源一致性
- [ ] 无密钥 / Token / 真实用户数据入库；`node_modules`、`.next`、`*.log`、`*.tsbuildinfo` 未被跟踪；`npm run check:root-clean` 通过，临时文件统一在 `.tmp/` / `reports/`
- [ ] 生产就绪：按 [`PRODUCTION-READINESS.md` §9](./PRODUCTION-READINESS.md) 逐项核对（覆盖率、依赖漏洞、性能指标、日志、回滚、版本展示位、文档登记）；带 `[待实现]` 的项已评估并在 PR 中透明记录

### 9.3 功能验证（GUI 手动冒烟，跨平台）

`repo-manager` 的 GUI 交互无自动化测试，发版前须在**隔离的临时目录**（放 1～2 个含 `package.json` 的假项目）手动验证：

| 检查项 | 预期 |
|--------|------|
| `repo-manager -h` | 输出帮助信息后退出 |
| `repo-manager -lang en` | `config.ini` 写入 `LANG=en` 且不产生重复行 |
| GUI 语言按钮 | 简中 / 繁中 / 英文三语切换生效，执行日志同步切换，语言偏好持久化 |
| GUI 关于按钮 | 弹窗展示功能说明，不出现版本号 |
| 功能 1（确认） | 生成 `.npmrc`、清理本地 `node_modules`；取消则返回且不删除 |
| 功能 3 → 功能 4 | 生成带隐藏属性的 `.next` 文件；功能 4 能成功删除该文件（无残留、无报错） |
| 功能 6 / 7 | 仅清理白名单目录；含空格路径可正确处理 |
| 路径不存在 | 弹窗报错并允许重新输入 |

> 严禁在真实仓库根目录做删除类验证。

---

## 10. 文档维护约定

1. 改功能 → 同步 `docs/OPENSPEC.md` + 根 `README.MD`；改 UI → 同步 `docs/DESIGN-SYSTEM.md` + `docs/tokens.json`；改流程 → 同步本文；改文档结构 → 同步 `docs/README.md`；维护者信息与发布校验 → 同步 `docs/MAINTAINERS.md`；症状处理 → 同步 `docs/TROUBLESHOOTING.md`。完整对照见 [`MAINTAINERS.md` §4](./MAINTAINERS.md)。
2. 文档中的**路径、命令、文件名、版本号必须与仓库实况一致**；写文档前先 `git log -1` 与 `node -p "require('./package.json').version"` 实查，禁止凭记忆填写；本文档与代码冲突时以代码为准。
3. 禁止在文档中保留「规划中但未实现」的内容而不加标注；要么删除，要么明确写「规划，当前未实现」，或记入 `docs/DESIGN-SYSTEM.md` §7 差距清单 / `docs/MAINTAINERS.md` §6 已知差距。
4. 数据型描述（数量、行数、命令数）须实算或引用唯一来源，禁止估算或虚构；历史归档（`CHANGELOG.MD` 旧版本段等）不追溯改写。
5. 文档文件不拆分；超长文档保持单文件完整。
6. 禁止创建 `.github/readme.md`；`*.code-workspace` 是工作区配置，禁止删除（当前文件名 `RepoManager.code-workspace`）。
7. 改动文档后必须执行 `npm run check:docs`（相对链接含大小写）与 `npm run check:root-clean`（根目录无散落条目），校验方法见 [`docs/README.md` §5](./README.md)。
8. **任务单一来源**：未完成 / 进行中事项只写 `docs/tasks.md`，完成后归档到 `docs/tasks-done.md`；其他文档不得另建待办清单。
9. **新增 / 删除 / 重命名文档**按 §3.2 规则 7 同步三处登记，并在 `CHANGELOG.MD` 对应版本的 `Docs` 分组记录。
10. 版本同步清单按 §8.2 执行，含 `docs/TROUBLESHOOTING.md` / `docs/MAINTAINERS.md` / `landing/README.md` 标题与 `docs/tasks.md` 头部版本说明。
11. **本地记忆文档**：`.codebuddy/memory/`（已被 `.gitignore` 忽略，不入库）记录 AI 协作的背景、核心决策与时间线，是**参考信息而非规范来源**，不得覆盖本文 / `OPENSPEC.md` / `DESIGN-SYSTEM.md`。记忆文档须按主题归类、标注时效性与可信度，并引用本文 §1 作为规范入口；本文只在 §1.1 表中提及它，**不建立指向未入库文件的可点击链接**（GitHub 上会 404）。

---

## 11. 安全红线

- 代码、配置、日志中不得包含密钥 / Token / 凭据。
- 用户输入必须校验：目录用 `os.Stat` + `IsDir` 判定；删除目标只取 `platform` 枚举出的路径，禁止拼接用户输入后删除。
- 前端不拼接未净化的 HTML（本仓库无 `dangerouslySetInnerHTML`）。
- 删除类操作必须二次确认，且不得在真实仓库根目录做验证。
- 发现安全漏洞按 [SECURITY.md](../.github/SECURITY.md) 私密报告，禁止公开披露可复现细节。
- 依赖定期审计：Go 用 `go list -m -u all`，Node 用 `npm outdated`；升级 `fyne.io/fyne/v2` 或 `next` 等核心依赖须跑完整门禁。
