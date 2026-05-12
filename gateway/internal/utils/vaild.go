package utils

import "regexp"

func IsValidEmail(email string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return reg.MatchString(email)
}

// IsValidPhone 验证手机号格式：支持为空（可选填），非空时需符合手机号规则
func IsValidPhone(phone string) bool {
	// 1. 先判断是否为空：空字符串直接返回有效（支持可选填）
	if phone == "" {
		return true
	}

	// 2. 非空时，用正则验证手机号格式（原逻辑不变）
	// 正则规则：以1开头，第二位3-9，后续9位数字（共11位）
	// 正则表达式
	// ^1[3-9]\d{9}$ 解释：
	// ^1          ：以1开头
	// [3-9]       ：第二位数字是3-9
	// \d{9}       ：后面跟9个数字
	// $           ：结束符
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(phone)
}
