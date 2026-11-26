//go:build !nogui

// internal/ui/runner.go v2.3.6

package ui

import (
	"fmt"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"repo-manager/internal/config"
	"repo-manager/internal/engine"
	"repo-manager/internal/i18n"
	"repo-manager/internal/logs"
	"repo-manager/internal/platform"
)

// functionRunner 封装一次功能执行所需的 UI 依赖与并发互斥状态。
// running 由 Run 作用域创建并跨 render 复用：语言切换重绘时仍阻止并发危险操作。
type functionRunner struct {
	w              fyne.Window
	p              platform.Platform
	cfg            *config.Config
	cfgPath        string
	dirEntry       *widget.Entry
	dryRunCheck    *widget.Check
	recursiveCheck *widget.Check
	progress       *progressView
	logger         *guiLogger
	buttons        []*widget.Button
	running        *atomic.Bool
	lang           i18n.Lang
	dict           i18n.Dict
	d              i18n.LogDict
}

// run 执行功能：危险操作先二次确认，预演模式跳过；执行期间互斥并禁用全部功能按钮。
func (r *functionRunner) run(fn int) {
	root, ok := resolveRoot(r.w, r.dirEntry.Text)
	if !ok {
		return
	}
	r.cfg.RepoRoot = root
	opts := engine.Options{
		DryRun:         r.dryRunCheck.Checked,
		Recursive:      r.recursiveCheck.Checked,
		ExtraCleanDirs: r.cfg.ExtraCleanDirs,
	}
	do := func(confirmed bool) {
		if !confirmed || !r.running.CompareAndSwap(false, true) {
			return
		}
		r.setButtonsEnabled(false)
		r.progress.Reset()
		go func() {
			defer r.running.Store(false)
			defer r.setButtonsEnabled(true)
			defer r.progress.Finish()
			// V4：执行日志同时追加写入 repo-manager.log。
			fileLog := logs.NewFileLogger(logs.DefaultPath(), r.logger)
			defer fileLog.Close()
			if err := engine.Run(r.p, fn, root, r.lang, opts, fileLog); err != nil {
				fileLog.Logf(r.d.ErrorPrefix, err)
			}
		}()
	}
	if engine.IsDangerous(fn) && !opts.DryRun {
		warning := r.dict.WarningF1
		if fn == 6 {
			warning = r.dict.WarningAll
		} else if fn == 7 {
			warning = r.dict.WarningVite
		}
		dialog.ShowConfirm(warning, r.dict.Confirm, do, r.w)
		return
	}
	do(true)
}

// setButtonsEnabled 在 UI 线程批量切换功能按钮可用状态。
func (r *functionRunner) setButtonsEnabled(enabled bool) {
	fyne.Do(func() {
		for _, b := range r.buttons {
			if enabled {
				b.Enable()
			} else {
				b.Disable()
			}
		}
	})
}

// newFnButtons 按当前语言生成功能按钮，并挂接执行器。
func (r *functionRunner) newFnButtons() *fyne.Container {
	grid := container.NewGridWithColumns(2)
	for i, name := range r.dict.FuncNames {
		fn := i + 1
		btn := widget.NewButton(fmt.Sprintf("[%d] %s", fn, name), func() { r.run(fn) })
		grid.Add(btn)
		r.buttons = append(r.buttons, btn)
	}
	return grid
}
