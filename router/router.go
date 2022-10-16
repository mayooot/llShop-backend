package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"shop-backend/controller"
	_ "shop-backend/docs"
	"shop-backend/logger"
	"shop-backend/middleware"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// 将模式设置为发布模式，控制台不会打印日志
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// 添加跨域、日志、恢复中间件
	r.Use(middleware.Cors(), logger.GinLogger(), logger.GinRecovery(true))
	// swagger接口文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// v1路由组不使用校验JWT中间件
	v1 := r.Group("/api/v1")
	{
		// 获取验证码
		v1.GET("/phone", controller.SendVerifyCodeHandler)
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
	// 刷新AccessToken
	v2.PUT("refreshToken/:id", middleware.JWTAuthRefreshMiddleware())
	// 添加校验JWT、限制终端设备中间件
	v2.Use(middleware.JWTAuthMiddleware(), middleware.JWTLimitLoginMiddleware())
	{
		v3 := v2.Use(middleware.JWTCheckUID())
		{
			v3.DELETE("/exit:id", controller.SignOutHandler)
			// 获取用户简略信息，用于商城header显示
			v3.GET("/someinfo/:id", controller.SomeInfoHandler)
			// 获取用户个人信息，用于个人资料显示
			v3.GET("/infos/:id", controller.UserInfosHandler)
		}

		// 用户修改个人资料
		v2.PUT("/infos/update", controller.UserInfosUpdateHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		controller.ResponseErrorWithMsg(c, http.StatusBadRequest, gin.H{"msg": "404"})
	})
	return r
}
