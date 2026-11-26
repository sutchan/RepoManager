// internal/engine/context_test.go v2.3.12

package engine

import (
	"context"
	"errors"
	"testing"

	"repo-manager/internal/i18n"
)

// TestRunStopsOnCancelledContext 校验 GUI 传入的取消上下文在执行前即阻断副作用（S5）。
func TestRunStopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := newStub(t)
	var log captureLogger
	err := Run(s, 1, s.DefaultRepoRoot(), i18n.Zh, Options{Context: ctx}, &log)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("取消后应返回 context.Canceled，实际: %v", err)
	}
	if len(s.envs) != 0 {
		t.Fatalf("取消后不应写入环境变量，实际: %v", s.envs)
	}
	if len(log.lines) != 0 {
		t.Fatalf("取消后不应产生执行日志，实际: %v", log.lines)
	}
}

// TestInstallGlobalDepsStopsOnCancelledContext 校验取消后不再启动 npm 分组（S5）。
func TestInstallGlobalDepsStopsOnCancelledContext(t *testing.T) {
	saved := npmRunner
	t.Cleanup(func() { npmRunner = saved })
	called := false
	npmRunner = func(context.Context, []string) (string, error) {
		called = true
		return "", nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var log captureLogger
	err := InstallGlobalDeps(i18n.Zh, Options{Context: ctx}, &log)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("取消后应返回 context.Canceled，实际: %v", err)
	}
	if called {
		t.Fatal("取消后不应启动 npm 安装")
	}
}
