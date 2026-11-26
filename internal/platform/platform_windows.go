//go:build windows

// internal/platform/platform_windows.go v2.3.12

package platform

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// windowsPlatform Windows 实现：attrib 隐藏属性、setx/reg 环境变量持久化。
type windowsPlatform struct{}

// New 返回当前平台的 Platform 实现（由 build tag 分发）。
func New() Platform { return windowsPlatform{} }

func (windowsPlatform) DefaultRepoRoot() string { return `E:\Github` }

func (windowsPlatform) ListProjectDirs(repoRoot string) ([]string, error) {
	return listProjectDirs(repoRoot)
}

func (windowsPlatform) RemoveAll(path string) error  { return os.RemoveAll(path) }
func (windowsPlatform) RemoveFile(path string) error { return os.Remove(path) }

// SetHidden 对应 v1.x `attrib +h +s +r`。
func (windowsPlatform) SetHidden(path string) error {
	return run("attrib", "+h", "+s", "+r", path)
}

// ClearHidden 对应 v1.x `attrib -h -s -r`（功能 4 必须先清除属性再删除）。
func (windowsPlatform) ClearHidden(path string) error {
	return run("attrib", "-h", "-s", "-r", path)
}

func (windowsPlatform) TempDir() string { return os.TempDir() }

// SetUserEnv 直写 HKCU\Environment（S2），规避 setx 对长值的截断；失败时回退 setx。
func (windowsPlatform) SetUserEnv(key, value string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err == nil {
		defer k.Close()
		if err := k.SetStringValue(key, value); err == nil {
			return nil
		}
	}
	return run("setx", key, value)
}

// RemoveUserEnv 直删 HKCU\Environment 值（S2）；不存在视为成功，失败时回退旧 setx/reg 路径。
func (windowsPlatform) RemoveUserEnv(key string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err == nil {
		defer k.Close()
		if err := k.DeleteValue(key); err == nil || errors.Is(err, registry.ErrNotExist) {
			return nil
		}
	}
	clearErr := run("setx", key, "")
	delErr := run("reg", "delete", `HKCU\Environment`, "/v", key, "/f")
	if delErr != nil {
		return errors.Join(clearErr, delErr)
	}
	return nil
}

func (windowsPlatform) GlobalNodeModulesPath() (string, error) {
	return globalNodeModulesPath()
}

// run 执行外部命令；失败时附带输出便于定位。
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s %s: %v: %s",
			ErrCommand, name, strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return nil
}
