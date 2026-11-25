# Repo Manager 接口文档（CLI 契约）v2.3.7

> **关于 OpenAPI**：OpenAPI 规范描述 HTTP API。Repo Manager 是本地桌面 / CLI 工具，**没有 HTTP 接口**（无网络监听、无 RPC，见 [ARCHITECTURE.md §1.3](./ARCHITECTURE.md)），因此**不适用 OpenAPI 规范**，本文以等效的 CLI 接口契约承载「路径 / 方法 / 参数 / 返回 / 错误码 / 调用示例」信息。若未来引入 HTTP API（如本地守护服务），须另立 OpenAPI 文档并登记 [tasks.md](./tasks.md)。
> 接口事实来源为 `cmd/repo-manager/main.go` 与 `internal/engine/options.go`；本文与代码冲突时以代码为准。

## 目录

- [1. 命令语法与参数](#1-命令语法与参数)
- [2. 退出码定义](#2-退出码定义)
- [3. 输出流与日志约定](#3-输出流与日志约定)
- [4. 配置文件契约](#4-配置文件契约)
- [5. 调用示例](#5-调用示例)

---

## 1. 命令语法与参数

```
repo-manager [flags] [功能号 1-7] [路径]
```

- 无功能号 → 启动 GUI；有功能号 → CLI 模式执行对应功能。
- 选项优先级：**命令行 > `config.ini` > 默认关闭**；`-recursive` 与配置值是「或」关系。
- 路径优先级：**命令行路径 > `config.ini` 的 `REPO_ROOT` > 平台默认**（Windows `E:\Github`，Unix `~/Github`）。
- 功能号 1-7 的名称与行为见 [ARCHITECTURE.md §1.2](./ARCHITECTURE.md) / [OPENSPEC.md §3](./OPENSPEC.md)；功能 1、6、7 为危险操作。

### 1.1 Flags

| Flag | 别名 | 取值 | 说明 |
|------|------|------|------|
| `-lang` | — | `zh`（默认）/ `zh-Hant` / `en`；`tw`、`zh-TW` 归一化为 `zh-Hant`，非法值回落 `zh` | 界面语言；**会持久化写回 `config.ini`** |
| `-dry-run` | — | 布尔 | 预演模式：功能 1-7 全部改为输出 `[预演] 将...`，零副作用，且跳过危险操作确认 |
| `-recursive` | — | 布尔 | 项目目录递归扫描（深度 3 层） |
| `-h` | `-help` | — | 输出帮助与版本页脚（`vX.Y.Z · build yymmddhhmm`）后退出 |

> 无 `-version` flag；版本信息仅出现在 `-h` 的 `buildinfo.BuildInfo()` 输出中（`main.go:132`）。

### 1.2 位置参数

| 位置 | 取值 | 说明 |
|------|------|------|
| 1 | `1`-`7` | 功能编号；不是 1-7 数字的首参不作为功能号（进 GUI） |
| 2 | 目录路径 | 本次执行的 `REPO_ROOT`，执行成功后写回 `config.ini` |

---

## 2. 退出码定义

仅 3 个退出码（`main.go`）：

| 退出码 | 场景 |
|:------:|------|
| `0` | 成功结束；**含用户在确认提示回答非 `Y`（打印「已取消」后正常返回）** |
| `1` | 三类失败：① 根目录校验失败（不存在 / 非目录 / 文件系统根）；② 非交互环境读确认输入失败（stdin EOF / 关闭，明确**不算用户取消**）；③ `engine.Run` 执行出错 |

脚本侧应把非 0 一律视为失败；「取消」与「成功」同为 0，需依赖输出文案区分（改进建议见 [ARCHITECTURE.md §7.2 S7](./ARCHITECTURE.md)）。

---

## 3. 输出流与日志约定

| 内容 | 流 | 证据 |
|------|----|------|
| 错误信息（`[错误] ...`） | stderr | `log.go:94` |
| 危险操作确认提示（`是否继续? [Y/N]: `） | stdout | `log.go:97` |
| 正常执行日志（含 `[预演]` / `[n/m]` 进度 / `[警告]`） | stdout | `engine.go:23-27` |
| 全部执行日志（带 `2006-01-02 15:04:05` 时间戳） | 同时追加写 exe 同目录 `repo-manager.log` | `main.go:114-116`、`logs.go` |

确认交互规则：`TrimSpace` 后必须 `EqualFold "Y"` 才执行；stdin 读取失败 → stderr 提示并退出 1（不静默当作取消）。

---

## 4. 配置文件契约

- 路径：可执行文件同目录 `config.ini`（便携式）；格式 INI，`键=值`，`#` 注释保留。
- 键恰为 4 个：`REPO_ROOT` / `LANG` / `RECURSIVE` / `EXTRA_CLEAN_DIRS`，语义见 [INSTALL.md §3](./INSTALL.md)。
- 写入为原子替换（临时文件 + `os.Rename`）；读取失败（非「文件不存在」）时 CLI 禁止回写并报错退出，防止损坏配置被覆盖。
- `RECURSIVE` 接受 `true` / `1` / `yes`；`EXTRA_CLEAN_DIRS` 逐项安全校验，不合格直接丢弃并告警。

---

## 5. 调用示例

```bash
# 帮助与版本
repo-manager -h

# 预演功能 1（零副作用，必先执行）
repo-manager -dry-run 1 E:\Github

# 危险操作：交互确认后执行（回答 Y 才继续）
repo-manager 1 E:\Github

# 脚本集成：非交互 + 显式路径（无 stdin 时危险操作退出码 1）
repo-manager 6 ~/Github < /dev/null; echo "exit=$?"

# 递归清理 Vite / 构建缓存（功能 7），英文界面
repo-manager -lang en -recursive 7 ~/Github

# 语言偏好持久化（写回 config.ini）
repo-manager -lang zh-Hant
```

```powershell
# Windows PowerShell 校验退出码
repo-manager -dry-run 6 E:\Github
if ($LASTEXITCODE -ne 0) { "预演失败" }
```

> 所有示例严禁指向真实项目根目录做删除验证；先在隔离临时目录演练（[INSTALL.md §5](./INSTALL.md)）。