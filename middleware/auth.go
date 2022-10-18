package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/controller"
	"shop-backend/dao/redis"
	"shop-backend/utils/check"
	"shop-backend/utils/gen"
	"strings"
)

var (
	CtxUserIdKey = "uid"
	CtxAToken    = "atoken"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 用户需要携带AccessToken和RefreshToken
		// Authorization: Bearer AccessToken&RefreshToken
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// 未携带Token
			zap.L().Error("c.Request.Header.Get(\"Authorization\") is nil")
			controller.ResponseError(c, controller.CodeTokenIsEmpty)
			c.Abort()
			return
		}

		// 按照空格进行分割
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 携带的Token格式有误
			zap.L().Error("trings.SplitN(authHeader, \" \", 2) failed", zap.String("authHeader", authHeader))
			controller.ResponseError(c, controller.CodeTokenIsWrongFormat)
			c.Abort()
			return
		}

		// 解析前端传递的两个Token
		tokens := strings.SplitN(parts[1], "&", 2)
		accessToken := tokens[0]
		refreshToken := tokens[1]
		// 解析AccessToken
		mc, err := check.CheckToken(accessToken)
		if err != nil {
			if err == check.ErrorATokenExpired {
				// 如果错误类型为accessToken过期错误，那么需要使用refreshToken协助刷新
				// zap.L().Info("accessToken is expired", zap.Int64("uid", mc.UserID))
				_, err = check.CheckToken(refreshToken)
				if err != nil {
					// 如果解析refreshToken出现错误
					zap.L().Error("parseRefreshToken error", zap.Int64("uid", mc.UserID))
					controller.ResponseError(c, controller.CodeNeedReLogin)
					c.Abort()
					return
				} else {
					var newToken string
					newToken, _, err = gen.GenToken(mc.UserID)
					if err != nil {
						// 生成AccessToken失败
						zap.L().Error("use refreshToken refresh accessToken failed")
						controller.ResponseError(c, controller.CodeTokenRefreshFailed)
						c.Abort()
						return
					}
					controller.ResponseErrorWithMsg(c,
						controller.CodeFrontEndNeedUseNewToken,
						gin.H{"accessToken": newToken})
					c.Next()
					return
				}
			} else {
				// 解析失败，Token不合法
				zap.L().Error("received an illegal token", zap.String("accessToken", accessToken))
				controller.ResponseError(c, controller.CodeTokenIsInvalid)
				c.Abort()
				return
			}
		}
		// AccessToken正确且未过期
		// 将JWT中携带的用户ID和AccessToken存到中间件链路中
		c.Set(CtxUserIdKey, mc.UserID)
		c.Set(CtxAToken, accessToken)
		c.Next()
	}
}

// JWTLimitLoginMiddleware 限制同一账号同一时间只能一台设备登录
func JWTLimitLoginMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 从Redis中获取AccessToken
		token, err := redis.GetAccessTokenByUID(c.GetString(CtxUserIdKey))
		if err != nil {
			// 获取失败
			controller.ResponseError(c, controller.CodeServeBusy)
			c.Abort()
			return
		}
		// 获取上一步JWT存储的AccessToken
		aToken, ok := c.Get(CtxAToken)
		if !ok {
			// 获取失败
			controller.ResponseError(c, controller.CodeServeBusy)
			c.Abort()
			return
		}
		if aToken != token {
			// 如果两次AccessToken不同，说明超过最大登录终端数量
			controller.ResponseError(c, controller.CodeExceedMaxTerminalNum)
			c.Abort()
			return
		}
		c.Next()
	}
}
