// internal/config/save.go v2.3.14

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Save 以「读旧 → 过滤目标键 → 补写 → 临时文件 → 回写」方式保存，去重同名键。
// 读旧失败且非「文件不存在」时直接返回，禁止用内存中的配置覆盖磁盘文件。
func Save(path string, cfg Config) error {
	var lines []string
	data, err := os.ReadFile(path)
	if err == nil {
		lines = strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("读取现有配置失败，已取消保存: %w", err)
	}
	repoWritten, langWritten, recWritten, extraWritten, depthWritten := false, false, false, false, false
	extraLine := "EXTRA_CLEAN_DIRS=" + strings.Join(cfg.ExtraCleanDirs, ",")
	recLine := "RECURSIVE=false"
	depthLine := "SCAN_DEPTH=" + strconv.Itoa(cfg.ScanDepth)
	if cfg.Recursive {
		recLine = "RECURSIVE=true"
	}
	for i, line := range lines {
		key, _, ok := parseKeyValue(line)
		if !ok {
			continue
		}
		switch key {
		case "REPO_ROOT":
			if !repoWritten {
				lines[i] = "REPO_ROOT=" + cfg.RepoRoot
				repoWritten = true
			} else {
				lines[i] = ""
			}
		case "LANG":
			if !langWritten {
				lines[i] = "LANG=" + cfg.Lang
				langWritten = true
			} else {
				lines[i] = ""
			}
		case "RECURSIVE":
			if !recWritten {
				lines[i] = recLine
				recWritten = true
			} else {
				lines[i] = ""
			}
		case "EXTRA_CLEAN_DIRS":
			if !extraWritten {
				lines[i] = extraLine
				extraWritten = true
			} else {
				lines[i] = ""
			}
		case "SCAN_DEPTH":
			if !depthWritten {
				lines[i] = depthLine
				depthWritten = true
			} else {
				lines[i] = ""
			}
		}
	}
	if !repoWritten {
		lines = append(lines, "REPO_ROOT="+cfg.RepoRoot)
	}
	if !langWritten {
		lines = append(lines, "LANG="+cfg.Lang)
	}
	if !recWritten {
		lines = append(lines, recLine)
	}
	if !extraWritten {
		lines = append(lines, extraLine)
	}
	if !depthWritten {
		lines = append(lines, depthLine)
	}
	content := strings.Join(lines, "\n")
	// 唯一临时文件名 + Sync：多实例并发保存时不再争用固定的 config.ini.tmp，
	// 崩溃时也不会留下半截配置；失败时清理残留。
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.tmp")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		return fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("同步临时文件失败: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("关闭临时文件失败: %w", err)
	}
	return os.Rename(tmpName, path)
}
