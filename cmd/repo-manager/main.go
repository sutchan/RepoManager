// cmd/repo-manager/main.go v2.3.16

// 命令 repo-manager 是产品入口：解析参数后分流到命令行模式或 Fyne GUI。
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"repo-manager/internal/buildinfo"
	"repo-manager/internal/config"
	"repo-manager/internal/engine"
	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

func main() {
	langFlag := flag.String("lang", "", "界面语言 zh|zh-Hant|en（非预演时写入 config.ini）")
	dryRunFlag := flag.Bool("dry-run", false, "预演模式：仅显示将要执行的操作，不修改任何文件")
	recursiveFlag := flag.Bool("recursive", false, "递归扫描 monorepo 深层项目（深度可配，默认 3 层）")
	scanDepthFlag := flag.Int("scan-depth", 0, "递归扫描深度（1-8；0 使用 config.ini 或默认 3）")
	jsonFlag := flag.Bool("json", false, "CLI 结构化 JSON 输出到 stdout（退出码不变）")
	showHelp := flag.Bool("h", false, "显示帮助")
	flag.BoolVar(showHelp, "help", false, "显示帮助")
	showVersion := flag.Bool("version", false, "显示版本与构建信息")
	flag.Parse()
	if *showVersion {
		fmt.Println(buildinfo.BuildInfo())
		return
	}

	cfgPath := config.Path()
	cfg, err := config.Load(cfgPath)
	if *showHelp {
		printHelp()
		return
	}
	lang := i18n.Normalize(cfg.Lang)
	d := i18n.Logs(lang)
	if err != nil {
		// 配置文件读取失败时禁止任何回写，避免用默认值覆盖用户设置（对齐 OPENSPEC §5.3）。
		fmt.Fprintf(os.Stderr, d.LoadCfgErr, err)
		// 配置损坏属于真实失败，必须以非零码退出，避免脚本把失败误判为成功。
		os.Exit(1)
	}
	lang = applyLanguageOption(cfgPath, &cfg, *langFlag, *dryRunFlag)
	// 选项优先级：命令行参数 > config.ini > 默认关闭（对齐 OPENSPEC §4.5 / §5.2）。
	scanDepth := cfg.ScanDepth
	if *scanDepthFlag > 0 {
		scanDepth = *scanDepthFlag
	}
	opts := engine.Options{
		DryRun:         *dryRunFlag,
		Recursive:      *recursiveFlag || cfg.Recursive,
		ExtraCleanDirs: cfg.ExtraCleanDirs,
		ScanDepth:      scanDepth,
	}
	p := platform.New()
	fn, pathArg := parseArgs(flag.Args())
	// S12：非法功能号与被 GUI 忽略的位置参数必须可见，避免用户输入静默失效。
	if len(flag.Args()) > 0 {
		if n, err := strconv.Atoi(flag.Args()[0]); err == nil && (n < 1 || n > 7) {
			fmt.Fprintf(os.Stderr, "invalid function number %d (valid: 1-7)\n", n)
			os.Exit(1)
		}
		if fn == 0 && pathArg != "" {
			fmt.Fprintf(os.Stderr, "note: no function number given; positional path %q is ignored in GUI mode\n", pathArg)
		}
	}
	root := resolveRoot(p, cfg, pathArg)

	if fn > 0 {
		runCLI(p, cfgPath, &cfg, fn, root, lang, opts, *jsonFlag)
		return
	}
	runGUI(cfg)
}

// applyLanguageOption 应用 -lang 并按需持久化；预演模式只影响本次运行，不回写配置（对齐 OPENSPEC §4.7）。
func applyLanguageOption(cfgPath string, cfg *config.Config, requested string, dryRun bool) i18n.Lang {
	lang := i18n.Normalize(cfg.Lang)
	if requested == "" {
		return lang
	}
	lang = i18n.Normalize(requested)
	cfg.Lang = string(lang)
	if dryRun {
		return lang
	}
	if saveErr := config.Save(cfgPath, *cfg); saveErr != nil {
		fmt.Fprintf(os.Stderr, i18n.Logs(lang).SaveCfgErr, saveErr)
	}
	return lang
}

// parseArgs 解析位置参数：[功能号 1-7] [路径]。
func parseArgs(args []string) (fn int, path string) {
	if len(args) == 0 {
		return 0, ""
	}
	if n, err := strconv.Atoi(args[0]); err == nil && n >= 1 && n <= 7 {
		fn = n
		args = args[1:]
	}
	if len(args) > 0 {
		path = args[0]
	}
	return fn, path
}

// resolveRoot 确定仓库根目录：命令行路径 → config.ini → 平台默认值（对齐 OPENSPEC §4.2）。
func resolveRoot(p platform.Platform, cfg config.Config, pathArg string) string {
	if pathArg != "" {
		return platform.NormalizeRoot(pathArg)
	}
	if cfg.RepoRoot != "" {
		return platform.NormalizeRoot(cfg.RepoRoot)
	}
	return p.DefaultRepoRoot()
}
