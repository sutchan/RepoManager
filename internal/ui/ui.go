//go:build !nogui

// internal/ui/ui.go v2.3.12

// Package ui 实现 Fyne 桌面 GUI（对齐 docs/DESIGN-SYSTEM.md §3.0）。
// 本包需要 C 工具链编译（Fyne 依赖 OpenGL）；无 C 工具链环境用 `go build -tags nogui` 校验。
package ui

import (
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"repo-manager/internal/config"
	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// langOptions 语言切换顺序：简中 → 繁中 → 英文（覆盖全部支持语言）。
var langOptions = []struct {
	label string
	lang  i18n.Lang
}{
	{"简体中文", i18n.Zh},
	{"繁體中文", i18n.Hant},
	{"English", i18n.En},
}

// Run 启动 Fyne GUI（阻塞直到窗口关闭）。
func Run(cfg config.Config) {
	a := app.NewWithID("repo-manager")
	a.Settings().SetTheme(newTheme())
	w := a.NewWindow("Repo Manager")
	w.Resize(fyne.NewSize(760, 600))

	p := platform.New()
	cfgPath := config.Path()
	if cfg.RepoRoot == "" {
		cfg.RepoRoot = p.DefaultRepoRoot()
	}
	var running atomic.Bool // 跨 render 复用，语言切换重绘时也阻止并发执行

	// render 按当前语言重建窗口内容（语言切换即整体重绘）。
	var render func()
	render = func() {
		lang := i18n.Normalize(cfg.Lang)
		dict := i18n.T(lang)
		w.SetTitle(dict.AppTitle)

		output := widget.NewMultiLineEntry()
		output.Wrapping = fyne.TextWrapWord
		output.Disable()
		progress := newProgressView()
		cancelBtn := widget.NewButton(dict.Cancel, nil)
		cancelBtn.Disable()
		logger := &guiLogger{out: output, mu: &sync.Mutex{}, progress: progress}

		dirEntry := widget.NewEntry()
		dirEntry.SetText(cfg.RepoRoot)

		langSelect := newLangSelect(dict, cfg.Lang, func(next i18n.Lang) {
			cfg.Lang = string(next)
			if err := config.Save(cfgPath, cfg); err != nil {
				dialog.ShowError(err, w)
				return
			}
			render()
		})

		browse := widget.NewButton(dict.ChooseDir, func() {
			dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
				if err != nil || uri == nil {
					return
				}
				dirEntry.SetText(uri.Path())
			}, w)
		})

		saveBtn := widget.NewButton(dict.Save, func() {
			root, ok := resolveRoot(w, dirEntry.Text, lang)
			if !ok {
				return
			}
			cfg.RepoRoot = root
			if err := config.Save(cfgPath, cfg); err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation(dict.Ok, dict.Saved, w)
		})

		aboutBtn := widget.NewButton(dict.About, func() {
			dialog.ShowInformation(dict.About, dict.AboutText, w)
		})

		// 预演 / 递归为会话级选项，不写入 config.ini（对齐 OPENSPEC §4.7 / §3.0）。
		dryRunCheck := widget.NewCheck(dict.DryRunLabel, nil)
		dryRunCheck.SetChecked(false)
		recursiveCheck := widget.NewCheck(dict.RecursiveLabel, nil)
		recursiveCheck.SetChecked(cfg.Recursive)

		runner := &functionRunner{
			w: w, p: p, cfg: &cfg, cfgPath: cfgPath,
			dirEntry: dirEntry, dryRunCheck: dryRunCheck, recursiveCheck: recursiveCheck,
			progress: progress, logger: logger, cancelBtn: cancelBtn, langSelect: langSelect, running: &running,
			lang: lang, dict: dict, d: i18n.Logs(lang),
		}
		cancelBtn.OnTapped = runner.requestCancel
		fnButtons := runner.newFnButtons()

		top := container.NewVBox(
			container.NewHBox(widget.NewLabel(dict.TargetDir), dirEntry, browse, saveBtn),
			container.NewHBox(widget.NewLabel(dict.LangLabel), langSelect, aboutBtn),
			container.NewHBox(widget.NewLabel(dict.RunHint), dryRunCheck, recursiveCheck),
			container.NewVBox(progress.label, progress.bar, cancelBtn),
		)
		// 页脚区（Border bottom）：版本号 + build 编号的唯一展示位置，页眉不放版本号。
		w.SetContent(container.NewBorder(top,
			widget.NewLabel(footerText()), nil, nil,
			container.NewVScroll(container.NewVBox(fnButtons,
				widget.NewLabel(dict.OutputTitle), output))),
		)
	}
	render()
	w.ShowAndRun()
}

// newLangSelect 构造语言下拉框；选中未变化的语言时不触发回调。
func newLangSelect(dict i18n.Dict, current string, onChange func(i18n.Lang)) *widget.Select {
	labels := make([]string, len(langOptions))
	for i, option := range langOptions {
		labels[i] = option.label
	}
	sel := widget.NewSelect(labels, func(s string) {
		for _, option := range langOptions {
			if option.label == s {
				onChange(option.lang)
				return
			}
		}
	})
	for _, option := range langOptions {
		if option.lang == i18n.Normalize(current) {
			sel.SetSelected(option.label)
			break
		}
	}
	sel.PlaceHolder = dict.LangLabel
	return sel
}
