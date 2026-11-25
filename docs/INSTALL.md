# Repo Manager 安装部署指南 v2.3.7

> 适用对象：从零把 Repo Manager 跑起来的使用者与 CI 维护者。产品行为契约见 [OPENSPEC.md](./OPENSPEC.md)；接口参数见 [CLI-REFERENCE.md](./CLI-REFERENCE.md)；操作步骤见 [USER-GUIDE.md](./USER-GUIDE.md)。

## 目录

- [1. 环境依赖清单](#1-环境依赖清单)
- [2. 获取与安装](#2-获取与安装)
- [3. 配置步骤](#3-配置步骤)
- [4. 多环境部署方案](#4-多环境部署方案)
- [5. 部署验证方法](#5-部署验证方法)
- [6. 升级与卸载](#6-升级与卸载)

---

## 1. 环境依赖清单

### 1.1 运行时（使用发布产物）

Repo Manager 是**单文件可执行程序**，运行时零依赖——无需额外运行库、JVM 或 .NET。功能 1 / 5 需要系统已安装 **Node.js（含 npm）**，因为其实现调用 `npm root -g` 与 `npm install -g`。

| 平台 | 运行时要求 | 说明 |
|------|------------|------|
| Windows 10/11（amd64） | 无 | `repo-manager-windows-amd64.exe` 单文件 |
| macOS 12+（amd64 / arm64） | 无 | darwin 产物，首次运行需在系统设置中放行 Gatekeeper |
| Linux（amd64，桌面） | 桌面环境的 X11 / Wayland | GUI 窗口依赖显示服务；CLI 亦可用 |

> 当前发布矩阵不含 `windows/arm64` 与 `linux/arm64`（托管 runner 缺交叉 / 多架构工具链，见 [MAINTAINERS.md §6](./MAINTAINERS.md)）；本地 `scripts/build.ps1` 可补齐。

### 1.2 源码构建（开发者）

| 依赖 | 版本要求 | 用途 |
|------|----------|------|
| Go | **1.25+**（`go.mod:3`） | 编译 |
| C 工具链 | gcc / clang（CGO_ENABLED=1） | 默认 GUI 标签经 Fyne + OpenGL，必须可链接 C |
| OpenGL / 窗口库开发头 | 见 1.3 | Linux 默认标签编译必需 |
| Node.js + npm | ≥ 18 | 参考界面工程门禁（`npm ci` / `lint` / `type-check`）与文档门禁 |
| PowerShell 5.1+ / bash | — | 构建脚本 `build.ps1` / `build.sh` |

### 1.3 Linux GUI 构建的 X11 / Wayland 开发库

`go-gl/glfw v3.4` 在 Linux 默认标签同时编译 X11 与 Wayland 后端，缺库会在 `wl_platform.h` 处失败（[TROUBLESHOOTING.md §2](./TROUBLESHOOTING.md)）：

```bash
sudo apt-get install -y \
  libgl1-mesa-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev \
  libwayland-dev libxkbcommon-dev
```

---

## 2. 获取与安装

### 2.1 方式一：下载 Release 归档（推荐）

1. 进入仓库 [Releases](https://github.com/sutchan/RepoManager/releases) 页面；
2. 选择对应平台的归档：`RepoManager-v<版本>-windows-x64.zip` / `-macos-x64.zip` / `-macos-arm64.zip` / `-linux-x64.tar.gz`（命名规则见 [MAINTAINERS.md §3](./MAINTAINERS.md)）；
3. 下载 `SHA256SUMS.txt` 校验完整性：

   ```bash
   # Linux / macOS
   sha256sum --check SHA256SUMS.txt --ignore-missing
   # Windows PowerShell
   Get-FileHash RepoManager-v*-windows-x64.zip -Algorithm SHA256
   ```

4. 解压到任意目录（如 `D:\Tools\RepoManager\`），归档内含可执行文件、`README.MD`、`LICENSE`、`config.example.ini`；
5. （可选）把该目录加入 `PATH`，以便直接运行 `repo-manager`。

### 2.2 方式二：源码构建

```bash
git clone https://github.com/sutchan/RepoManager.git
cd RepoManager
go build -trimpath -o bin/repo-manager ./cmd/repo-manager        # GUI（需 C 工具链）
go build -tags nogui -o bin/repo-manager-cli ./cmd/repo-manager  # 无 C 工具链的校验构建
```

批量交叉编译 6 目标（缺对应 C 编译器时**跳过并警告**，不产出坏文件）：

```powershell
./scripts/build.ps1   # Windows；输出 bin/repo-manager-<os>-<arch>[.exe]
./scripts/build.sh    # macOS / Linux
```

build 编号（10 位 `yymmddhhmm`）与版本号由脚本经 `-ldflags -X repo-manager/internal/buildinfo.BuildID=...` 注入，同一批次共用同一编号（[CONTRIBUTING.md §8.2](./CONTRIBUTING.md)）。

### 2.3 方式三：CI 构建（Actions）

- 推送 `vX.Y.Z` tag → `release.yml` 产出归档并发布 Release（[CONTRIBUTING.md §8.4](./CONTRIBUTING.md)）；
- 手动触发（`workflow_dispatch`）→ 只产出 `release-archives` 压缩包 artifact，不创建 Release，适合出测试包。

---

## 3. 配置步骤

配置文件 `config.ini` 位于**可执行文件同目录**（便携式），由程序自动创建与维护；手工编辑是安全的（程序保存时保留注释与未知行，见 [ARCHITECTURE.md §3.3](./ARCHITECTURE.md)）。模板见 [config.example.ini](./config.example.ini)。

| 键 | 类型与取值 | 默认 | 说明 |
|----|------------|------|------|
| `REPO_ROOT` | 目录路径 | Windows `E:\Github`，Unix `~/Github` | 项目根目录（多个项目的父目录）；CLI 位置参数可临时覆盖 |
| `LANG` | `zh` / `zh-Hant` / `en`（`tw`、`zh-TW` 归一化为 `zh-Hant`，非法值回落 `zh`） | `zh` | 界面语言；`-lang` 会持久化回写 |
| `RECURSIVE` | `true` / `1` / `yes` | 关闭 | 项目目录递归扫描（深度 3 层）；`-recursive` 与其「或」 |
| `EXTRA_CLEAN_DIRS` | 逗号分隔的单级目录名 | 空 | 并入功能 7 清理白名单；**仅接受安全目录名**，`..`、绝对路径、Windows 保留名等会被丢弃并告警（[ARCHITECTURE.md §6.1](./ARCHITECTURE.md)） |

两个随产物落地的数据文件（均自动生成，删除即重置）：

- `config.ini` — 上述 4 键；
- `repo-manager.log` — 运行日志（带时间戳追加；GUI / CLI 执行均写入）。

## 4. 多环境部署方案

| 环境 | 部署方式 | 验证要点 |
|------|----------|----------|
| Windows 桌面（日常使用） | Release zip 解压至固定目录，可创建快捷方式 | `-h` 正常；GUI 页脚显示 `vX.Y.Z · build ...` |
| macOS 桌面 | Release zip 解压；arm64 机型选 `-macos-arm64` | Gatekeeper 放行；GUI 正常渲染 |
| Linux 桌面 | Release tar.gz 解压；`chmod +x repo-manager` | 桌面会话中 GUI 打开；无显示服务时用 CLI |
| 便携 / U 盘 | 整目录拷贝（exe + config.ini 随行） | 换机后 `REPO_ROOT` 需按新机器调整 |
| CI / 自动化 | 源码构建 `-tags nogui` 校验 + `windows-smoke` 冒烟（[ci.yml](../.github/workflows/ci.yml)） | 门禁全绿；发布按 §8.4 闭环 |

> 没有服务器 / 容器部署形态：产品是本地桌面工具，无网络服务（[ARCHITECTURE.md §1.3](./ARCHITECTURE.md)）。

---

## 5. 部署验证方法

安装 / 构建完成后按顺序验证（**严禁在真实项目根目录做删除类验证**，先建含 1-2 个假项目 `package.json` 的临时目录）：

1. **版本自检**：`repo-manager -h`，应输出帮助与页脚 `v2.3.7 · build <10 位数字>`；
2. **语言持久化**：`repo-manager -lang en`，`config.ini` 出现 `LANG=en` 且无重复行；
3. **零副作用预演**：`repo-manager -dry-run 1 <临时目录>`，输出全部为 `[预演]` 且无任何文件改动；
4. **GUI 冒烟**：双击运行，按 [CONTRIBUTING.md §9.3](./CONTRIBUTING.md) 的 8 项清单逐项核对（语言切换、关于弹窗、功能 3 → 4 隐藏文件闭环等）；
5. **日志落盘**：确认 exe 同目录生成 `repo-manager.log`。

任一步失败 → 按 [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) 症状定位。

---

## 6. 升级与卸载

- **升级**：下载新版归档，覆盖可执行文件即可；`config.ini` 与 `repo-manager.log` 保留。跨大版本升级无迁移步骤（配置仅 4 键）。
- **卸载**：删除程序目录（含 `config.ini`、`repo-manager.log`）即可，程序不写注册表、不写系统目录（除功能 1 写入的用户环境变量 `NODE_PATH` 与各项目 `.npmrc`——卸载前如需还原请先执行一次功能 2）。
