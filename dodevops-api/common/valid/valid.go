package valid

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func ErrorToText(e validator.ValidationErrors) string {
	if len(e) == 0 {
		return ""
	}
	var sb strings.Builder
	first := true
	for _, v := range e {
		if !first {
			sb.WriteString(";")
		}
		switch v.Tag() {
		case "required":
			sb.WriteString(v.Field() + ":字段是必填项")
		case "email":
			sb.WriteString(v.Field() + ":邮箱格式不正确")
		case "gte":
			sb.WriteString(v.Field() + ":必须 ≥ " + v.Param())
		case "lte":
			sb.WriteString(v.Field() + ":必须 ≤ " + v.Param())
		default:
			sb.WriteString(e.Error()) // 默认返回原始错误信息
		}
		first = false
	}
	return sb.String()
}
