package util

import (
	"regexp"
	"strings"
	"unicode"
)

// GenerateAppCode 根据应用名称生成应用编码
func GenerateAppCode(name string) string {
	if name == "" {
		return ""
	}

	// 1. 移除特殊字符，只保留字母、数字、中文、空格和连字符
	reg := regexp.MustCompile(`[^\p{L}\p{N}\s\-_]`)
	cleaned := reg.ReplaceAllString(name, "")

	// 2. 处理中文字符转拼音首字母（简单映射）
	code := processChinese(cleaned)

	// 3. 转换为小写
	code = strings.ToLower(code)

	// 4. 替换空格和下划线为连字符
	code = regexp.MustCompile(`[\s_]+`).ReplaceAllString(code, "-")

	// 5. 移除多余的连字符
	code = regexp.MustCompile(`-+`).ReplaceAllString(code, "-")

	// 6. 移除首尾连字符
	code = strings.Trim(code, "-")

	// 7. 限制长度，最长32个字符
	if len(code) > 32 {
		code = code[:32]
		// 确保不以连字符结尾
		code = strings.TrimRight(code, "-")
	}

	// 8. 如果为空，返回默认值
	if code == "" {
		code = "app"
	}

	return code
}

// processChinese 处理中文字符，转换为拼音首字母或常见缩写
func processChinese(text string) string {
	// 常见中文词汇映射表
	chineseMap := map[string]string{
		"管理":   "mgmt",
		"系统":   "sys",
		"平台":   "platform",
		"服务":   "service",
		"应用":   "app",
		"接口":   "api",
		"网关":   "gateway",
		"中心":   "center",
		"后台":   "admin",
		"前端":   "frontend",
		"后端":   "backend",
		"数据库": "db",
		"缓存":   "cache",
		"消息":   "msg",
		"队列":   "queue",
		"监控":   "monitor",
		"日志":   "log",
		"文件":   "file",
		"图片":   "image",
		"视频":   "video",
		"用户":   "user",
		"订单":   "order",
		"商品":   "product",
		"支付":   "pay",
		"财务":   "finance",
		"库存":   "inventory",
		"配置":   "config",
		"权限":   "auth",
		"认证":   "auth",
		"登录":   "login",
		"注册":   "register",
		"搜索":   "search",
		"推荐":   "recommend",
		"统计":   "stats",
		"报表":   "report",
		"分析":   "analysis",
		"测试":   "test",
		"开发":   "dev",
		"生产":   "prod",
		"预发":   "staging",
	}

	result := text

	// 替换常见的中文词汇
	for chinese, english := range chineseMap {
		result = strings.ReplaceAll(result, chinese, english)
	}

	// 处理剩余的中文字符，转换为拼音首字母（简化处理）
	var builder strings.Builder
	for _, char := range result {
		if unicode.Is(unicode.Han, char) {
			// 简单的中文字符处理，这里可以集成更完善的拼音库
			// 暂时跳过单个中文字符
			continue
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

// ValidateAppCode 验证应用编码格式
func ValidateAppCode(code string) bool {
	if code == "" {
		return false
	}

	// 只允许小写字母、数字和连字符
	matched, _ := regexp.MatchString(`^[a-z0-9]([a-z0-9\-]*[a-z0-9])?$`, code)
	return matched && len(code) <= 32
}

// GenerateUniqueAppCode 生成唯一的应用编码，支持冲突处理
func GenerateUniqueAppCode(name string, existsChecker func(string) bool) string {
	baseCode := GenerateAppCode(name)

	// 检查基础编码是否已存在
	if !existsChecker(baseCode) {
		return baseCode
	}

	// 如果存在冲突，添加数字后缀
	for i := 1; i <= 999; i++ {
		candidate := baseCode + "-" + string(rune('0'+i/100)) + string(rune('0'+(i%100)/10)) + string(rune('0'+i%10))
		if len(candidate) > 32 {
			// 如果太长，截断基础部分
			maxBase := 32 - 4 // 为 "-xxx" 预留空间
			if maxBase > 0 {
				candidate = baseCode[:maxBase] + "-" + string(rune('0'+i/100)) + string(rune('0'+(i%100)/10)) + string(rune('0'+i%10))
			} else {
				candidate = "app-" + string(rune('0'+i/100)) + string(rune('0'+(i%100)/10)) + string(rune('0'+i%10))
			}
		}

		if !existsChecker(candidate) {
			return candidate
		}
	}

	// 如果仍然冲突，使用时间戳
	return baseCode + "-" + string(rune('0'+len(baseCode)%10))
}