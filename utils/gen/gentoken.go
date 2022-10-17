package gen

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

// ATokenExpireDuration AToken存活时间，1个小时
const ATokenExpireDuration = time.Hour * 24

// RTokenExpireDuration AToken存活时间，30天
const RTokenExpireDuration = time.Hour * 24 * 30

// MyIssuer 签发人
var myIssuer = "Bertram Li"

// 秘钥
var mySecret = []byte("shop-backend")

// MyClaims 自定义Token结构体，包含用户ID和用户手机号信息
type MyClaims struct {
	UserId int64  `json:"userId"`
	Phone  string `json:"phone"`
	jwt.StandardClaims
}

// GenToken 生成AToken和RToken
func GenToken(userId int64, phone string) (aToken string, rToken string, err error) {
	c := &MyClaims{
		UserId: userId,
		Phone:  phone,
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
		zap.L().Error("jwt create AccessToken failed",
			zap.Error(err),
			zap.Int64("uid", userId))
		zap.String("phone", phone)
		return
	}
	// 生成不携带用户信息的RToken
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RTokenExpireDuration).Unix(),
		Issuer:    myIssuer,
	}).SignedString(mySecret)
	if err != nil {
		zap.L().Error("jwt create AccessToken failed",
			zap.Error(err),
			zap.Int64("uid", userId))
		zap.String("phone", phone)
		return
	}
	return
}

// ParseToken 解析Token
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 初始化claims
	claims = new(MyClaims)
	var token *jwt.Token
	// 解析token
	token, err = jwt.ParseWithClaims(tokenString, claims, GetSecret())

	if err != nil {
		return
	}
	// 校验token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// GetSecret 返回自定义秘钥
func GetSecret() func(*jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	}
}

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken string, err error) {
	if _, err = jwt.Parse(rToken, GetSecret()); err != nil {
		// RefreshToken已经失效或错误，那么就不能用来刷新AccessToken
		return
	}
	// 解析AccessToken
	var claims *MyClaims
	claims, err = ParseToken(aToken)
	v, ok := err.(*jwt.ValidationError)
	if !ok {
		// 如果不是ValidationError类型的错误
		return
	}
	if v.Errors == jwt.ValidationErrorExpired {
		// 如果AccessToke是过期错误，且RefreshToken没有过期。创建一个新的AccessToken返回
		newAToken, _, err = GenToken(claims.UserId, claims.Phone)
		return
	}
	return
}
