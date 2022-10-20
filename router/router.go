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

	// 普通路由组，只包含跨域、日志、恢复中间件
	commonGroup := r.Group("/api")
	umsGroup := commonGroup.Group("/user")
	{
		// 获取验证码
		umsGroup.GET("/phone", controller.SendVerifyCodeHandler)
		// 注册
		umsGroup.POST("/signup", controller.SignUpHandler)
		// 登录
		umsGroup.POST("/login", controller.LoginHandler)
	}
	// 鉴权路由组，包含JWT校验中间件，限制用户多端登录中间件
	jwtGroup := umsGroup
	jwtGroup.Use(middleware.JWTAuthMiddleware())
	{
		// 获取用户简略信息，用于商城header显示
		jwtGroup.GET("/someinfo", controller.SomeInfoHandler)
		// 获取用户个人信息，用于个人资料显示
		jwtGroup.GET("/infos", controller.UserInfosHandler)
		// 用户修改个人资料
		jwtGroup.PUT("/infos/update", controller.UserInfosUpdateHandler)
		// 用户修改头像
		jwtGroup.POST("/infos/update/avatar", controller.UserInfoUpdateAvatarHandler)
		// 用户退出
		jwtGroup.DELETE("/exit", controller.SignOutHandler)
	}

	// 商品路由组
	pmsGroup := commonGroup.Group("/pms/product")
	{
		pmsGroup.GET("/category/list", controller.CategoryListHandler)
		pmsGroup.GET("/attribute/bycategory/:categoryID", controller.AttributeByCategoryIDHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		controller.ResponseErrorWithMsg(c, http.StatusBadRequest, gin.H{"msg": "404"})
	})
	return r
}
