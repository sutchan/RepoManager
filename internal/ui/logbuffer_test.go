// internal/ui/logbuffer_test.go v2.3.6

package ui

import (
	"strconv"
	"testing"
)

func TestAppendBoundedLog(t *testing.T) {
	var lines []string
	for i := 0; i < maxLogLines+10; i++ {
		lines = appendBoundedLog(lines, strconv.Itoa(i), maxLogLines)
	}
	if len(lines) != maxLogLines {
		t.Fatalf("缓冲应被限制为 %d 行，得到 %d", maxLogLines, len(lines))
	}
	if lines[0] != "10" || lines[len(lines)-1] != strconv.Itoa(maxLogLines+9) {
		t.Fatalf("应保留最近 %d 行，得到首=%q 末=%q", maxLogLines, lines[0], lines[len(lines)-1])
	}
	if got := appendBoundedLog(nil, "only", 0); len(got) != 1 || got[0] != "only" {
		t.Fatalf("max<=0 应保留最后一行，得到 %v", got)
	}
}
