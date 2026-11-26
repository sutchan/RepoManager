// internal/engine/failures.go v2.3.12

package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// failureCollector 汇总局部失败项（S10）。
// 危险操作逐项执行时先记录失败、继续处理其余项，结束时统一返回非零状态，
// 避免「部分未删除但 CLI 退出 0」的静默失败。
type failureCollector struct {
	items []string
}

// add 记录一条失败；format 与 args 按 fmt.Sprintf 约定。
func (c *failureCollector) add(format string, args ...any) {
	c.items = append(c.items, fmt.Sprintf(format, args...))
}

// err 返回汇总错误；无失败项时返回 nil。
func (c *failureCollector) err(lang i18n.Lang) error {
	if len(c.items) == 0 {
		return nil
	}
	return fmt.Errorf(i18n.Logs(lang).PartialFailSummary, len(c.items), strings.Join(c.items, "; "))
}

// DescribeError 把平台层哨兵错误与 context 错误翻译为当前语言文案（S17）。
// 平台层保持无 i18n 依赖，所有面向用户的错误在 engine 边界统一格式化。
func DescribeError(lang i18n.Lang, err error) string {
	if err == nil {
		return ""
	}
	d := i18n.Logs(lang)
	switch {
	case errors.Is(err, platform.ErrRootFilesystem):
		return fmt.Sprintf(d.RootFilesystemErr, err)
	case errors.Is(err, platform.ErrRootNotExist):
		return fmt.Sprintf(d.RootNotExistErr, err)
	case errors.Is(err, platform.ErrRootNotDir):
		return fmt.Sprintf(d.RootNotDirErr, err)
	case errors.Is(err, platform.ErrNpmRoot):
		return fmt.Sprintf(d.NpmRootErr, err)
	case errors.Is(err, context.Canceled):
		return d.CancelledErr
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Sprintf(d.TimeoutErr, err)
	default:
		return fmt.Sprintf(d.ErrorPrefix, err)
	}
}
