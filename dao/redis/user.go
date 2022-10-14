package redis

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"shop-backend/utils/concatstr"
	"time"
)

var (
	verifyCodeKey  = "user:"
	codeLivingTime = time.Minute * 5
)

// SetVerifyCode 将手机验证码存入Redis中，有效期五分钟
func SetVerifyCode(phone, code string) (err error) {
	key := concatstr.ConcatString(verifyCodeKey, phone)
	if err = rdb.Set(key, code, codeLivingTime).Err(); err != nil {
		return
	}
	return
}

// GetVerifyCode 通过手机号获取验证码
func GetVerifyCode(phone string) (code string, err error) {
	key := concatstr.ConcatString(verifyCodeKey, phone)
	if code, err = rdb.Get(key).Result(); err != nil {
		// 验证码已过期或不存在
		if err == redis.Nil {
			err = nil
		}
		zap.L().Error("rdb.Get(key).Result() failed", zap.Error(err))
		return
	}
	return
}
