package redis

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"shop-backend/utils/concatstr"
	"shop-backend/utils/gen"
	"strconv"
	"time"
)

var (
	userPrefix     = "user:"
	verifyStr      = "verify:"
	tokenStr       = "token:"
	codeLivingTime = time.Minute * 5
)

// SetVerifyCode 将手机验证码存入Redis中，有效期五分钟
func SetVerifyCode(phone, code string) (err error) {
	key := concatstr.ConcatString(userPrefix, verifyStr, phone)
	if err = rdb.Set(key, code, codeLivingTime).Err(); err != nil {
		zap.L().Error("SetAccessToken failed", zap.Error(err))
		return
	}
	return
}

// GetVerifyCode 通过手机号获取验证码
func GetVerifyCode(phone string) (code string, err error) {
	key := concatstr.ConcatString(userPrefix, verifyStr, phone)
	if code, err = rdb.Get(key).Result(); err != nil {
		// 验证码已过期或不存在
		if err == redis.Nil {
			err = nil
		}
		return
	}
	return
}

// SetAccessToken 将AccessToken存入Redis中 K: userID V: aToken
func SetAccessToken(userID int64, aToken string) (err error) {
	key := concatstr.ConcatString(userPrefix, tokenStr, strconv.FormatInt(userID, 10))
	if err = rdb.Set(key, aToken, gen.ATokenExpireDuration).Err(); err != nil {
		zap.L().Error("SetAccessToken failed", zap.Error(err))
		return
	}
	return
}

// GetAccessTokenByUID 通过userID从Redis中获取对应的Access Token
func GetAccessTokenByUID(uid string) (token string, err error) {
	key := concatstr.ConcatString(userPrefix, tokenStr, uid)
	token, err = rdb.Get(key).Result()
	if err != nil {
		return
	}
	return
}

// DelAccessTokenByUID 根据uid删除AccessToken
func DelAccessTokenByUID(idStr string) (err error) {
	key := concatstr.ConcatString(userPrefix, tokenStr, idStr)
	_, err = rdb.Del(key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
	}
	return
}
