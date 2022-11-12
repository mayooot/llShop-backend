package middleware

import (
	"github.com/gin-gonic/gin"
	"shop-backend/controller"
	"shop-backend/dao/redis"
)

// RateLimitRedisMiddleware 限制用户5秒钟只能抢购一次商品
func RateLimitRedisMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		ok := redis.SetNXSecKillUID(c.GetInt64("uid"))
		if !ok {
			// 用户已经在5秒钟内抢购过商品
			controller.ResponseBadError(c, controller.CodeToManyRequest)
			c.Abort()
			return
		}
		// 用户还未抢购过，放行
		c.Next()
	}
}
