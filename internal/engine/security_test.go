// internal/engine/security_test.go v2.3.14

package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"repo-manager/internal/i18n"
)

// TestEnableGlobalNodeModulesSkipsNonNodeProjects 确认功能 1 只删除 Node 项目的本地依赖。
func TestEnableGlobalNodeModulesSkipsNonNodeProjects(t *testing.T) {
	s := newStub(t)
	node := newProject(t, s.root, "app", map[string][]byte{"package.json": []byte("{}")})
	plain := newProject(t, s.root, "docs", nil)
	for _, dir := range []string{node, plain} {
		if err := os.MkdirAll(filepath.Join(dir, "node_modules"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := EnableGlobalNodeModules(s, s.root, i18n.Zh, Options{}, &captureLogger{}); err != nil {
		t.Fatal(err)
	}
	if exists(t, filepath.Join(node, "node_modules")) {
		t.Fatal("Node 项目的 node_modules 应被删除")
	}
	if !exists(t, filepath.Join(plain, "node_modules")) {
		t.Fatal("非 Node 目录下的 node_modules 必须保留")
	}
}

// TestCleanNextCacheSymlinkDoesNotOverwriteTarget 确认 .next 符号链接不破坏外部文件。
func TestCleanNextCacheSymlinkDoesNotOverwriteTarget(t *testing.T) {
	s := newStub(t)
	project := newProject(t, s.root, "web", map[string][]byte{"next.config.js": []byte("")})
	outside := filepath.Join(t.TempDir(), "important.txt")
	if err := os.WriteFile(outside, []byte("keep-me"), 0o644); err != nil {
		t.Fatal(err)
	}
	next := filepath.Join(project, ".next")
	if err := os.Symlink(outside, next); err != nil {
		t.Skip("当前环境不支持创建符号链接:", err)
	}
	if err := CleanNextCache(s, s.root, i18n.Zh, Options{}, &captureLogger{}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(outside)
	if err != nil || string(data) != "keep-me" {
		t.Fatalf("符号链接指向的外部文件被破坏: %q, %v", data, err)
	}
	info, err := os.Lstat(next)
	if err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf(".next 应被替换为普通占位文件，得到 %v / %v", info, err)
	}
}

// TestRunRecursiveLogUsesConfiguredDepth 确认递归日志深度与配置一致。
func TestRunRecursiveLogUsesConfiguredDepth(t *testing.T) {
	s := newStub(t)
	logger := &captureLogger{}
	opts := Options{Recursive: true, ScanDepth: 7}
	if err := Run(s, 6, s.root, i18n.Zh, opts, logger); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(logger.lines, "\n"), "深度 7 层") {
		t.Fatalf("日志应报告配置深度 7，得到 %v", logger.lines)
	}
}
