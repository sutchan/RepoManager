// internal/platform/security_test.go v2.3.14

package platform

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveCleanTargetRejectsSymlink 确认中间符号链接不能把清理目标导向项目外。
func TestResolveCleanTargetRejectsSymlink(t *testing.T) {
	base := t.TempDir()
	project := filepath.Join(base, "project")
	outside := filepath.Join(base, "outside")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(project, "link")
	if err := os.Symlink(outside, link); err != nil {
		t.Skip("当前环境不支持创建符号链接:", err)
	}
	if _, ok := ResolveCleanTarget(project, "link/bar"); ok {
		t.Fatal("指向项目外的符号链接子路径必须被拒绝")
	}
	target, ok := ResolveCleanTarget(project, "apps/web/node_modules")
	if !ok {
		t.Fatal("普通子路径应被允许")
	}
	resolvedProject, err := filepath.EvalSymlinks(project)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(target, resolvedProject+string(filepath.Separator)) {
		t.Fatalf("目标未落在项目内: %q", target)
	}
}

// TestValidateRootRejectsSymlinkToFilesystemRoot 确认符号链接不能绕过根边界校验。
func TestValidateRootRejectsSymlinkToFilesystemRoot(t *testing.T) {
	root := filepath.VolumeName(t.TempDir()) + string(filepath.Separator)
	link := filepath.Join(t.TempDir(), "root-link")
	if err := os.Symlink(root, link); err != nil {
		t.Skip("当前环境不支持创建符号链接:", err)
	}
	if err := ValidateRoot(link); !errors.Is(err, ErrRootFilesystem) {
		t.Fatalf("指向文件系统根的符号链接应被拒绝，得到 %v", err)
	}
}
