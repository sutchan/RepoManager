//go:build !nogui

// internal/ui/root.go v2.3.10

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"repo-manager/internal/platform"
)

// resolveRoot 归一化并校验仓库根目录；校验失败时弹窗提示并返回 ok=false。
// 「保存配置」与「执行功能」共用同一套校验逻辑，避免两处重复实现导致行为漂移。
func resolveRoot(w fyne.Window, raw string) (string, bool) {
	root := platform.NormalizeRoot(raw)
	if err := platform.ValidateRoot(root); err != nil {
		dialog.ShowError(err, w)
		return "", false
	}
	return root, true
}
