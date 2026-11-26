// internal/engine/install_timeout_test.go v2.3.12

package engine

import (
	"context"
	"strings"
	"testing"
	"time"

	"repo-manager/internal/i18n"
)

// TestNpmInstallTimeoutConfig 校验超时默认 10 min、环境变量覆盖与非法值回退（S11）。
func TestNpmInstallTimeoutConfig(t *testing.T) {
	t.Setenv(npmTimeoutEnv, "")
	if got := npmInstallTimeout(); got != installTimeout {
		t.Fatalf("空值应回退默认超时，实际: %s", got)
	}
	if installTimeout != defaultNpmTimeout {
		t.Fatalf("默认超时应为 10 min，实际: %s", installTimeout)
	}
	t.Setenv(npmTimeoutEnv, "25")
	if got := npmInstallTimeout(); got != 25*time.Second {
		t.Fatalf("环境变量应生效，实际: %s", got)
	}
	t.Setenv(npmTimeoutEnv, "abc")
	if got := npmInstallTimeout(); got != installTimeout {
		t.Fatalf("非法值应回退默认超时，实际: %s", got)
	}
	t.Setenv(npmTimeoutEnv, "-1")
	if got := npmInstallTimeout(); got != installTimeout {
		t.Fatalf("非正数应回退默认超时，实际: %s", got)
	}
}

// TestRunNpmInstallTimeoutViaEnv 校验超时后返回带「超时」标识的错误（S11）。
func TestRunNpmInstallTimeoutViaEnv(t *testing.T) {
	t.Setenv(npmTimeoutEnv, "1")
	saved := npmExec
	t.Cleanup(func() { npmExec = saved })
	npmExec = func(ctx context.Context, _ string, _ ...string) ([]byte, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	_, err := runNpmInstall([]string{"react"})
	if err == nil {
		t.Fatal("超时后应返回错误")
	}
	if !strings.Contains(err.Error(), "超时") {
		t.Fatalf("错误信息应标注超时，实际: %v", err)
	}
}

// TestInstallGlobalDepsTimeoutPropagates 校验安装组超时被汇总进失败组（S10 + S11）。
func TestInstallGlobalDepsTimeoutPropagates(t *testing.T) {
	t.Setenv(npmTimeoutEnv, "1")
	savedRunner, savedExec := npmRunner, npmExec
	t.Cleanup(func() { npmRunner, npmExec = savedRunner, savedExec })
	npmExec = func(ctx context.Context, _ string, _ ...string) ([]byte, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	npmRunner = runNpmInstallContext

	var log captureLogger
	err := InstallGlobalDeps(i18n.Zh, Options{}, &log)
	if err == nil {
		t.Fatal("安装组超时汇总后应返回错误")
	}
	if !strings.Contains(err.Error(), "以下分组安装失败") {
		t.Fatalf("错误信息应含失败组汇总，实际: %v", err)
	}
}
