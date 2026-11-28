// internal/platform/platform.go v2.3.12

// Package platform 收口全部操作系统差异（目录扫描、删除、隐藏属性、环境变量持久化）。
// engine 只依赖本接口，禁止直接调用 OS 专属命令；新增平台仅需增加对应 build tag 实现。
package platform

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// 平台层只返回英文哨兵错误，文案由上层按语言翻译；本包禁止依赖 i18n（S17）。
var (
	ErrRootFilesystem = errors.New("filesystem root rejected")
	ErrRootNotExist   = errors.New("path does not exist")
	ErrRootNotDir     = errors.New("path is not a directory")
	ErrNpmRoot        = errors.New("npm root -g failed")
	ErrCommand        = errors.New("command failed")
	ErrOpenFile       = errors.New("open file failed")
)

// listProjectDirs 枚举 repoRoot 的直接子目录；Windows / Unix 实现共用。
func listProjectDirs(repoRoot string) ([]string, error) {
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join(repoRoot, e.Name()))
		}
	}
	return dirs, nil
}

// globalNodeModulesPath 返回 npm 全局模块路径；Windows / Unix 实现共用。
func globalNodeModulesPath() (string, error) {
	out, err := exec.Command("npm", "root", "-g").Output()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNpmRoot, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// Platform 抽象操作系统差异，是 engine 与操作系统之间的唯一边界。
type Platform interface {
	// DefaultRepoRoot 返回默认仓库根目录（Windows: E:\Github；Unix: ~/Github）。
	DefaultRepoRoot() string
	// ListProjectDirs 枚举 repoRoot 的直接子目录（仅目录，不递归）。
	ListProjectDirs(repoRoot string) ([]string, error)
	// RemoveAll 删除目录及其全部内容。
	RemoveAll(path string) error
	// RemoveFile 删除单个文件。
	RemoveFile(path string) error
	// SetHidden 设置隐藏属性（功能 3）。
	SetHidden(path string) error
	// ClearHidden 清除隐藏属性（功能 4，必须先于删除，否则删除静默失败）。
	ClearHidden(path string) error
	// TempDir 返回系统临时目录（写入 .npmrc 的 tmp=）。
	TempDir() string
	// SetUserEnv 持久化写入用户级环境变量。
	SetUserEnv(key, value string) error
	// RemoveUserEnv 删除用户级环境变量。
	RemoveUserEnv(key string) error
	// GlobalNodeModulesPath 返回 `npm root -g` 的全局模块路径。
	GlobalNodeModulesPath() (string, error)
}

// deepSkipDirs 递归枚举时跳过的依赖 / 构建产物目录名（避免深入依赖树与输出目录）。
var deepSkipDirs = map[string]bool{
	"node_modules": true, ".next": true, ".nuxt": true, "dist": true, "build": true,
	".vite": true, ".turbo": true, ".angular": true, ".svelte-kit": true,
	"coverage": true, ".nyc_output": true, ".git": true,
}

// ListProjectDirsDeep 递归枚举 repoRoot 下 depthLimit 层内的目录（1 = 仅直接子目录，
// 2 = 含孙目录，V1 monorepo 支持）。跳过 deepSkipDirs 中的目录；单个子树读取失败
// 不影响其余结果。返回排序后的路径列表（v2.3.0，对齐 OPENSPEC §3.0）。
func ListProjectDirsDeep(repoRoot string, depthLimit int) ([]string, error) {
	var dirs []string
	seen := map[string]bool{}
	var walk func(dir string, depth int)
	walk = func(dir string, depth int) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return // 子树不可读时跳过，不影响其余枚举
		}
		for _, e := range entries {
			if !e.IsDir() || deepSkipDirs[e.Name()] {
				continue
			}
			full := filepath.Join(dir, e.Name())
			if !seen[full] {
				seen[full] = true
				dirs = append(dirs, full)
			}
			if depth < depthLimit {
				walk(full, depth+1)
			}
		}
	}
	walk(repoRoot, 1)
	sort.Strings(dirs)
	return dirs, nil
}

// HasFile 判断 dir 下是否存在指定文件（探测 package.json / next.config.*）。
func HasFile(dir, name string) bool {
	info, err := os.Stat(filepath.Join(dir, name))
	return err == nil && !info.IsDir()
}

// IsNodeProject 判断目录是否为 Node 项目（含 package.json）。
func IsNodeProject(dir string) bool { return HasFile(dir, "package.json") }

// IsNextProject 判断目录是否为 Next.js 项目（含 next.config.js 或 next.config.ts）。
func IsNextProject(dir string) bool {
	return HasFile(dir, "next.config.js") || HasFile(dir, "next.config.ts")
}

// NormalizeRoot 规范化仓库根目录：去空白、去末尾分隔符、转绝对路径（对齐 OPENSPEC §7.1）。
// 文件系统根（如 "/" 或 "X:\"）去分隔符后会变空；此时保留根本身，禁止退化为当前工作目录。
func NormalizeRoot(root string) string {
	root = strings.TrimSpace(root)
	if root == "" {
		return ""
	}
	trimmed := strings.TrimRight(root, `\/`)
	if trimmed == "" {
		// 输入仅为分隔符 → 是文件系统根，用 filepath.Clean 保留根语义。
		if cleaned := filepath.Clean(root); cleaned != "." {
			return cleaned
		}
		return root
	}
	if abs, err := filepath.Abs(trimmed); err == nil {
		return abs
	}
	return trimmed
}

// ValidateRoot 校验仓库根目录存在、为目录且不是文件系统根（危险操作的最小范围约束）。
func ValidateRoot(root string) error {
	if IsFilesystemRoot(root) {
		return fmt.Errorf("%w: %s", ErrRootFilesystem, root)
	}
	// 符号链接可能指向文件系统根；词法判定会被绕过，解析后再校验一次。
	if resolved, err := filepath.EvalSymlinks(root); err == nil && IsFilesystemRoot(resolved) {
		return fmt.Errorf("%w: %s", ErrRootFilesystem, root)
	}
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRootNotExist, root)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrRootNotDir, root)
	}
	return nil
}

// IsFilesystemRoot 判断路径是否为文件系统根（"/"、"X:\" 或 UNC 共享根）。
// 用 filepath.Dir 自反性判定：根目录的父目录仍是自身。
func IsFilesystemRoot(root string) bool {
	if root == "" {
		return false
	}
	cleaned := filepath.Clean(root)
	return filepath.Dir(cleaned) == cleaned
}
