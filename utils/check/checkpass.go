package check

import (
	"fmt"
	"regexp"
	"shop-backend/settings"
)

// CheckPass 校验密码强度，必须存在特殊字符、大小写字母和数字，校验密码长度
func CheckPass(pass string) error {
	minLength := settings.Conf.UserConfig.MinPassLen
	maxLength := settings.Conf.UserConfig.MaxPassLen
	if len(pass) < minLength {
		return fmt.Errorf("password len is < %d", minLength)
	}
	if len(pass) > maxLength {
		return fmt.Errorf("password len is > %d", maxLength)
	}

	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`

	if b, err := regexp.MatchString(num, pass); !b || err != nil {
		return fmt.Errorf("password need num :%v", err)
	}
	if b, err := regexp.MatchString(a_z, pass); !b || err != nil {
		return fmt.Errorf("password need a_z :%v", err)
	}
	if b, err := regexp.MatchString(A_Z, pass); !b || err != nil {
		return fmt.Errorf("password need A_Z :%v", err)
	}
	if b, err := regexp.MatchString(symbol, pass); !b || err != nil {
		return fmt.Errorf("password need symbol :%v", err)
	}
	return nil
}
