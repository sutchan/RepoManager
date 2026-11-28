// internal/platform/dirs.go v2.3.11

package platform

import (
	"os"
	"path/filepath"
	"strings"
)

// IsSafeCleanDirName 校验自定义清理目录是否为「项目目录内的单个普通目录名」。
// 拒绝路径穿越、绝对路径、盘符、分隔符、空白与 Windows 保留设备名，
// 保证 filepath.Join(project, name) 不会指向 project 之外。
func IsSafeCleanDirName(name string) bool {
	if name == "" || name != strings.TrimSpace(name) {
		return false
	}
	if name == "." || name == ".." || len(name) > 255 {
		return false
	}
	if strings.ContainsAny(name, `/\`) || strings.Contains(name, ":") {
		return false
	}
	if filepath.IsAbs(name) || filepath.Base(name) != name {
		return false
	}
	if isReservedDeviceName(name) {
		return false
	}
	// 控制字符会造成跨平台路径解析异常，直接拒绝。
	for _, r := range name {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

// IsSafeCleanDirPath 校验受限相对子路径（S3）：斜杠 / 反斜杠分隔的每一段都要满足
// IsSafeCleanDirName，拒绝绝对路径、盘符与以分隔符开头的输入，杜绝路径穿越。
// 单级目录名同样合法，兼容 v2.3.10 之前的行为。
func IsSafeCleanDirPath(name string) bool {
	if name == "" || strings.HasPrefix(name, "/") || strings.HasPrefix(name, `\`) || strings.Contains(name, ":") {
		return false
	}
	// 连续分隔符会产生空段，明确拒绝以免绕过逐段校验。
	if strings.Contains(name, "//") || strings.Contains(name, `\\`) {
		return false
	}
	parts := strings.FieldsFunc(name, func(r rune) bool { return r == '/' || r == '\\' })
	if len(parts) == 0 {
		return false
	}
	for _, part := range parts {
		if !IsSafeCleanDirName(part) {
			return false
		}
	}
	return true
}

// ResolveCleanTarget 把受控相对路径解析为项目内的绝对清理目标。
// 逐段拒绝符号链接，防止中间链接把 RemoveAll 导向项目外的任意目录；
// 解析后的最终路径也必须仍位于项目目录内。
// 返回 false 表示该目标不安全，调用方必须跳过。
func ResolveCleanTarget(project, name string) (string, bool) {
	if !IsSafeCleanDirPath(name) {
		return "", false
	}
	base, err := filepath.EvalSymlinks(project)
	if err != nil {
		return "", false
	}
	target := base
	for _, part := range strings.FieldsFunc(name, func(r rune) bool { return r == '/' || r == '\\' }) {
		target = filepath.Join(target, part)
		if info, err := os.Lstat(target); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return "", false
		}
	}
	if resolved, err := filepath.EvalSymlinks(target); err == nil {
		rel, err := filepath.Rel(base, resolved)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", false
		}
	}
	return target, true
}

// isReservedDeviceName 判断 Windows 设备保留名（不区分大小写，忽略扩展名）。
func isReservedDeviceName(name string) bool {
	base := strings.ToLower(name)
	if dot := strings.IndexByte(base, '.'); dot >= 0 {
		base = base[:dot]
	}
	switch base {
	case "con", "prn", "aux", "nul":
		return true
	}
	if len(base) == 4 && (strings.HasPrefix(base, "com") || strings.HasPrefix(base, "lpt")) {
		return base[3] >= '1' && base[3] <= '9'
	}
	return false
}
