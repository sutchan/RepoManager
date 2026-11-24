// cmd/repo-manager/main_test.go v2.3.6

package main

import (
	"testing"

	"repo-manager/internal/config"
	"repo-manager/internal/platform"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantFn   int
		wantPath string
	}{
		{"empty", nil, 0, ""},
		{"only path", []string{`E:\Projects`}, 0, `E:\Projects`},
		{"fn only", []string{"3"}, 3, ""},
		{"fn and path", []string{"5", `E:\Projects`}, 5, `E:\Projects`},
		{"fn out of range", []string{"9"}, 0, "9"},
		{"fn non numeric", []string{"abc"}, 0, "abc"},
		{"extra args ignored", []string{"1", "path", "extra"}, 1, "path"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, path := parseArgs(tt.args)
			if fn != tt.wantFn || path != tt.wantPath {
				t.Fatalf("parseArgs(%v) = (%d, %q), want (%d, %q)", tt.args, fn, path, tt.wantFn, tt.wantPath)
			}
		})
	}
}

func TestResolveRootPriority(t *testing.T) {
	p := platform.New()
	if fn, _ := parseArgs([]string{"0"}); fn != 0 {
		t.Fatal("功能号 0 应视为无功能")
	}

	// 命令行路径优先于配置；配置优先于平台默认。
	argPath := t.TempDir()
	if got := resolveRoot(p, config.Config{RepoRoot: t.TempDir()}, argPath); got != platform.NormalizeRoot(argPath) {
		t.Fatalf("命令行路径应优先，得到 %q", got)
	}
	cfgPath := t.TempDir()
	if got := resolveRoot(p, config.Config{RepoRoot: cfgPath}, ""); got != platform.NormalizeRoot(cfgPath) {
		t.Fatalf("配置路径应优先于默认值，得到 %q", got)
	}
	if got := resolveRoot(p, config.Config{}, ""); got != p.DefaultRepoRoot() {
		t.Fatalf("无参数应取平台默认，得到 %q", got)
	}
}
