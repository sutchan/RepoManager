// internal/engine/install_test.go v2.3.14

package engine

import (
	"strings"
	"testing"
	"time"
)

// TestRunNpmInstallTimeout 验证超时作为明确错误返回（S11）。
func TestRunNpmInstallTimeout(t *testing.T) {
	old := installTimeout
	installTimeout = time.Nanosecond
	defer func() { installTimeout = old }()

	_, err := runNpmInstall([]string{"no-such-package"})
	if err == nil || !strings.Contains(err.Error(), "超时") {
		t.Fatalf("期望明确超时错误，得到 %v", err)
	}
}
