// internal/engine/failures_test.go v2.3.12

package engine

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// removeFailPlatform 让删除类操作失败，用于验证 S10 的失败汇总行为。
type removeFailPlatform struct {
	*stubPlatform
}

func (*removeFailPlatform) RemoveAll(string) error { return errors.New("permission denied") }

// envFailPlatform 让环境变量写入失败，验证危险操作不会静默成功。
type envFailPlatform struct {
	*stubPlatform
}

func (*envFailPlatform) SetUserEnv(string, string) error { return errors.New("registry locked") }

// TestFailureCollectorEmpty 校验无失败项时不返回错误。
func TestFailureCollectorEmpty(t *testing.T) {
	var failed failureCollector
	if err := failed.err(i18n.Zh); err != nil {
		t.Fatalf("无失败项时应返回 nil，实际: %v", err)
	}
}

// TestFailureCollectorFormat 校验失败项按 fmt 约定格式化并可汇总。
func TestFailureCollectorFormat(t *testing.T) {
	var failed failureCollector
	failed.add("%s: %v", "a", errors.New("x"))
	failed.add("%s: %v", "b", errors.New("y"))
	err := failed.err(i18n.Zh)
	if err == nil {
		t.Fatal("应返回汇总错误")
	}
	msg := err.Error()
	for _, want := range []string{"部分操作失败（2 项）", "a: x", "b: y"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("汇总信息缺少 %q，实际: %s", want, msg)
		}
	}
}

// TestCleanAllCachesFailureSummary 校验功能 6 局部删除失败时返回非零码且汇总含路径（S10）。
func TestCleanAllCachesFailureSummary(t *testing.T) {
	s := newStub(t)
	proj := filepath.Join(s.DefaultRepoRoot(), "proj-a", "node_modules")
	newProject(t, s.DefaultRepoRoot(), "proj-a", map[string][]byte{
		"package.json":        []byte("{}"),
		"node_modules/lib.js": []byte("x"),
	})
	fp := &removeFailPlatform{stubPlatform: s}
	var log captureLogger
	err := CleanAllCaches(fp, s.DefaultRepoRoot(), i18n.Zh, Options{}, &log)
	if err == nil {
		t.Fatal("局部删除失败时应返回汇总错误（S10）")
	}
	if !strings.Contains(err.Error(), "部分操作失败") {
		t.Fatalf("错误信息应含失败汇总，实际: %v", err)
	}
	if !strings.Contains(err.Error(), proj) {
		t.Fatalf("错误信息应含失败路径 %q，实际: %v", proj, err)
	}
}

// TestEnableGlobalNodeModulesFailureSummary 校验功能 1 环境变量写入失败时返回非零码（S10）。
func TestEnableGlobalNodeModulesFailureSummary(t *testing.T) {
	s := newStub(t)
	newProject(t, s.DefaultRepoRoot(), "proj-a", map[string][]byte{"package.json": []byte("{}")})
	fp := &envFailPlatform{stubPlatform: s}
	var log captureLogger
	err := EnableGlobalNodeModules(fp, s.DefaultRepoRoot(), i18n.Zh, Options{}, &log)
	if err == nil {
		t.Fatal("环境变量写入失败时应返回汇总错误（S10）")
	}
	if !strings.Contains(err.Error(), "NODE_PATH") {
		t.Fatalf("错误信息应含 NODE_PATH 上下文，实际: %v", err)
	}
}

// TestDescribeErrorLocalizesSentinels 校验平台哨兵错误按语言翻译（S17）。
func TestDescribeErrorLocalizesSentinels(t *testing.T) {
	rootErr := platform.ValidateRoot(filepath.Join(t.TempDir(), "missing"))
	zh := DescribeError(i18n.Zh, rootErr)
	if !strings.Contains(zh, "仓库根目录不存在") {
		t.Fatalf("简中应翻译路径不存在，实际: %s", zh)
	}
	en := DescribeError(i18n.En, rootErr)
	if !strings.Contains(en, "Repository root does not exist") {
		t.Fatalf("英文应翻译路径不存在，实际: %s", en)
	}
	if got := DescribeError(i18n.Zh, context.Canceled); !strings.Contains(got, "操作已取消") {
		t.Fatalf("取消错误应本地化，实际: %s", got)
	}
}
