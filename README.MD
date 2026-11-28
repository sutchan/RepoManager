# 🧰 Repo Manager

<div align="center">

[![Platform](https://img.shields.io/badge/Platform-Win%20%7C%20macOS%20%7C%20Linux-blue?style=flat)](https://github.com/sutchan/RepoManager)
[![Go](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)
[![Version](https://img.shields.io/badge/Version-2.3.13-green?style=flat)](https://github.com/sutchan/RepoManager)
[![License](https://img.shields.io/badge/License-MIT-orange?style=flat)](https://github.com/sutchan/RepoManager)
[![i18n](https://img.shields.io/badge/i18n-简中%2F繁中%2FEnglish-9C27B0?style=flat)](https://github.com/sutchan/RepoManager)

</div>

> 一键让所有前端项目保持干净：无冗余 `node_modules`、无 `.next` 缓存、无打包垃圾。
> 适配 **Next.js / React / Vue / Vite / Angular / Svelte / 常规 Node 项目**；提供桌面 GUI 与命令行两种用法。

---

## 📑 目录

- [它能做什么](#-它能做什么)
- [安装与构建](#-安装与构建)
- [快速上手](#-快速上手)
- [功能详解](#-功能详解)
- [配置与命令行](#-配置与命令行)
- [常见问题](#-常见问题)
- [注意事项与更多资源](#-注意事项与更多资源)

---

## ✨ 它能做什么

| 痛点 | Repo Manager 的做法 |
|------|---------------------|
| 每个项目重复生成巨大的 `node_modules` | 启用全局 `node_modules`，项目不再生成本地依赖 |
| Next.js 的 `.next` 缓存占用空间 | 一键禁用 / 恢复缓存 |
| `dist` / `build` 等打包产物堆积 | 一键清理所有项目缓存 |
| Vite / Turborepo / Angular / Svelte 缓存膨胀 | 专项清理 `.vite` / `.turbo` / `.angular` / `.svelte-kit` 等目录 |
| 依赖版本混乱 | 统一安装常用全局依赖 |

- 🖥️ Go 编译的**单文件可执行程序**（Windows / macOS / Linux），零第三方运行时依赖
- 🎨 Fyne 桌面 GUI 一键操作；给出功能号 `1`–`7` 时进入非交互命令行模式
- 🌏 简中 / 繁中 / 英文三语界面与执行日志
- 🧪 预演模式（dry-run）：先看将删除什么，再决定是否执行
- 🌲 递归扫描（monorepo）：深度 3 层，默认关闭
- 🗂️ 可追加自定义清理目录（仅限项目目录内部）
- 🛡️ 危险操作（功能 1 / 6 / 7）前二次确认，避免误删
- 📝 执行日志写入可执行文件同目录的 `repo-manager.log`

> 删除操作不可逆。请先在隔离目录用 `-dry-run` 预演，确认后再执行。

---

## 📦 安装与构建

### 下载预编译版

从 [GitHub Releases](https://github.com/sutchan/RepoManager/releases) 下载对应平台的归档：

```
RepoManager-v<版本>-<os>-<arch>.zip      # Windows / macOS
RepoManager-v<版本>-<os>-<arch>.tar.gz  # Linux
```

解压得到 `repo-manager`（Windows 为 `repo-manager.exe`），用随附的 `SHA256SUMS.txt` 校验：

```bash
sha256sum -c SHA256SUMS.txt        # Linux / macOS
```

```powershell
Get-Content .\SHA256SUMS.txt | ForEach-Object { $h, $f = $_ -split '\s+', 2; if ((Get-FileHash $f -Algorithm SHA256).Hash.ToLower() -eq $h.ToLower()) { "OK $f" } else { "FAIL $f" } }
```

> 当前发布矩阵覆盖 windows-amd64、darwin-amd64 / arm64、linux-amd64；`windows/arm64` 与 `linux/arm64` 暂不产出，可用本地脚本补齐。

### 从源码构建

**前置要求**

| 用途 | 要求 |
|------|------|
| 运行功能 1 / 5（依赖管理） | Node.js ≥ 18 与 npm |
| 构建产品本体 | Go ≥ 1.25 |
| 构建完整 GUI | C 工具链 + 系统 OpenGL；Linux 还需 X11 / Wayland 开发库 |

```bash
# 完整 GUI（需 C 工具链与 OpenGL）
go build -o repo-manager ./cmd/repo-manager

# 无 C 工具链时：仅命令行功能号模式，GUI 不可用
go build -tags nogui -o repo-manager ./cmd/repo-manager

# 交叉编译全部 6 个目标到 bin/（Windows 用 ./scripts/build.ps1）
./scripts/build.sh
```

> 各平台的具体依赖与编译报错处理见[排障手册](./docs/TROUBLESHOOTING.md) §2。

---

## 🚀 快速上手

1. **运行**：不带参数启动 GUI（Windows 可双击）；带功能号 `1`–`7` 直接执行对应功能。
2. **选择仓库根目录**：GUI 输入框 / 目录选择按钮（可保存），或命令行位置路径（须与功能号同用）。
3. **选择功能**：GUI 点按钮，或 CLI 写功能号。单独给路径会启动 GUI 并忽略该路径。

| 场景 | 建议操作 |
|------|----------|
| 首次配置 | `[1]` 启用全局 node_modules → `[5]` 安装常用依赖 → `[3]` 禁用 `.next` 缓存 |
| 磁盘空间不足 | `[6]` 清理所有缓存 |
| Vite / Turborepo / Angular / Svelte | `[7]` 清理构建缓存 |
| 想先看影响 | 任意功能加 `-dry-run`，不修改任何文件 |

应用页脚展示 `vX.Y.Z · build yymmddhhmm`（10 位构建编号），用于核对版本与构建时间。

---

## 📖 功能详解

| 编号 | 功能 | 作用与注意 |
|------|------|-----------|
| 1 | 启用全局 node_modules | 获取 `npm root -g`、设置 `NODE_PATH`、为项目生成 `.npmrc` 并清理本地 `node_modules`。**先运行功能 5 安装依赖**，否则项目会找不到模块。需二次确认 |
| 2 | 重置 Node 配置 | 清除 `NODE_PATH` 与各项目 `.npmrc`，恢复 npm 默认行为 |
| 3 | 禁用 `.next` 缓存 | 删除 `.next` 目录并创建隐藏只读的 `.next` 文件阻止重建；会影响 Next.js 热更新 |
| 4 | 恢复 Next.js 缓存 | 删除该占位文件，恢复缓存与热更新 |
| 5 | 安装常用全局依赖 | `react react-dom next`、`vue @vue/cli`、`axios express`、`pnpm yarn`、`cross-env dotenv nodemon pm2 rimraf mkdirp`；逐组检查退出码并汇总失败 |
| 6 | 清理所有缓存 | `node_modules`、`.next`、`.nuxt`、`dist`、`build` 及 npm/yarn/pnpm 锁文件。需二次确认 |
| 7 | 清理 Vite / 构建缓存 | `.vite`、`.turbo`、`.nuxt`、`.angular`、`.svelte-kit`、`dist`、`build`、`coverage`、`.nyc_output` |

> 默认只扫描仓库根目录的**直接子目录**；monorepo 深层项目需开启递归（见下节）。

---

## ⚙️ 配置与命令行

### 仓库根目录从哪来

优先级：**CLI 位置路径 → `config.ini` 的 `REPO_ROOT` → 平台默认**（Windows `E:\Github`，macOS / Linux `~/Github`）。

`config.ini` 位于**可执行文件同目录**（不入库），可复制 [`docs/config.example.ini`](./docs/config.example.ini) 作为起点：

```ini
REPO_ROOT=E:\Github
LANG=zh
RECURSIVE=false
EXTRA_CLEAN_DIRS=
```

| 键 | 说明 |
|----|------|
| `REPO_ROOT` | 仓库根目录；GUI「保存」或 `-lang` 会写回本文件 |
| `LANG` | `zh`（简中）/ `zh-Hant`（繁中）/ `en`（英文）；`tw`、`zh-TW` 归繁中，非法值回落简中 |
| `RECURSIVE` | `true` / `1` / `yes` 开启递归扫描（深度 3 层）；CLI `-recursive` 优先 |
| `EXTRA_CLEAN_DIRS` | 逗号分隔的自定义目录**名**，并入功能 6 / 7 白名单；禁止路径分隔符与 `..`，仅匹配项目目录内部 |

### 参数速查

```bash
repo-manager [选项] [功能号] [目录路径]
```

| 参数 | 说明 |
|------|------|
| `-lang zh\|zh-Hant\|en` | 设置语言并写入 `config.ini` |
| `-dry-run` | 预演：只输出将执行的操作，不修改文件与环境变量 |
| `-recursive` | 递归扫描深层项目（深度 3 层） |
| `-h` / `-help` / `--help` | 帮助信息（固定英文，含版本与 build） |
| `1`–`7` | 直接运行对应功能，进入命令行模式 |
| 目录路径 | 指定仓库根目录，**仅在给出功能号时生效** |

无 `-root` / `-version` 标志，也没有交互式文本菜单；超出 `1`–`7` 的数字会被当作路径。

### 常用示例

```bash
repo-manager                                                       # 启动 GUI
repo-manager -dry-run 6 "D:\My Projects"                           # 预演清理缓存（含空格路径）
repo-manager -lang zh-Hant -dry-run -recursive 6 "D:\My Projects"  # 繁中 + 预演 + 递归
repo-manager 5 E:\Projects                                         # 为指定目录安装全局依赖
repo-manager -h                                                    # 查看版本与构建信息
```

首次删除前建议按「预演 → 确认执行 → 查看 `repo-manager.log`」三步走一遍。完整参数契约见[功能规格 §4.3](./docs/OPENSPEC.md)。

---

## ❓ 常见问题

**Q：功能 1 会影响现有项目吗？**
A：启用后项目不再生成本地 `node_modules`，统一从全局读取；请先运行功能 5，否则可能因找不到模块而报错。

**Q：功能 3 让热更新变慢？**
A：禁用缓存会略微影响热更新；开发完成后运行功能 4 恢复。

**Q：功能 6 / 7 删除后能恢复吗？**
A：删除不可逆。需要重新 `npm install`；建议确认不再开发或已备份时再执行。

**Q：为什么 `repo-manager E:\Github` 没切换目录？**
A：不带功能号会启动 GUI 并忽略位置路径。写成 `repo-manager 6 E:\Github`，或在 GUI 中选择目录后保存。

**Q：有交互式文本菜单吗？**
A：没有。只有 Fyne GUI 与「功能号 CLI」两种模式。

**Q：怎么知道会删哪些文件？**
A：加 `-dry-run` 或勾选 GUI「预演模式」；只输出计划操作，不修改文件。

**Q：日志在哪里？**
A：可执行文件同目录的 `repo-manager.log`，追加写入并带时间戳；公开提交前请脱敏。

更多症状与按主题检索的 FAQ 索引见[排障手册 §0](./docs/TROUBLESHOOTING.md#0-常见问题快速索引)。

---

## ⚠️ 注意事项与更多资源

1. 删除不可逆，执行前确认已备份；首次使用先用 `-dry-run`。
2. 默认只扫描直接子目录，深层项目需显式开启递归。
3. 功能 3 影响 Next.js 热更新，可在开发结束后用功能 4 恢复。
4. 功能 1 前先运行功能 5；全局依赖路径变更需重开终端。

**更多资源**

- [文档中心](./docs/README.md)：全部规范文档与按角色阅读路径（开发 / 维护者资料在这里）
- [功能规格](./docs/OPENSPEC.md)：7 项功能的完整行为契约
- [支持指南](./.github/SUPPORT.md) / [Issue 模板](https://github.com/sutchan/RepoManager/issues/new/choose)：提问与报缺陷
- [安全策略](./.github/SECURITY.md)：**发现漏洞请勿公开提 Issue**，按私密渠道报告
- [更新日志](./CHANGELOG.MD)：版本变更历史

---

<div align="center">

Repo Manager v2.3.13 · MIT License · [MIT 许可证](./LICENSE) · [文档中心](./docs/README.md)

</div>
