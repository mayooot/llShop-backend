package check

import (
	"regexp"
)

// VerifyMobileFormat 校验手机号格式是否正确
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// VerifyEmailFormat 校验电子邮箱格式是否正确
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyUsernameFormat 校验用户名长度
func VerifyUsernameFormat(username string) bool {
	if len(username) >= 5 && len(username) <= 100 {
		return true
	}
	return false
}
