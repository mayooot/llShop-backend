package middleware

import (
	"github.com/gin-gonic/gin"
	"shop-backend/controller"
	"shop-backend/dao/redis"
	"shop-backend/utils/check"
	"shop-backend/utils/gen"
	"strconv"
	"strings"
)

var (
	CtxUserIdKey = "uid"
	CtxAToken    = "atoken"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式：1.放在请求头中 2.放在请求体中 3.放在URI中
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization: Bearer xxx.xxx.xxx
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// 未携带Token
			controller.ResponseError(c, controller.CodeTokenIsEmpty)
			c.Abort()
			return
		}
		// 按照空格进行分割
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 发送Token格式有误
			controller.ResponseError(c, controller.CodeTokenIsWrongFormat)
			c.Abort()
			return
		}
		// 解析前端传递的Token
		mc, err := check.CheckAToken(parts[1])
		if err != nil {
			if err == check.ErrorATokenExpired {
				// Token过期错误
				controller.ResponseError(c, controller.CodeTokenExpire)
				c.Abort()
				return
			}
			// 解析失败，Token不合法
			controller.ResponseError(c, controller.CodeTokenIsInvalid)
			c.Abort()
			return
		}
		// 将当前请求的UserID信息保存到请求的上下文中
		// 将uid转换成string
		uidStr := strconv.FormatInt(mc.UserId, 10)
		c.Set(CtxUserIdKey, uidStr)
		c.Set(CtxAToken, parts[1])
		// 后续的请求可以通过c.Get(CtxUserIdKey)来获取当前请求的用户信息
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

// JWTAuthRefreshMiddleware 刷新AccessToken
func JWTAuthRefreshMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeTokenIsEmpty)
			c.Abort()
			return
		}
		// 按照空格进行分割
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			controller.ResponseError(c, controller.CodeTokenIsWrongFormat)
			c.Abort()
			return
		}
		// 前端传递过来的两个Token使用&拼接
		tokens := strings.SplitN(parts[1], "&", 2)
		if len(tokens) != 2 {
			controller.ResponseError(c, controller.CodeTokenIsWrongFormat)
			c.Abort()
			return
		}
		aToken, err := gen.RefreshToken(tokens[0], tokens[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeTokenRefreshFailed)
			c.Abort()
			return
		}
		if aToken == "" {
			controller.ResponseError(c, controller.CodeAccessTokenIsLiving)
			c.Abort()
			return
		}
		controller.ResponseSuccess(c, gin.H{
			"AccessToken": aToken,
		})
		c.Next()
	}
}

// JWTCheckUID 检查前端传递的用户ID是否和JWTAuthMiddleware存储的ID相同
func JWTCheckUID() func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr != c.GetString("uid") {
			controller.ResponseError(c, controller.CodeServeBusy)
			c.Abort()
			return
		}

		c.Next()
	}
}
