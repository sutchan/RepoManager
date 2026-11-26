//go:build !nogui

// internal/ui/runner.go v2.3.12

package ui

import (
	"context"
	"fmt"
	"sync"
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
	cancelBtn      *widget.Button
	cancelFn       context.CancelFunc
	cancelMu       sync.Mutex
	running        *atomic.Bool
	lang           i18n.Lang
	dict           i18n.Dict
	d              i18n.LogDict
	// runEngine 允许测试注入可控执行器；生产路径回落到 engine.Run（S4 + S5）。
	runEngine func(context.Context, int, string, engine.Options, engine.Logger) error
}

// run 执行功能：危险操作先二次确认，预演模式跳过；执行期间互斥并禁用全部功能按钮。
func (r *functionRunner) run(fn int) {
	root, ok := resolveRoot(r.w, r.dirEntry.Text, r.lang)
	if !ok {
		return
	}
	r.cfg.RepoRoot = root
	opts := engine.Options{
		DryRun:         r.dryRunCheck.Checked,
		Recursive:      r.recursiveCheck.Checked,
		ExtraCleanDirs: r.cfg.ExtraCleanDirs,
		ScanDepth:      r.cfg.ScanDepth,
	}
	do := func(confirmed bool) {
		if !confirmed || !r.running.CompareAndSwap(false, true) {
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		r.setCancel(cancel)
		r.setButtonsEnabled(false)
		r.setCancelEnabled(true)
		r.progress.Reset()
		go func() {
			defer r.clearCancel()
			defer r.running.Store(false)
			defer r.setButtonsEnabled(true)
			defer r.progress.Finish()
			// V4：执行日志同时追加写入 repo-manager.log。
			fileLog := logs.NewFileLogger(logs.DefaultPath(), r.logger)
			defer fileLog.Close()
			if err := r.execute(ctx, fn, root, opts, fileLog); err != nil {
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

// execute 调用 engine；测试可注入 runEngine 以观察取消信号（S4 + S5）。
func (r *functionRunner) execute(ctx context.Context, fn int, root string, opts engine.Options, logger engine.Logger) error {
	if r.runEngine != nil {
		return r.runEngine(ctx, fn, root, opts, logger)
	}
	opts.Context = ctx
	return engine.Run(r.p, fn, root, r.lang, opts, logger)
}

// setCancel 保存当前执行的取消函数。
func (r *functionRunner) setCancel(cancel context.CancelFunc) {
	r.cancelMu.Lock()
	r.cancelFn = cancel
	r.cancelMu.Unlock()
}

// requestCancel 由「取消执行」按钮触发，取消当前上下文。
func (r *functionRunner) requestCancel() {
	r.cancelMu.Lock()
	cancel := r.cancelFn
	r.cancelMu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// clearCancel 结束执行后释放取消函数并禁用按钮。
func (r *functionRunner) clearCancel() {
	r.cancelMu.Lock()
	r.cancelFn = nil
	r.cancelMu.Unlock()
	r.setCancelEnabled(false)
}

// setCancelEnabled 在 UI 线程切换「取消执行」按钮可用状态。
func (r *functionRunner) setCancelEnabled(enabled bool) {
	if r.cancelBtn == nil {
		return
	}
	fyne.Do(func() {
		if enabled {
			r.cancelBtn.Enable()
		} else {
			r.cancelBtn.Disable()
		}
	})
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
