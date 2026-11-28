// internal/engine/whitelist.go v2.3.14

package engine

import (
	"os"
	"path/filepath"

	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// CleanAllCaches 功能 6：清理所有项目缓存（对齐 OPENSPEC §3.6）。
func CleanAllCaches(p platform.Platform, repoRoot string, lang i18n.Lang, opts Options, log Logger) error {
	return cleanByWhitelist(p, repoRoot, cleanDirsF6, lockFilesF6, lang, opts, log)
}

// CleanViteCache 功能 7：清理 Vite/构建缓存（对齐 OPENSPEC §3.7）。
func CleanViteCache(p platform.Platform, repoRoot string, lang i18n.Lang, opts Options, log Logger) error {
	return cleanByWhitelist(p, repoRoot, cleanDirsF7, nil, lang, opts, log)
}

// cleanByWhitelist 遍历 REPO_ROOT 直接子目录，仅清理白名单内的目录与文件。
// DryRun 下仅输出预演日志，不删除任何文件。
func cleanByWhitelist(p platform.Platform, repoRoot string, dirNames, fileNames []string, lang i18n.Lang, opts Options, log Logger) error {
	d := i18n.Logs(lang)
	ctx := opts.ctx()
	var failed failureCollector
	if err := ctx.Err(); err != nil {
		return err
	}
	projects, err := listDirs(p, repoRoot, opts)
	if err != nil {
		return err
	}
	if len(opts.ExtraCleanDirs) > 0 {
		// V3：自定义清理目录仅限项目目录内部；逐项校验，拒绝路径穿越。
		dirNames = append([]string{}, dirNames...)
		for _, name := range opts.ExtraCleanDirs {
			if !platform.IsSafeCleanDirPath(name) {
				log.Logf(d.UnsafeCleanDir, name)
				continue
			}
			dirNames = append(dirNames, name)
		}
	}
	for _, project := range projects {
		if err := ctx.Err(); err != nil {
			return err
		}
		// 删除前逐段拒绝符号链接并确认目标仍位于项目内，避免越界删除。
		type targetEntry struct {
			name string
			path string
		}
		targets := make([]targetEntry, 0, len(dirNames))
		for _, name := range dirNames {
			target, ok := platform.ResolveCleanTarget(project, name)
			if !ok {
				log.Logf(d.UnsafeCleanDir, name)
				continue
			}
			targets = append(targets, targetEntry{name: name, path: target})
		}
		if opts.DryRun {
			for _, t := range targets {
				if _, err := os.Stat(t.path); err == nil {
					log.Logf(d.DryRunDelete, t.path)
				}
			}
			for _, name := range fileNames {
				target := filepath.Join(project, name)
				if _, err := os.Stat(target); err == nil {
					log.Logf(d.DryRunDelete, target)
				}
			}
			continue
		}
		log.Logf(d.Cleaning, filepath.Base(project))
		for _, t := range targets {
			if _, err := os.Stat(t.path); err == nil {
				if err := p.RemoveAll(t.path); err != nil {
					log.Logf(d.RemoveNameErr, t.name, err)
					failed.add("%s: %v", t.path, err)
				}
			}
		}
		for _, name := range fileNames {
			target := filepath.Join(project, name)
			if err := p.RemoveFile(target); err == nil {
				log.Logf(d.Removed, name)
			} else if !os.IsNotExist(err) {
				failed.add("%s: %v", target, err)
			}
		}
	}
	return failed.err(lang)
}
