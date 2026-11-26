// cmd/repo-manager/json_output_test.go v2.3.11

package main

import (
	"encoding/json"
	"testing"
)

// TestJSONResultShape 校验 -json 输出契约字段（S7）。
func TestJSONResultShape(t *testing.T) {
	mem := &memLogger{}
	mem.Logf("已处理 %d 项", 3)
	if len(mem.msgs) != 1 || mem.msgs[0] != "已处理 3 项" {
		t.Fatalf("内存日志收集异常: %v", mem.msgs)
	}

	r := jsonResult{Function: 6, Root: `/tmp/repo`, OK: false, Messages: mem.msgs, Error: "boom"}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"function", "root", "ok", "messages", "error"} {
		if _, ok := got[key]; !ok {
			t.Fatalf("JSON 输出缺少字段 %q: %s", key, data)
		}
	}
	if got["ok"] != false || got["function"] != float64(6) {
		t.Fatalf("字段值异常: %s", data)
	}
}

// TestTeeLogger 校验日志同时写入两个下游。
func TestTeeLogger(t *testing.T) {
	a, b := &memLogger{}, &memLogger{}
	teeLogger{primary: a, secondary: b}.Logf("line-%d", 1)
	if len(a.msgs) != 1 || len(b.msgs) != 1 || a.msgs[0] != b.msgs[0] {
		t.Fatalf("双写不一致: a=%v b=%v", a.msgs, b.msgs)
	}
}
