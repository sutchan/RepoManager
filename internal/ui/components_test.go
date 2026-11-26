//go:build !nogui

// internal/ui/components_test.go v2.3.12

package ui

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"repo-manager/internal/config"
	"repo-manager/internal/engine"
	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// TestNewLangSelectComponent 校验语言下拉框三语覆盖与回调（S4）。
func TestNewLangSelectComponent(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	var got i18n.Lang
	sel := newLangSelect(i18n.T(i18n.Zh), string(i18n.Zh), func(lang i18n.Lang) { got = lang })
	if len(sel.Options) != 3 {
		t.Fatalf("语言选项应覆盖三语，实际: %v", sel.Options)
	}
	if sel.Selected != "简体中文" {
		t.Fatalf("初始语言应选中简体中文，实际: %q", sel.Selected)
	}
	sel.SetSelected("English")
	if got != i18n.En {
		t.Fatalf("选择 English 应回调 i18n.En，实际: %q", got)
	}
}

// TestNewFnButtonsComponent 校验功能按钮按当前语言生成且数量正确（S4）。
func TestNewFnButtonsComponent(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	dict := i18n.T(i18n.Zh)
	r := &functionRunner{dict: dict, running: &atomic.Bool{}}
	grid := r.newFnButtons()
	if len(r.buttons) != 7 || len(grid.Objects) != 7 {
		t.Fatalf("应生成 7 个功能按钮，实际 buttons=%d objects=%d", len(r.buttons), len(grid.Objects))
	}
	if want := "[1] " + dict.FuncNames[0]; r.buttons[0].Text != want {
		t.Fatalf("按钮文案异常，want %q got %q", want, r.buttons[0].Text)
	}
}

// TestFunctionRunnerCancellation 校验取消按钮把 context 取消传回 engine（S4 + S5）。
func TestFunctionRunnerCancellation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	root := t.TempDir()
	dirEntry := widget.NewEntry()
	dirEntry.SetText(root)
	progress := newProgressView()
	cancelBtn := widget.NewButton("cancel", nil)
	cancelBtn.Disable()
	logger := &guiLogger{out: widget.NewMultiLineEntry(), mu: &sync.Mutex{}, progress: progress}
	started := make(chan struct{})
	cancelSeen := make(chan error, 1)

	r := &functionRunner{
		w: test.NewWindow(nil), p: platform.New(), cfg: &config.Config{RepoRoot: root},
		dirEntry: dirEntry, dryRunCheck: widget.NewCheck("", nil), recursiveCheck: widget.NewCheck("", nil),
		progress: progress, logger: logger, cancelBtn: cancelBtn, running: &atomic.Bool{},
		lang: i18n.Zh, dict: i18n.T(i18n.Zh), d: i18n.Logs(i18n.Zh),
		runEngine: func(ctx context.Context, _ int, _ string, _ engine.Options, _ engine.Logger) error {
			close(started)
			<-ctx.Done()
			cancelSeen <- ctx.Err()
			return ctx.Err()
		},
	}
	cancelBtn.OnTapped = r.requestCancel
	if !cancelBtn.Disabled() {
		t.Fatal("初始状态取消按钮应禁用")
	}
	r.run(3)

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("功能执行未启动")
	}
	if cancelBtn.Disabled() {
		t.Fatal("执行中取消按钮应启用")
	}
	test.Tap(cancelBtn)
	select {
	case err := <-cancelSeen:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("执行器应收到 context.Canceled，实际: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("取消信号未传到执行器")
	}
	for i := 0; i < 200; i++ {
		if !r.running.Load() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("执行结束后互斥状态未释放")
}
