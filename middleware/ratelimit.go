package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"shop-backend/controller"
	"shop-backend/utils/gen"
	"time"
)

// RateLimitMiddleware 令牌桶限流中间件。令牌桶按固定的速率往桶中放入令牌，并且只能从桶中取出令牌后才能通过
func RateLimitMiddleware() func(c *gin.Context) {
	// 创建填充速度为指定速率和容量大小的令牌桶
	// NewBucketWithRate(0.1, 200) 表示每秒填充20个令牌
	bucket := ratelimit.NewBucketWithRate(0.1, 200)
	return func(c *gin.Context) {
		// 如果取不到令牌，最大等待5秒。如果5秒后仍然没有取到令牌，则中断本次请求
		_, ok := bucket.TakeMaxDuration(1, time.Second*5)
		if !ok {
			controller.ResponseBadError(c, controller.CodeToManyRequest)
			c.Abort()
			return
		}
		// 成功取到令牌就放行
		c.Next()
		c.Set("uid", gen.GenSnowflakeID())
	}
}
