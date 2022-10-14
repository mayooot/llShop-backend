package check

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"shop-backend/utils/gen"
)

// CheckAToken 解析AToken
func CheckAToken(tokenString string) (claims *gen.MyClaims, err error) {
	// 初始化claims
	claims = new(gen.MyClaims)
	var token *jwt.Token
	// 解析token
	token, err = jwt.ParseWithClaims(tokenString, claims, gen.GetSecret())

	if err != nil {
		// 解析错误
		zap.L().Error("jwt.ParseWithClaims failed",
			zap.Error(err),
			zap.Int64("uid", claims.UserId),
			zap.String("phone", claims.Phone),
		)
		return nil, err
	}
	// 校验token
	if !token.Valid {
		err = errors.New("invalid token")
		return nil, err
	}
	return claims, nil
}
