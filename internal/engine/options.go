// internal/engine/options.go v2.3.12

package engine

import (
	"context"
	"fmt"

	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// Options 功能执行选项。
type Options struct {
	// DryRun 预演模式：仅输出将要执行的操作，不修改任何文件与环境变量
	// （对齐 OPENSPEC §4.7）；危险操作在预演下跳过二次确认。
	DryRun bool
	// Recursive 递归扫描 monorepo 深层项目（V1，3 层深度，默认关闭；
	// 对齐 OPENSPEC §3.0）。
	Recursive bool
	// ExtraCleanDirs 自定义清理目录名 / 受限相对子路径（V3 / S3），
	// 并入功能 6 / 7 白名单；对齐 OPENSPEC §3.6 / §3.7。
	ExtraCleanDirs []string
	// ScanDepth 递归扫描深度（S6，1–8，0 表示回落到默认 3）。
	ScanDepth int
	// Context 由 CLI / GUI 注入的取消上下文（S5）；nil 时等价于 context.Background。
	Context context.Context
}

// ctx 返回可用的上下文，避免调用方显式判空。
func (o Options) ctx() context.Context {
	if o.Context != nil {
		return o.Context
	}
	return context.Background()
}

// scanDepth 返回有效扫描深度（1-8，非法值回落默认 3），枚举与日志共用同一口径。
func scanDepth(opts Options) int {
	if opts.ScanDepth < 1 || opts.ScanDepth > 8 {
		return 3
	}
	return opts.ScanDepth
}

// listDirs 按选项选择枚举方式：Recursive 启用指定深度的递归（默认 3 层，覆盖
// REPO_ROOT/monorepo/packages/app 这类典型 monorepo 结构），否则仅直接子目录。
func listDirs(p platform.Platform, repoRoot string, opts Options) ([]string, error) {
	if opts.Recursive {
		return platform.ListProjectDirsDeep(repoRoot, scanDepth(opts))
	}
	return p.ListProjectDirs(repoRoot)
}

// Run 按功能编号（1-7）执行对应功能，lang 决定执行日志语言。
// 前置条件：调用方须先用 platform.NormalizeRoot + platform.ValidateRoot 校验 repoRoot。
func Run(p platform.Platform, fn int, repoRoot string, lang i18n.Lang, opts Options, log Logger) error {
	if err := opts.ctx().Err(); err != nil {
		return err
	}
	d := i18n.Logs(lang)
	if opts.DryRun {
		log.Logf("%s", d.DryRunOn)
	}
	if opts.Recursive {
		log.Logf(d.RecursiveOn, scanDepth(opts))
	}
	switch fn {
	case 1:
		return EnableGlobalNodeModules(p, repoRoot, lang, opts, log)
	case 2:
		return ResetNodeConfig(p, repoRoot, lang, opts, log)
	case 3:
		return CleanNextCache(p, repoRoot, lang, opts, log)
	case 4:
		return RestoreNextCache(p, repoRoot, lang, opts, log)
	case 5:
		return InstallGlobalDeps(lang, opts, log)
	case 6:
		return CleanAllCaches(p, repoRoot, lang, opts, log)
	case 7:
		return CleanViteCache(p, repoRoot, lang, opts, log)
	default:
		return fmt.Errorf("unknown function: %d", fn)
	}
}

// IsDangerous 返回功能是否为危险操作（需要二次确认）：1 / 6 / 7。
func IsDangerous(fn int) bool { return fn == 1 || fn == 6 || fn == 7 }
