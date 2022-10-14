package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-backend/controller"
	"shop-backend/logger"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// 将模式设置为发布模式，控制台不会打印日志
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/v1")
	{
		// 获取验证码
		v1.POST("/phone", controller.SendVerifyCodeHandler)
		// 注册
		v1.POST("/signup", controller.SignUpHandler)
		v1.POST("/login", controller.LoginHandler)

		// 测试
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "pong"})
		})
	}

	r.NoRoute(func(c *gin.Context) {
		controller.ResponseErrorWithMsg(c, http.StatusBadRequest, gin.H{"msg": "404"})
	})
	return r
}
