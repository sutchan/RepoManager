// internal/engine/failures.go v2.3.11

package engine

import (
	"fmt"
	"strings"

	"repo-manager/internal/i18n"
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
