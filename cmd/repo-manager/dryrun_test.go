// cmd/repo-manager/dryrun_test.go v2.3.16

package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"repo-manager/internal/config"
	"repo-manager/internal/engine"
	"repo-manager/internal/i18n"
	"repo-manager/internal/platform"
)

// TestRunCLIDryRunDoesNotWriteConfig 确认预演不回写 config.ini（对齐 OPENSPEC §4.7）。
func TestRunCLIDryRunDoesNotWriteConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.ini")
	original := "REPO_ROOT=old\nLANG=zh\n"
	if err := os.WriteFile(cfgPath, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(dir, "projects")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg.RepoRoot = root
	runCLI(platform.New(), cfgPath, &cfg, 6, root, i18n.Zh, engine.Options{DryRun: true}, false)
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != original {
		t.Fatalf("预演不应修改 config.ini，得到 %q", data)
	}
}

// TestRunCLIJSONStdoutIsPureJSON 确认 -json 模式 stdout 只有单个 JSON 对象，日志不混入。
func TestRunCLIJSONStdoutIsPureJSON(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.ini")
	root := filepath.Join(dir, "projects")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = writer
	defer func() { os.Stdout = original }()
	runCLI(platform.New(), cfgPath, &cfg, 6, root, i18n.Zh, engine.Options{DryRun: true}, true)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	var result jsonResult
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("stdout 应为单个 JSON 文档，得到 %q: %v", out, err)
	}
	if !result.OK || result.Function != 6 || len(result.Messages) == 0 {
		t.Fatalf("JSON 结果不完整: %+v", result)
	}
}

// TestRunCLIFunction5SkipsRootValidation 功能 5 不访问 REPO_ROOT，目录不存在也应正常预演（对齐 OPENSPEC §3.5）。
func TestRunCLIFunction5SkipsRootValidation(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.ini")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	missingRoot := filepath.Join(dir, "does-not-exist")
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = writer
	defer func() { os.Stdout = original }()
	// 若仍校验根目录，runCLI 会 os.Exit(1) 并直接终止测试进程。
	runCLI(platform.New(), cfgPath, &cfg, 5, missingRoot, i18n.Zh, engine.Options{DryRun: true}, true)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	var result jsonResult
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("stdout 应为单个 JSON 文档，得到 %q: %v", out, err)
	}
	if !result.OK || result.Function != 5 || result.Root != missingRoot {
		t.Fatalf("功能 5 预演结果异常: %+v", result)
	}
}
