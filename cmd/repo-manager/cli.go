// cmd/repo-manager/cli.go v2.3.14

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"repo-manager/internal/buildinfo"
	"repo-manager/internal/config"
	"repo-manager/internal/engine"
	"repo-manager/internal/i18n"
	"repo-manager/internal/logs"
	"repo-manager/internal/platform"
)

// runCLI 命令行模式：危险操作保留二次确认，预演模式跳过（对齐 OPENSPEC §4.5 / §4.7）。
func runCLI(p platform.Platform, cfgPath string, cfg *config.Config, fn int, root string, lang i18n.Lang, opts engine.Options, useJSON bool) {
	d := i18n.Logs(lang)
	if err := platform.ValidateRoot(root); err != nil {
		// S17：平台层只给哨兵错误，CLI 按当前语言本地化后再输出。
		localized := errors.New(engine.DescribeError(lang, err))
		if useJSON {
			emitJSONResult(fn, root, nil, localized)
		} else {
			fmt.Fprintln(os.Stderr, localized)
		}
		os.Exit(1)
	}
	if engine.IsDangerous(fn) && !opts.DryRun {
		if useJSON {
			fmt.Fprintln(os.Stderr, d.ConfirmCLI)
		} else {
			fmt.Print(d.ConfirmCLI)
		}
		reader := bufio.NewReader(os.Stdin)
		answer, readErr := reader.ReadString('\n')
		if !strings.EqualFold(strings.TrimSpace(answer), "Y") {
			// 非交互环境（stdin 关闭 / EOF）不能静默当作「用户取消」，否则脚本会误判成功。
			if readErr != nil {
				fmt.Fprintln(os.Stderr, d.ConfirmUnavailable)
				os.Exit(1)
			}
			if useJSON {
				emitJSONResult(fn, root, []string{d.Cancelled}, nil)
			} else {
				fmt.Println(d.Cancelled)
			}
			return
		}
	}
	// V4：执行日志同时追加写入可执行文件同目录的 repo-manager.log。
	// -json 模式下 stdout 必须是单个 JSON 文档，因此文件 logger 不再转发到 stdout；
	// 日志只进文件与内存收集器，最终由 emitJSONResult 统一输出。
	var downstream engine.Logger = engine.StdLogger{}
	if useJSON {
		downstream = nil
	}
	fileLog := logs.NewFileLogger(logs.DefaultPath(), downstream)
	defer fileLog.Close()
	logger := engine.Logger(fileLog)
	var mem *memLogger
	if useJSON {
		mem = &memLogger{}
		logger = teeLogger{primary: fileLog, secondary: mem}
	}
	runErr := engine.Run(p, fn, root, lang, opts, logger)
	if useJSON {
		var msgs []string
		if mem != nil {
			msgs = mem.msgs
		}
		emitJSONResult(fn, root, msgs, runErr)
	}
	if runErr != nil {
		if !useJSON {
			fmt.Fprintf(os.Stderr, d.ErrorPrefix+"\n", runErr)
		}
		os.Exit(1)
	}
	// 命令行模式记住本次使用的目录；预演不写盘（对齐 OPENSPEC §4.7「不修改任何文件」）。
	if opts.DryRun {
		return
	}
	cfg.RepoRoot = root
	if err := config.Save(cfgPath, *cfg); err != nil {
		fmt.Fprintf(os.Stderr, d.SaveCfgErr, err)
	}
}

// printHelp 输出帮助信息；CLI 帮助固定英文，界面文案取自 internal/i18n。
func printHelp() {
	fmt.Println("Repo Manager")
	fmt.Println()
	fmt.Printf("Version: %s\n\n", buildinfo.BuildInfo())
	fmt.Println("Usage: repo-manager [options] [function] [path]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  1-7       Run corresponding function directly (CLI mode)")
	fmt.Println("  -lang zh  Set language to Simplified Chinese")
	fmt.Println("  -lang zh-Hant  Set language to Traditional Chinese")
	fmt.Println("  -lang en  Set language to English")
	fmt.Println("  -dry-run  Preview actions without modifying anything")
	fmt.Println("  -recursive  Recursive scan (monorepo, default depth 3)")
	fmt.Println("  -scan-depth N  Recursive depth 1-8 (0 uses config.ini / default 3)")
	fmt.Println("  -version  Show version and build info")
	fmt.Println("  -json     Emit structured JSON on stdout (CLI mode only)")
	fmt.Println("  -h        Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  repo-manager                    # GUI mode")
	fmt.Println("  repo-manager 1                  # CLI: enable global node_modules")
	fmt.Println("  repo-manager 5 E:\\Projects      # CLI with path")
}
