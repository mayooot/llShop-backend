package middleware

import (
	"github.com/gin-gonic/gin"
	"shop-backend/controller"
	"shop-backend/utils/check"
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
			// 解析失败，Token不合法
			controller.ResponseError(c, controller.CodeTokenIsInvalid)
			c.Abort()
			return
		}
		// 将当前请求的UserID信息保存到请求的上下文中
		c.Set(CtxUserIdKey, mc.UserId)
		c.Set(CtxAToken, parts[1])
		// 后续的请求可以通过c.Get(CtxUserIdKey)来获取当前请求的用户信息
		c.Next()
	}
}
