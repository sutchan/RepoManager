// internal/engine/engine.go v2.3.12

// Package engine 实现功能 1 ～ 7 的平台无关业务逻辑；OS 差异一律经 platform.Platform。
package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// nextCachePath 返回 Next.js 项目目录的缓存路径；功能 3/4 共用，避免重复拼接。
func nextCachePath(projectDir string) string { return filepath.Join(projectDir, ".next") }

// Logger 由 UI / CLI 注入，用于输出执行日志。
type Logger interface {
	Logf(format string, args ...any)
}

// StdLogger 默认输出到 stdout（命令行模式）。
type StdLogger struct{}

// Logf 实现 Logger。
func (StdLogger) Logf(format string, args ...any) { fmt.Printf(format+"\n", args...) }

// 功能 6 清理白名单（对齐 OPENSPEC §3.6）。
var cleanDirsF6 = []string{"node_modules", ".next", ".nuxt", "dist", "build"}

// 功能 6 清理锁文件白名单。
var lockFilesF6 = []string{"package-lock.json", "yarn.lock", "pnpm-lock.yaml"}

// 功能 7 清理白名单（对齐 OPENSPEC §3.7；v2.2.0 起含 Angular / Svelte 缓存）。
var cleanDirsF7 = []string{".vite", ".turbo", ".nuxt", ".angular", ".svelte-kit", "dist", "build", "coverage", ".nyc_output"}

// npmGroups 功能 5 安装分组（对齐 OPENSPEC §3.5）。
var npmGroups = [][]string{
	{"react", "react-dom", "next"},
	{"vue", "@vue/cli"},
	{"axios", "express"},
	{"pnpm", "yarn"},
	{"cross-env", "dotenv", "nodemon", "pm2", "rimraf", "mkdirp"},
}

// EnableGlobalNodeModules 功能 1：启用全局 node_modules（对齐 OPENSPEC §3.1）。
// DryRun 下仅输出预演日志，不写环境变量 / .npmrc / 删除目录。
// 局部失败逐项记录并继续，结束时统一返回汇总错误（S10）。
func EnableGlobalNodeModules(p platform.Platform, repoRoot string, lang i18n.Lang, opts Options, log Logger) error {
	d := i18n.Logs(lang)
	ctx := opts.ctx()
	var failed failureCollector
	// S5：每个外部副作用前检查取消，避免用户取消后继续删除或写入。
	if err := ctx.Err(); err != nil {
		return err
	}
	log.Logf("%s", d.F1Step1)
	globalPath, err := p.GlobalNodeModulesPath()
	if err != nil {
		return err
	}
	log.Logf(d.F1Path, globalPath)

	log.Logf("%s", d.F1Step2)
	if err := ctx.Err(); err != nil {
		return err
	}
	if opts.DryRun {
		log.Logf(d.DryRunEnv, "NODE_PATH="+globalPath)
	} else if err := p.SetUserEnv("NODE_PATH", globalPath); err != nil {
		log.Logf(d.EnvSetWarn, err)
		failed.add("NODE_PATH: %v", err)
	}

	log.Logf("%s", d.F1Step3)
	dirs, err := listDirs(p, repoRoot, opts)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if !platform.IsNodeProject(dir) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if opts.DryRun {
			log.Logf(d.DryRunWrite, filepath.Join(dir, ".npmrc"))
			continue
		}
		if err := writeNpmrc(p, dir, globalPath); err != nil {
			log.Logf(d.WriteNpmrcErr, dir, err)
			failed.add("%s: %v", filepath.Join(dir, ".npmrc"), err)
			continue
		}
		log.Logf(d.Configured, filepath.Base(dir))
	}

	log.Logf("%s", d.F1Step4)
	for _, dir := range dirs {
		if err := ctx.Err(); err != nil {
			return err
		}
		nm := filepath.Join(dir, "node_modules")
		if info, err := os.Stat(nm); err == nil && info.IsDir() {
			if opts.DryRun {
				log.Logf(d.DryRunDelete, nm)
				continue
			}
			log.Logf(d.Cleaning, filepath.Base(dir))
			if err := p.RemoveAll(nm); err != nil {
				log.Logf(d.DeleteErr, err)
				failed.add("%s: %v", nm, err)
			}
		}
	}
	return failed.err(lang)
}

// writeNpmrc 写入 .npmrc；tmp= 值已加引号（v2.0.0 修复 L1，见 OPENSPEC §7.3）。
func writeNpmrc(p platform.Platform, projectDir, globalPath string) error {
	content := fmt.Sprintf("prefix=%s\ncache=%s\ntmp=%q\nglobal=true\nprefer-global=true\n",
		globalPath, filepath.Join(globalPath, ".npm-cache"), p.TempDir())
	return os.WriteFile(filepath.Join(projectDir, ".npmrc"), []byte(content), 0o644)
}

// ResetNodeConfig 功能 2：重置 Node 配置（对齐 OPENSPEC §3.2）。
// DryRun 下仅输出预演日志，不删除环境变量 / .npmrc。
func ResetNodeConfig(p platform.Platform, repoRoot string, lang i18n.Lang, opts Options, log Logger) error {
	d := i18n.Logs(lang)
	ctx := opts.ctx()
	var failed failureCollector
	if err := ctx.Err(); err != nil {
		return err
	}
	log.Logf("%s", d.F2Step1)
	if opts.DryRun {
		log.Logf(d.DryRunEnvDel, "NODE_PATH")
	} else if err := p.RemoveUserEnv("NODE_PATH"); err != nil {
		log.Logf(d.EnvDelWarn, err)
		failed.add("NODE_PATH: %v", err)
	}
	log.Logf("%s", d.F2Step2)
	if err := ctx.Err(); err != nil {
		return err
	}
	dirs, err := listDirs(p, repoRoot, opts)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if err := ctx.Err(); err != nil {
			return err
		}
		npmrc := filepath.Join(dir, ".npmrc")
		if opts.DryRun {
			if _, err := os.Stat(npmrc); err == nil {
				log.Logf(d.DryRunDelete, npmrc)
			}
			continue
		}
		if err := p.RemoveFile(npmrc); err == nil {
			log.Logf(d.Removed, filepath.Base(dir))
		} else if !os.IsNotExist(err) {
			failed.add("%s: %v", npmrc, err)
		}
	}
	return failed.err(lang)
}
