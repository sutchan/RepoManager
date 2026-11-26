// internal/i18n/parity_test.go v2.3.11

package i18n

import (
	"reflect"
	"testing"
)

// TestFillMissing 验证 string 与数组字段的空缺由简中逐字段补齐（S13）。
func TestFillMissing(t *testing.T) {
	dst := Dict{AppTitle: "自定义", FuncNames: [7]string{"一"}}
	fillMissing(&dst, dicts[Zh])
	if dst.AppTitle != "自定义" {
		t.Fatalf("非空字段不应被覆盖: %q", dst.AppTitle)
	}
	if dst.TargetDir != dicts[Zh].TargetDir {
		t.Fatalf("空字段应由简中补齐: %q", dst.TargetDir)
	}
	if dst.FuncNames[0] != "一" || dst.FuncNames[6] != dicts[Zh].FuncNames[6] {
		t.Fatalf("数组字段应逐项补齐: %v", dst.FuncNames)
	}
}

// TestAllLocalesComplete 校验三语 Dict / LogDict 所有字段均非空（字段级回退生效）。
func TestAllLocalesComplete(t *testing.T) {
	for _, lang := range []Lang{Zh, Hant, En} {
		assertNoEmptyFields(t, string(lang)+"/Dict", T(lang))
		assertNoEmptyFields(t, string(lang)+"/LogDict", Logs(lang))
	}
}

func assertNoEmptyFields(t *testing.T, label string, value any) {
	t.Helper()
	rv := reflect.ValueOf(value)
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				t.Errorf("%s 字段 %s 为空", label, rv.Type().Field(i).Name)
			}
		case reflect.Array:
			for j := 0; j < field.Len(); j++ {
				if field.Index(j).Kind() == reflect.String && field.Index(j).String() == "" {
					t.Errorf("%s 字段 %s[%d] 为空", label, rv.Type().Field(i).Name, j)
				}
			}
		}
	}
}
