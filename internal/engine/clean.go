// internal/engine/clean.go v2.3.12

package engine

import (
	"os"
	"path/filepath"

	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// CleanNextCache 功能 3：清理并禁用 .next 缓存（对齐 OPENSPEC §3.3）。
// DryRun 下仅输出预演日志，不删除 / 不写入任何文件。
func CleanNextCache(p platform.Platform, repoRoot string, lang i18n.Lang, opts Options, log Logger) error {
	d := i18n.Logs(lang)
	ctx := opts.ctx()
	var failed failureCollector
	if err := ctx.Err(); err != nil {
		return err
	}
	log.Logf("%s", d.F3Start)
	dirs, err := listDirs(p, repoRoot, opts)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if !platform.IsNextProject(dir) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		next := nextCachePath(dir)
		if opts.DryRun {
			if _, err := os.Stat(next); err == nil {
				log.Logf(d.DryRunDelete, next)
			}
			log.Logf(d.DryRunWrite, next)
			continue
		}
		log.Logf(d.Processing, filepath.Base(dir))
		// Lstat 不跟随符号链接：.next 若被替换为链接，只删除链接本身，
		// 避免下面的写入与属性修改破坏链接指向的项目外文件。
		if info, err := os.Lstat(next); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				if err := p.RemoveFile(next); err != nil {
					log.Logf(d.RemoveErr, err)
					failed.add("%s: %v", next, err)
					continue
				}
			} else if info.IsDir() {
				if err := p.RemoveAll(next); err != nil {
					log.Logf(d.RemoveDirErr, err)
					failed.add("%s: %v", next, err)
					continue
				}
			} else if err := p.ClearHidden(next); err != nil {
				log.Logf(d.ClearAttrWarn, err)
				failed.add("%s ClearHidden: %v", next, err)
			}
		}
		// 写入同名文件阻止 Next.js 重建缓存目录，再设置隐藏属性
		if err := os.WriteFile(next, []byte("\n"), 0o644); err != nil {
			log.Logf(d.WriteNextErr, err)
			failed.add("%s: %v", next, err)
			continue
		}
		if err := p.SetHidden(next); err != nil {
			log.Logf(d.SetAttrWarn, err)
			failed.add("%s SetHidden: %v", next, err)
		}
	}
	return failed.err(lang)
}

// RestoreNextCache 功能 4：恢复 Next.js 缓存（对齐 OPENSPEC §3.4）。
// 关键：必须先清除隐藏属性再删除，否则删除静默失败（沿用 v1.2.3 修复约束）。
// DryRun 下仅输出预演日志，不删除任何文件。
func RestoreNextCache(p platform.Platform, repoRoot string, lang i18n.Lang, opts Options, log Logger) error {
	d := i18n.Logs(lang)
	ctx := opts.ctx()
	var failed failureCollector
	if err := ctx.Err(); err != nil {
		return err
	}
	log.Logf("%s", d.F4Start)
	dirs, err := listDirs(p, repoRoot, opts)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if !platform.IsNextProject(dir) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		next := nextCachePath(dir)
		// Lstat 不跟随符号链接，避免把链接指向的外部文件当作占位文件处理。
		info, err := os.Lstat(next)
		if err != nil || info.IsDir() {
			continue // 不存在或是正常缓存目录，无需处理
		}
		if opts.DryRun {
			log.Logf(d.DryRunDelete, next)
			continue
		}
		log.Logf(d.Restoring, filepath.Base(dir))
		if info.Mode()&os.ModeSymlink != 0 {
			// 功能 3 写入的占位若为链接，只删除链接本身。
			if err := p.RemoveFile(next); err != nil && !os.IsNotExist(err) {
				log.Logf(d.RemoveErr, err)
				failed.add("%s: %v", next, err)
			}
			continue
		}
		if err := p.ClearHidden(next); err != nil {
			log.Logf(d.ClearAttrWarn, err)
			failed.add("%s ClearHidden: %v", next, err)
		}
		if err := p.RemoveFile(next); err != nil && !os.IsNotExist(err) {
			log.Logf(d.RemoveErr, err)
			failed.add("%s: %v", next, err)
		}
	}
	return failed.err(lang)
}

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
