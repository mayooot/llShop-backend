package gen

import (
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

// ATokenExpireDuration AToken存活时间，1个小时
const ATokenExpireDuration = time.Second

// RTokenExpireDuration AToken存活时间，30天
const RTokenExpireDuration = time.Hour * 24 * 30

// MyIssuer 签发人
var myIssuer = "Bertram Li"

// 秘钥
var mySecret = []byte("shop-backend")

// MyClaims 自定义Token结构体，包含用户ID和用户手机号信息
type MyClaims struct {
	UserID int64  `json:"userId"`
	Phone  string `json:"phone"`
	jwt.StandardClaims
}

// GenToken 生成AToken和RToken
func GenToken(userId int64) (aToken string, rToken string, err error) {
	c := &MyClaims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			// 设置过期时间
			ExpiresAt: time.Now().Add(ATokenExpireDuration).Unix(),
			// 设置签发人
			Issuer: myIssuer,
		},
	}
	// 生成携带用户信息的AToken
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	if err != nil {
		zap.L().Error("jwt create accessToken failed",
			zap.Error(err),
			zap.Int64("uid", userId))
		return
	}
	// 生成不携带用户信息的RToken
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RTokenExpireDuration).Unix(),
		Issuer:    myIssuer,
	}).SignedString(mySecret)
	if err != nil {
		zap.L().Error("jwt create refreshToken failed failed",
			zap.Error(err))
		return
	}
	return
}

// GetSecret 返回自定义秘钥
func GetSecret() func(*jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	}
}
