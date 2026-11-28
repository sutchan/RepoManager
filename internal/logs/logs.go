// internal/logs/logs.go v2.3.14

// Package logs 提供执行日志的文件落盘能力（V4）：在下游 Logger 输出之上
// 追加写入本地日志文件；文件打开 / 写入失败时仅保留下游输出，不中断功能。
package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"repo-manager/internal/engine"
)

// DefaultPath 返回与可执行文件同目录的日志文件路径（repo-manager.log）。
func DefaultPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "repo-manager.log"
	}
	return filepath.Join(filepath.Dir(exe), "repo-manager.log")
}

// maxLogSize 是单个日志文件大小上限（S1：5 MB 后归档为 .log.1）。
// 包级变量便于测试用较小阈值验证轮转边界，生产值固定。
var maxLogSize int64 = 5 << 20

// FileLogger 将日志同时写入文件与下游 Logger（engine.Logger）。
// 每次 Logf 带时间戳前缀；线程安全。
type FileLogger struct {
	mu   sync.Mutex
	f    *os.File
	next engine.Logger
}

// NewFileLogger 以追加模式打开日志文件；打开失败时返回仅含下游输出的 logger。
func NewFileLogger(path string, next engine.Logger) *FileLogger {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		f = nil
	}
	return &FileLogger{f: f, next: next}
}

// rotateIfNeeded 写入前检查大小：超限则关闭并归档为 .log.1（覆盖旧归档），再新建同名文件。
// 归档或重开失败时退回原路径继续写，保证日志降级不中断（对齐包级约定）。
func (l *FileLogger) rotateIfNeeded(incoming int) {
	if l.f == nil {
		return
	}
	info, err := l.f.Stat()
	if err != nil || info.Size()+int64(incoming) <= maxLogSize {
		return
	}
	path := l.f.Name()
	_ = l.f.Close()
	_ = os.Remove(path + ".1")
	if err := os.Rename(path, path+".1"); err != nil {
		l.f, _ = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		return
	}
	l.f, _ = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
}

// Logf 实现 engine.Logger：写文件（带时间戳）并转发下游。
func (l *FileLogger) Logf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f != nil {
		line := fmt.Sprintf("%s %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
		l.rotateIfNeeded(len(line))
		_, _ = fmt.Fprint(l.f, line)
	}
	if l.next != nil {
		l.next.Logf(format, args...)
	}
}

// Close 关闭日志文件（幂等）。
func (l *FileLogger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f != nil {
		_ = l.f.Close()
		l.f = nil
	}
}
