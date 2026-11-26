// internal/engine/install.go v2.3.11

package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"repo-manager/internal/i18n"
)

const (
	// defaultNpmTimeout 功能 5 单次 npm install -g 的默认超时（S11）。
	defaultNpmTimeout = 10 * time.Minute
	// npmTimeoutEnv 覆盖超时的环境变量，单位为秒；非法或非正数回退默认值。
	npmTimeoutEnv = "REPO_MANAGER_NPM_TIMEOUT"
)

// installTimeout 单次 npm 安装的默认超时；测试可直接覆盖（S11）。
var installTimeout = defaultNpmTimeout

// npmRunner 功能 5 的安装执行入口；测试可替换以避免真实联网安装。
var npmRunner = runNpmInstall

// npmExec 是 npm 命令的实际执行器；测试可替换以注入受控 context 行为。
var npmExec = func(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// npmInstallTimeout 返回本次安装超时：优先环境变量（秒），缺省或非法值回退 10 min。
func npmInstallTimeout() time.Duration {
	raw := strings.TrimSpace(os.Getenv(npmTimeoutEnv))
	if raw == "" {
		return installTimeout
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return installTimeout
	}
	return time.Duration(seconds) * time.Second
}

// InstallGlobalDeps 功能 5：分组安装常用全局依赖（对齐 OPENSPEC §3.5）。
// 单组失败不中断后续组；结束后汇总失败组并返回错误（v2.1.0 修复 L4）。
// DryRun 下仅列出各组将要安装的包，不执行 npm。
func InstallGlobalDeps(lang i18n.Lang, opts Options, log Logger) error {
	d := i18n.Logs(lang)
	var failed []string
	for i, group := range npmGroups {
		log.Logf("%s", d.GroupNames[i])
		if opts.DryRun {
			log.Logf(d.DryRunInstall, strings.Join(group, " "))
			continue
		}
		out, err := npmRunner(group)
		if err != nil {
			log.Logf(d.GroupFail, i+1, err)
			if out != "" {
				log.Logf("%s", out)
			}
			failed = append(failed, fmt.Sprintf("%d", i+1))
		}
	}
	if len(failed) > 0 {
		summary := fmt.Sprintf(d.GroupSummary, strings.Join(failed, ", "))
		log.Logf("%s", summary)
		return errors.New(summary)
	}
	return nil
}

// runNpmInstall 执行 `npm install -g <包...>`，超时由 context 控制（S11）。
func runNpmInstall(pkgs []string) (string, error) {
	timeout := npmInstallTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	args := append([]string{"install", "-g"}, pkgs...)
	out, err := npmExec(ctx, "npm", args...)
	if err != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return strings.TrimSpace(string(out)), fmt.Errorf("npm install 超时（%s）: %w", timeout, err)
	}
	return strings.TrimSpace(string(out)), err
}
