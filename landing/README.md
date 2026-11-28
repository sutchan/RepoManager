# Repo Manager 落地页目录 (landing/) v2.3.16

> **定位**：产品 Web 展示页（`index.html` + `assets/`）的唯一落点；面向想快速了解并下载 Repo Manager 的开发者。
> 本文件记录目标受众、页面结构、视觉与交互要求、内容维护方式；功能事实见 [`docs/OPENSPEC.md`](../docs/OPENSPEC.md)，技术栈见 [`docs/CONTRIBUTING.md` §2.1](../docs/CONTRIBUTING.md)，版本历史见 [`CHANGELOG.MD`](../CHANGELOG.MD)。部署 / 上线方式仍未确定，属规划项。

---

## 1. 目标受众与核心信息

| 项 | 内容 |
|----|------|
| 目标受众 | 前端 / 全栈开发者、同时维护多个 Node.js 项目的团队、monorepo 维护者（与 [`OPENSPEC.md` §1](../docs/OPENSPEC.md) 目标用户一致）；单项目用户同样受益 |
| 用户问题 | 每个项目重复生成 `node_modules`、Next.js `.next` 与 Vite / Turbo / Angular 缓存堆积、全局依赖版本混乱 |
| 核心信息 | 7 项功能统一管理全局 `node_modules` 与项目缓存；危险操作（1 / 6 / 7）二次确认、预演模式、执行日志；Windows / macOS / Linux 单二进制；GUI + CLI 双用法；简中 / 繁中 / 英文三语 |
| 转化目标 | 下载 Release → 按三步开始上手 → 阅读文档；漏洞走私密渠道，不公开提 Issue |
| 事实边界 | 不虚构容量、耗时、用户量等数据；页面文案与 `docs/` 规格冲突时以规格为准并修正页面 |

---

## 2. 页面结构规划

`index.html` 为单页锚点导航，主要容器与区块使用语义化 kebab-case `id`，便于可访问性关联与测试定位。

| 顺序 / id | 板块 | 作用 | 事实来源 |
|-----------|------|------|----------|
| 1 `#landing-hero` | 项目简介 | 一句话价值、下载与文档 CTA、三项概览 | `README.MD` 产品介绍 |
| 2 `#landing-audience` | 目标受众 | 三类典型用户与收益 | `OPENSPEC.md` §1 / §3 |
| 3 `#landing-features` | 功能亮点 | 7 项功能编号与简述，对齐 CLI 参数 | `OPENSPEC.md` §3.0 |
| 4 `#landing-demo` | 演示 | GUI 操作说明 + 真实 CLI 调用示例（含 `-dry-run`） | `CLI-REFERENCE.md` §5 / `README.MD` 常用示例 |
| 5 `#landing-workflow` | 使用方式 | 取得程序 → 选择目录 → 执行功能 | `README.MD` 快速上手 |
| 6 `#landing-stack` | 技术栈 | Go / Fyne / build tag / 原生前端 | `CONTRIBUTING.md` §2.1 |
| 7 `#landing-safety` | 安全设计 | 二次确认、预演、扫描深度、日志 | `OPENSPEC.md` §4 / §7 |
| 8 `#landing-updates` | 更新记录 | 最近三条摘要 + 完整 CHANGELOG 链接 | `CHANGELOG.MD` |
| 9 `#landing-faq` | 常见问题 | 联网、平台、语言三项 | `README.MD` FAQ / `SUPPORT.md` |
| 10 `#landing-contact` | 联系方式 | Issue 模板 / 安全私密渠道 / 贡献指南 | `.github/SUPPORT.md` / `SECURITY.md` |
| 11 `#landing-footer` | 页脚 | 唯一版本与构建编号展示位 | 本文 §3.5 |

---

## 3. 内容约定

1. 入口为 `landing/index.html`，静态资源放 `landing/assets/` 并用相对路径引用，可直接静态预览；产物写 `landing/dist/` 或 `landing/out/`（已被 `.gitignore` 忽略）。
2. 文案事实与 [`OPENSPEC.md`](../docs/OPENSPEC.md) 一致，禁止虚构容量、耗时、用户量；产品名用 `Repo Manager`，二进制 `repo-manager`、平台仓库名 `sutchan/RepoManager` 按实况书写；本地 Git 远程以 `git remote -v` 为准（当前 `sutchan/RepoCleaner.git`），对照见 [`MAINTAINERS.md` §1](../docs/MAINTAINERS.md)。
3. 视觉对齐 [`DESIGN-SYSTEM.md`](../docs/DESIGN-SYSTEM.md) 与 [`tokens.json`](../docs/tokens.json)；新板块沿用既有令牌，不引入第二套配色。
4. 多语言按「简中 → 繁中 → 英文」优先级：词条集中在 `assets/dict.js`，切换 UI 覆盖全部语言，偏好写入 `localStorage`（枚举校验），`<html lang>` 同步 `zh` / `zh-Hant` / `en`；繁中使用台湾用语（設定 / 搜尋 / 檔案 / 複製 / 匯入）。
5. 页脚 `#landing-version` 展示 `vX.Y.Z · build yymmddhhmm`（10 位）；版本号取自 [`internal/buildinfo/version.go`](../internal/buildinfo/version.go) 的 `Version`，build 编号由 `npm run landing:version` 构建期写入 `assets/version.js`；**禁止硬编码、禁止运行时生成**，注入缺失时显示降级文案；页眉与其他区块不得出现版本号。
6. 源码文件单文件 ≤ 200 行（词条已拆到 `dict.js`，逻辑留在 `app.js`）；文档不拆分；本次改动文件头注释同步至当前版本，未改动文件不批量刷写。

---

## 4. 视觉与交互要求

| 维度 | 要求 | 当前实现 |
|------|------|----------|
| 响应式布局 | 移动端单列、桌面多列，内容在 320px 宽无横向溢出 | `grid-template-columns: repeat(auto-fit, minmax(...))`；`max-width: 640px` 断点收紧导航 / 区块间距；命令块 `overflow-x: auto` |
| 加载性能 | 首屏内容在 CSS 到达后即可读；无外部字体 / 图片 / CDN 请求 | 单 CSS、无 Web Font、内联系统字体栈；3 个小体积本地 JS（`version` / `dict` / `app`）均 `defer`，不阻塞解析；零构建依赖 |
| 可访问性 | 语义化 HTML、键盘导航、WCAG AA 对比度、动效可关闭 | `skip-link`、`aria-labelledby` 区块标题、导航 `aria-label`、语言按钮 `aria-pressed`、`:focus-visible` 焦点环、`prefers-reduced-motion` 降级；文字与底色沿用高对比令牌 |
| 动效 | 克制、≤ 300ms、仅状态反馈 | `--dur: 250ms` / `--easing` 令牌，按钮 hover 位移；`prefers-reduced-motion` 下关闭 |
| 深色模式 | 跟随系统偏好，亮色默认 | `@media (prefers-color-scheme: dark)` 覆盖令牌 |

---

## 5. 内容更新与维护方式

### 5.1 按变更类型同步

| 改动 | 操作 | 必做校验 |
|------|------|----------|
| 功能 / 安全 / FAQ 文案 | 改 `index.html` 默认简中 + `assets/dict.js` 三语对应键；事实回到 `OPENSPEC.md` / `SUPPORT.md` 核对 | 三语键集合一致；盲测切换无空文案 |
| 新增板块 | 在 `index.html` 加语义化 `id`、导航锚点与三语键；在本文 §2 登记板块与事实来源 | `check:asset-refs`、`check:docs` |
| 技术栈 / 版本展示 | 改文案后核对 `CONTRIBUTING.md` §2.1；版本只改权威源并跑 `npm run landing:version` | `node scripts/gen-landing-version.mjs` 输出 10 位 build |
| 样式 / 交互 | 改 `assets/styles.css`，令牌值同步回 `docs/tokens.json` 与 `DESIGN-SYSTEM.md` | 320 / 640 / 桌面宽度目视检查；键盘 Tab 与 reduced-motion 验证 |
| 发版 / bump | 按 [`CONTRIBUTING.md` §8.2](../docs/CONTRIBUTING.md) 同步全局展示位，并在 `CHANGELOG.MD` 新增小节、`assets/version.js` 重新生成 | 见本文 §5.2 |

### 5.2 门禁命令（仓库根执行）

```bash
npm run check:root-clean   # 根目录无散落条目
npm run check:docs         # Markdown 相对链接存在且大小写正确
npm run check:asset-refs   # landing/assets 资源引用检查（告警模式）
npm run check:i18n         # 参考界面三语键一致性（落地页三语须人工核对 dict.js）
npm run lint               # ESLint
npm run type-check         # tsc --noEmit
```

### 5.3 维护注意

- 落地页不写死任务清单；发现待办或「待验证 / TODO」类表述，先登记 [`docs/tasks.md`](../docs/tasks.md) 唯一任务来源，再回到页面只保留问题证据与指向任务清单的链接。
- `assets/version.js` 为生成文件，禁止手工修改；`dict.js` 三语键增删须同步三处，缺失词条回退简中。
- 部署平台与上线方式未定；确定后同步本文与 [`MAINTAINERS.md`](../docs/MAINTAINERS.md)，不得写成既成事实。
- 一次性脚本与中间产物放仓库根 `.tmp/`（已 gitignore），用完即删；不得散落 `landing/`。
