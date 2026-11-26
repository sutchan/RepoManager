// internal/logs/logs_test.go v2.3.3

package logs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFileLogger V4：日志应同时落盘（带时间戳）并转发下游。
func TestFileLogger(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "repo-manager.log")
	var forwarded int
	l := NewFileLogger(path, &countLogger{&forwarded})
	l.Logf("step %d", 1)
	l.Logf("clean %s", "proj-a")
	l.Close()
	l.Logf("after close") // 关闭后不写文件，仍应转发

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if strings.Count(content, "\n") != 2 {
		t.Fatalf("文件应恰有 2 条日志，得到 %q", content)
	}
	if !strings.Contains(content, "step 1") || !strings.Contains(content, "clean proj-a") {
		t.Fatalf("文件应包含日志内容，得到 %q", content)
	}
	if strings.Contains(content, "after close") {
		t.Fatal("关闭后不应再写文件")
	}
	if forwarded != 3 {
		t.Fatalf("下游应收到全部 3 条日志，得到 %d", forwarded)
	}
}

// TestNewFileLoggerInvalidPath 打开失败时应降级为仅下游输出，不 panic。
func TestNewFileLoggerInvalidPath(t *testing.T) {
	n := 0
	l := NewFileLogger(filepath.Join(t.TempDir(), "no-such-dir", "x.log"), &countLogger{&n})
	l.Logf("hello")
	l.Close()
	if n != 1 {
		t.Fatalf("降级时下游仍应收到日志，得到 %d", n)
	}
}

// TestLogRotation 验证超过阈值后归档为 .log.1 并继续写入新文件（S1）。
func TestLogRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "repo-manager.log")
	old := maxLogSize
	maxLogSize = 80
	defer func() { maxLogSize = old }()

	n := 0
	l := NewFileLogger(path, &countLogger{&n})
	for i := 0; i < 20; i++ {
		l.Logf("line-%02d-%s", i, strings.Repeat("x", 30))
	}
	l.Close()

	if _, err := os.Stat(path + ".1"); err != nil {
		t.Fatalf("轮转归档 .log.1 未生成: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "line-00") {
		t.Fatalf("主日志仍包含轮转前内容: %q", string(data))
	}
	if n != 20 {
		t.Fatalf("下游应收到全部 20 条日志，得到 %d", n)
	}
}

type countLogger struct{ n *int }

func (c *countLogger) Logf(format string, args ...any) { *c.n++ }
