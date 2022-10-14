package gen

import (
	"fmt"
	"math/rand"
	"shop-backend/settings"
	"strings"
	"time"
)

// GenVerifyCode 生成width位验证码
func GenVerifyCode() (string, error) {
	width := settings.Conf.UserConfig.Width
	number := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	l := len(number)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%d", number[rand.Intn(l)])
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}
