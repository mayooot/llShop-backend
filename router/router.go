package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-backend/controller"
	"shop-backend/logger"
	"shop-backend/middleware"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// 将模式设置为发布模式，控制台不会打印日志
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// v1路由组不使用校验JWT中间件
	v1 := r.Group("/api/v1")
	{
		// 获取验证码
		v1.POST("/phone", controller.SendVerifyCodeHandler)
		// 注册
		v1.POST("/signup", controller.SignUpHandler)
		// 登录
		v1.POST("/login", controller.LoginHandler)
		// 测试
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "pong"})
		})
	}

	// v2路由组使用校验JWT中间件
	v2 := r.Group("/api/v1")
	v2.Use(middleware.JWTAuthMiddleware())
	{
		// 获取用户简略信息，用于商城header显示
		v2.GET("/someinfo/:id", controller.SomeInfoHandler)
		// 获取用户个人信息
		v2.GET("/infos")
	}

	r.NoRoute(func(c *gin.Context) {
		controller.ResponseErrorWithMsg(c, http.StatusBadRequest, gin.H{"msg": "404"})
	})
	return r
}
