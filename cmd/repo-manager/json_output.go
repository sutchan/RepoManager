// cmd/repo-manager/json_output.go v2.3.11

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"repo-manager/internal/engine"
)

// memLogger 收集执行日志供 -json 模式结构化输出（S7）。
type memLogger struct{ msgs []string }

func (m *memLogger) Logf(format string, args ...any) {
	m.msgs = append(m.msgs, fmt.Sprintf(format, args...))
}

// teeLogger 把同一条日志同时写入文件 logger 与内存收集器。
type teeLogger struct{ primary, secondary engine.Logger }

func (t teeLogger) Logf(format string, args ...any) {
	t.primary.Logf(format, args...)
	t.secondary.Logf(format, args...)
}

// jsonResult 是 -json 模式的标准输出契约；退出码仍沿用 0 / 1 定义。
type jsonResult struct {
	Function int      `json:"function"`
	Root     string   `json:"root"`
	OK       bool     `json:"ok"`
	Messages []string `json:"messages"`
	Error    string   `json:"error,omitempty"`
}

// emitJSONResult 向 stdout 输出缩进 JSON，不改写退出码（由调用方处理）。
func emitJSONResult(fn int, root string, msgs []string, err error) {
	r := jsonResult{Function: fn, Root: root, OK: err == nil, Messages: msgs}
	if err != nil {
		r.Error = err.Error()
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(r)
}
