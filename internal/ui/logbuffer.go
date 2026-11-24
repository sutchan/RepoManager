// internal/ui/logbuffer.go v2.3.6

package ui

// maxLogLines 是 GUI 输出区保留的最近日志行数；更早的日志仍完整落盘到 repo-manager.log。
// 有界缓冲避免逐行全量拼接导致的 O(n²) 内存拷贝（对齐审查报告 P2-4）。
const maxLogLines = 500

// appendBoundedLog 追加一行日志并裁掉最旧的行，使缓冲不超过 max 行。
// lines 为 nil 时自动创建；max <= 0 时退化为保留最后一行。
func appendBoundedLog(lines []string, msg string, max int) []string {
	lines = append(lines, msg)
	if max <= 0 {
		return lines[len(lines)-1:]
	}
	if len(lines) > max {
		lines = lines[len(lines)-max:]
	}
	return lines
}
