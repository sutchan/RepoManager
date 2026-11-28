// internal/config/config.go v2.3.14

// Package config 负责 config.ini 的读写：读旧 → 过滤目标键 → 补写 → 回写，
// 保留注释与未知行（对齐 OPENSPEC §5.3 与 CONTRIBUTING.md §4.5 约定）。
package config

import (
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
