# Repo Manager 接口文档（CLI 契约）v2.3.16

> **关于 OpenAPI**：OpenAPI 规范描述 HTTP API。Repo Manager 是本地桌面 / CLI 工具，**没有 HTTP 接口**（无网络监听、无 RPC，见 [ARCHITECTURE.md §1.3](./ARCHITECTURE.md)），因此**不适用 OpenAPI 规范**，本文以等效的 CLI 接口契约承载「路径 / 方法 / 参数 / 返回 / 错误码 / 调用示例」信息。若未来引入 HTTP API（如本地守护服务），须另立 OpenAPI 文档并登记 [tasks.md](./tasks.md)。
> 接口事实来源为 `cmd/repo-manager/main.go` 与 `internal/engine/options.go`；本文与代码冲突时以代码为准。

## 目录

- [1. 命令语法与参数](#1-命令语法与参数)
- [2. 退出码定义](#2-退出码定义)
- [3. 输出流与日志约定](#3-输出流与日志约定)
- [4. 配置文件契约](#4-配置文件契约)
- [5. 调用示例](#5-调用示例)
- [附录 A. OpenAPI 3.1 风格 CLI 契约映射（非 HTTP 接口）](#附录-a-openapi-31-风格-cli-契约映射非-http-接口)

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
| `-recursive` | — | 布尔 | 项目目录递归扫描（默认深度 3 层，1–8 可配） |
| `-scan-depth` | — | 整数 1–8（`0` 使用配置值/默认 3） | 递归扫描深度；优先级高于 `config.ini:SCAN_DEPTH`（S6） |
| `-json` | — | 布尔 | 仅 CLI 模式：stdout 输出结构化 JSON（`function` / `root` / `ok` / `messages` / `error`），退出码不变（S7） |
| `-version` | — | — | 直接输出 `buildinfo.BuildInfo()` 后退出 0（S12） |
| `-h` | `-help` | — | 输出帮助与版本页脚（`vX.Y.Z · build yymmddhhmm`）后退出 |

> `-json` 模式下危险操作确认提示与错误诊断仍走 stderr，避免污染 stdout 的 JSON 文档。

### 1.2 位置参数

| 位置 | 取值 | 说明 |
|------|------|------|
| 1 | `1`-`7` | 功能编号；数字但不在 1–7 → 报错退出 1；非数字 → 按路径解析（GUI 模式下被忽略并提示） |
| 2 | 目录路径 | 本次执行的 `REPO_ROOT`；仅 CLI 模式且非预演时执行成功后写回 `config.ini` |

---

## 2. 退出码定义

共 3 种退出码（`main.go`），`-json` 不改变该契约：

| 退出码 | 场景 |
|:------:|------|
| `0` | 成功结束；含用户在确认提示回答非 `Y`（正常文本模式打印「已取消」；`-json` 模式在 stdout 输出 `ok: true` 且 `messages` 含取消文案） |
| `1` | 五类失败：① `config.ini` 读取失败（禁止回写并以非零码终止；`-h` 仍返回 0）；② 根目录校验失败（不存在 / 非目录 / 文件系统根，含指向根的符号链接）；③ 非交互环境读确认输入失败（stdin EOF / 关闭，**不算用户取消**）；④ `engine.Run` 执行出错；⑤ 非法功能号（数字但不在 1–7） |
| `2` | 标准 `flag` 包解析错误（未知参数、缺少取值、整数参数非法）；由 Go `flag.CommandLine` 直接退出，不经应用自身错误处理 |

脚本侧应把非 0 一律视为失败；「取消」与「成功」同为 0，可用 `-json` 的 `messages` 字段区分（S7 已落地）。

---

## 3. 输出流与日志约定

| 内容 | 流 | 证据 |
|------|----|------|
| 错误信息（`[错误] ...`） | stderr | `internal/i18n.LogDict.ErrorPrefix` |
| 危险操作确认提示（`是否继续? [Y/N]: `） | 默认 stdout；`-json` 下走 stderr | `internal/i18n.LogDict.ConfirmCLI` |
| 正常执行日志（含 `[预演]` / `[n/m]` 进度 / `[警告]`） | 默认 stdout；`-json` 下不直接输出（进 `messages` 与日志文件） | `cmd/repo-manager/cli.go` 的 `runCLI` |
| 全部执行日志（带 `2006-01-02 15:04:05` 时间戳） | 同时追加写 exe 同目录 `repo-manager.log` | `logs.NewFileLogger`（CLI / GUI 入口各调用一次） |

确认交互规则：`TrimSpace` 后必须 `EqualFold "Y"` 才执行；stdin 读取失败 → stderr 提示并退出 1（不静默当作取消）。

`-json` 模式 stdout 固定为单个 JSON 对象；执行日志只进 `messages` 与 `repo-manager.log`，确认提示与错误走 stderr，不直接输出到 stdout。

---

## 4. 配置文件契约

- 路径：可执行文件同目录 `config.ini`（便携式）；格式 INI，`键=值`，`#` 注释保留。
- 键共 5 个：`REPO_ROOT` / `LANG` / `RECURSIVE` / `EXTRA_CLEAN_DIRS` / `SCAN_DEPTH`，语义见 [INSTALL.md §3](./INSTALL.md)。
- 写入为原子替换（唯一临时文件 + `Sync` + `os.Rename`）；读取失败（非「文件不存在」）时 CLI 禁止回写并报错退出，防止损坏配置被覆盖。预演模式不回写配置。
- `RECURSIVE` 接受 `true` / `1` / `yes`（均大小写不敏感）；`EXTRA_CLEAN_DIRS` 支持安全单级目录名与受限相对子路径（逐段校验，拒绝 `..` / 绝对路径 / 保留设备名；删除前逐段拒绝符号链接并确认目标仍在项目内），不合格直接丢弃并告警（S3）；`SCAN_DEPTH` 取 1–8，非法值回落 3（S6）。

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

# `-json` 脚本集成：stdout 固定单个 JSON 对象（日志进 messages，错误进 stderr）
repo-manager -json 3 ~/Github > result.json 2>err.log; echo "exit=$?"; jq -r '.ok, (.messages[]?)' result.json
```

```powershell
# Windows PowerShell 校验退出码
repo-manager -dry-run 6 E:\Github
if ($LASTEXITCODE -ne 0) { "预演失败" }
```

> 所有示例严禁指向真实项目根目录做删除验证；先在隔离临时目录演练（[INSTALL.md §5](./INSTALL.md)）。

---

## 附录 A. OpenAPI 3.1 风格 CLI 契约映射（非 HTTP 接口）

> **适用范围声明**：本产品**没有** HTTP 服务、监听端口、RPC 或请求体（见 [ARCHITECTURE.md §1.3](./ARCHITECTURE.md)）。本附录**不是 HTTP API 定义**，也不提供 `servers` 地址；它把 §1–§5 的 CLI 契约映射为 OpenAPI 3.1 的字段结构，便于按行业标准审阅「路径 / 方法 / 参数 / 返回 / 错误码 / 示例」六要素。不得据此发起网络请求。真正的机器可读结构化输出（`-json`）已由 S7 于 v2.3.11 实现并归档至 [tasks-done.md 第十四节](./tasks-done.md)。

### A.0 字段结构骨架（仅示意，非可请求文档）

```yaml
openapi: 3.1.0
info:
  title: Repo Manager CLI Contract (non-HTTP)
  version: X.Y.Z            # 占位；真实版本以 internal/buildinfo/version.go 为准
  description: 本地 CLI 行为的等效映射，无网络监听。
servers: []                 # 显式为空：产品不提供服务端点
paths:
  /cli/repo-manager:        # 逻辑路径，不是网络地址
    post:                   # 等价「启动一次本地进程」，非 HTTP POST
      parameters:           # 详见 A.2
        - name: lang
          in: flag
        - name: function
          in: positional
      responses:            # 详见 A.3
        "0": { description: 成功 }
        "1": { description: 失败（含局部失败汇总与超时） }
```

> 该骨架只表达字段层级；实际参数、退出码与示例以 A.1–A.4 为准。

### A.1 六要素映射表

| OpenAPI 要素 | 本产品对应物 | 现行文档 |
|--------------|--------------|----------|
| `paths` | 本地命令调用形态 `repo-manager [flags] [功能号 1-7] [路径]`；帮助路径 `repo-manager -h` | §1 |
| `operation.method` | **无 HTTP 方法**；等价动作为「启动一次本地进程并执行功能号」 | §1 / §1.2 |
| `parameters` | flags：`-lang` / `-dry-run` / `-recursive` / `-h`；位置参数：功能号、路径 | §1.1 / §1.2 |
| `requestBody` | 不适用；交互式输入为 stdin 的 `Y/N` 确认（仅功能 1/6/7） | §3 |
| `responses` | stdout 执行日志 + stderr 错误 + 进程退出码，无 MIME body | §2 / §3 |
| `error codes` | 退出码 `0` / `1` 与日志前缀 `[错误]` / `[警告]` / `[预演]` | §2 / §3 |
| `examples` | bash / PowerShell 调用示例 | §5 |

### A.2 参数定义（OpenAPI parameter 风格）

| name | in | required | schema | 约束与副作用 |
|------|----|:--------:|--------|--------------|
| `lang` | flag | 否 | enum `zh` \| `zh-Hant` \| `en` | `tw` / `zh-TW` 归一化为 `zh-Hant`；非法值回落 `zh`；**写回 `config.ini`** |
| `dry-run` | flag | 否 | boolean | 零副作用；危险操作跳过二次确认 |
| `recursive` | flag | 否 | boolean | 递归扫描；深度由 `scanDepth` 控制（默认 3） |
| `scanDepth` | flag | 否 | integer 1–8 | 递归深度；`0` 使用 `config.ini:SCAN_DEPTH` 或默认 3（S6） |
| `json` | flag | 否 | boolean | 仅 CLI：stdout 输出 `jsonResult`，退出码不变（S7） |
| `h` / `help` | flag | 否 | boolean | 输出帮助 + 版本页脚后退出 0；在读取配置之后、错误检查之前处理（仅 `-version` 在读取配置前） |
| `version` | flag | 否 | boolean | 直接输出 `buildinfo.BuildInfo()` 后退出 0；在读取配置之前处理 |
| `function` | positional | 否 | integer `1`–`7` | 缺省启动 GUI；数字但不在 1–7 → 报错退出 1；非数字按路径解析（GUI 模式下忽略并提示） |
| `repoRoot` | positional | 否 | string (path) | 优先级：命令行 > `config.ini:REPO_ROOT` > 平台默认 |

### A.3 返回定义（OpenAPI response 风格）

| 等价状态码 | 触发条件 | 输出 / 副作用 |
|:----------:|----------|---------------|
| `0` | 执行成功；或用户在确认提示回答非 `Y` | 默认 stdout：`[n/m]` 进度与结果日志；`-json` 下 stdout 固定单个 JSON 对象，日志进 `messages` 并追加写 `repo-manager.log`，确认提示与错误走 stderr |
| `1` | 配置读取失败 / 根目录无效 / stdin EOF / `engine.Run` 出错 / 非法功能号 | stderr：`[错误] ...`（`-json` 下亦走 stderr）；配置失败时不回写配置；`-json` 下 stdout 仍为单个 JSON 对象（`ok: false` + `error`） |
| `2` | `flag` 解析错误（未知参数 / 缺值 / 非法整数） | stderr：标准 `flag` 包错误；不经应用错误处理 |

日志行格式为 `YYYY-MM-DD HH:MM:SS <消息>`（`logs.go`）；`Content-Type` 不适用，因为响应是进程输出而非 HTTP 消息体。

`-json` 模式下 stdout 固定为单个对象（对应 `components.schemas.jsonResult`）：

```json
{
  "function": 6,
  "root": "E:\\Github",
  "ok": true,
  "messages": ["[n/m] ..."],
  "error": ""
}
```

字段：`function`（功能号，GUI 模式不适用）、`root`（本次根目录）、`ok`（等价退出状态）、`messages`（执行日志数组）、`error`（仅失败时出现）。危险操作确认提示与错误诊断仍走 stderr，不污染该对象。

### A.4 调用示例

与 §5 一致，例如：

```bash
# 预演（等价「只读调用」）：零副作用
repo-manager -dry-run 6 ~/Github

# 危险操作（等价「写入调用」）：stdin 回答 Y 才执行
repo-manager 1 E:\Github

# 显式语言偏好（副作用：写回 config.ini）
repo-manager -lang zh-Hant
```

> 若后续引入本地守护服务 / HTTP API，须另立符合 OpenAPI 3.1 的服务端文档（`servers` + `paths` + `components.schemas`），并同步登记 [tasks.md](./tasks.md)；本附录不得被改写成虚构的网络接口说明。