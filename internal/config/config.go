// internal/config/config.go v2.3.11

// Package config 负责 config.ini 的读写：读旧 → 过滤目标键 → 补写 → 回写，
// 保留注释与未知行（对齐 OPENSPEC §5.3 与 CONTRIBUTING.md §4.5 约定）。
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// Config 产品运行配置。
type Config struct {
	RepoRoot       string   // 仓库根目录
	Lang           string   // 界面语言：zh / zh-Hant / en
	Recursive      bool     // 递归扫描 monorepo 深层项目（V1，默认关闭）
	ExtraCleanDirs []string // 自定义清理目录名 / 受限相对子路径（V3）
	ScanDepth      int      // 递归扫描深度（S6，1–8，默认 3）
}

// 递归扫描深度边界（S6）：默认 3 层，最大 8 层。
const (
	defaultScanDepth = 3
	maxScanDepth     = 8
)

// Default 返回默认配置（REPO_ROOT 留空，由调用方按平台取默认目录）。
func Default() Config { return Config{RepoRoot: "", Lang: "zh", ScanDepth: 3} }

// Path 返回与可执行文件同目录的 config.ini 路径（对齐 v1.x 行为）。
func Path() string {
	exe, err := os.Executable()
	if err != nil {
		return "config.ini"
	}
	return filepath.Join(filepath.Dir(exe), "config.ini")
}

// Load 读取配置；文件不存在时返回默认值。
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := parseKeyValue(line)
		if !ok {
			continue
		}
		switch key {
		case "REPO_ROOT":
			cfg.RepoRoot = strings.TrimSpace(value)
		case "LANG":
			cfg.Lang = strings.TrimSpace(value)
		case "RECURSIVE":
			v := strings.TrimSpace(value)
			cfg.Recursive = strings.EqualFold(v, "true") || v == "1" || strings.EqualFold(v, "yes")
		case "EXTRA_CLEAN_DIRS":
			cfg.ExtraCleanDirs = parseDirList(value)
		case "SCAN_DEPTH":
			if n, err := strconv.Atoi(value); err == nil && n >= 1 && n <= maxScanDepth {
				cfg.ScanDepth = n
			} else {
				cfg.ScanDepth = defaultScanDepth
			}
		}
	}
	// 语言枚举校验：非法值回落简中，繁中别名归一化为 zh-Hant。
	cfg.Lang = string(i18n.Normalize(cfg.Lang))
	return cfg, nil
}

// parseKeyValue 解析 INI 行：跳过空行与注释，返回归一化 key 与去空白 value（S16）。
func parseKeyValue(line string) (key, value string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
		return "", "", false
	}
	key, value, ok = strings.Cut(line, "=")
	if !ok {
		return "", "", false
	}
	return strings.ToUpper(strings.TrimSpace(key)), strings.TrimSpace(value), true
}

// parseDirList 解析逗号分隔的目录名列表（去空白、去空项），
// 并丢弃会造成路径穿越的非法目录名，避免清理目标指向项目目录之外。
func parseDirList(value string) []string {
	var names []string
	for _, part := range strings.Split(value, ",") {
		name := strings.TrimSpace(part)
		if !platform.IsSafeCleanDirPath(name) {
			continue
		}
		names = append(names, name)
	}
	return names
}

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
