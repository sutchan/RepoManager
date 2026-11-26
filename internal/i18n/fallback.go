// internal/i18n/fallback.go v2.3.11

package i18n

import "reflect"

// fillMissing 用简中字典逐字段填补目标字典的空缺（S13：字段级回退）。
// 仅处理 string 与 [N]string 字段；新增字段类型时须扩展本函数并补测试。
func fillMissing(dst, src any) {
	target := reflect.ValueOf(dst).Elem()
	source := reflect.ValueOf(src)
	for i := 0; i < target.NumField(); i++ {
		field, base := target.Field(i), source.Field(i)
		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				field.SetString(base.String())
			}
		case reflect.Array:
			for j := 0; j < field.Len(); j++ {
				item, fallback := field.Index(j), base.Index(j)
				if item.Kind() == reflect.String && item.String() == "" {
					item.SetString(fallback.String())
				}
			}
		}
	}
}
