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
		umsGroup.POST("/signup", controller.UserSignUpHandler)
		// 登录
		umsGroup.POST("/login", controller.UserLoginHandler)
	}
	// 鉴权路由组，包含JWT校验中间件，限制用户多端登录中间件
	jwtGroup := umsGroup
	jwtGroup.Use(middleware.JWTAuthMiddleware())
	{
		// 获取用户简略信息，用于商城header显示
		jwtGroup.GET("/someinfo", controller.UserSomeInfoHandler)
		// 获取用户个人信息，用于个人资料显示
		jwtGroup.GET("/infos", controller.UserInfosHandler)
		// 用户修改个人资料
		jwtGroup.PUT("/infos/update", controller.UserInfosUpdateHandler)
		// 用户修改头像
		jwtGroup.POST("/infos/update/avatar", controller.UserInfoUpdateAvatarHandler)
		// 用户退出
		jwtGroup.DELETE("/exit", controller.UserSignOutHandler)
	}

	// 商品路由组
	pmsGroup := commonGroup.Group("/pms/product")
	{
		// 商品分类
		pmsGroup.GET("/category/list", controller.ProductCategoryListHandler)
		// 商品分类的属性列表
		pmsGroup.GET("/attribute/bycategory/:categoryID", controller.ProductAttributeByCategoryIDHandler)
		// 商品搜索接口
		pmsGroup.POST("/search", controller.ProductSearchHandler)
		// 商品详情接口
		pmsGroup.GET("/detail/:skuID", controller.ProductDetailHandler)
	}

	// 购物车路由组，需要鉴权
	cartGroup := commonGroup.Group("/oms/cart").Use(middleware.JWTAuthMiddleware())
	{
		// 获取用户购物车列表
		cartGroup.GET("/list", controller.OrderCartListHandler)
		// 添加商品到购物车
		cartGroup.POST("/add", controller.OrderCartAddHandler)
		// 从购物车中移除商品
		cartGroup.DELETE("/remove", controller.OrderCartRemoveHandler)
		// 获取用户购物车中商品的数量
		cartGroup.GET("/list/count", controller.OrderCartListCountHandler)
		// 修改购物车商品勾选状态
		cartGroup.PUT("/product/status", controller.OrderCartUpdateSelectedHandler)
	}

	// 订单路由组，需要鉴权
	orderGroup := commonGroup.Group("/oms/order").Use(middleware.JWTAuthMiddleware())
	{
		// 生成预提交订单
		orderGroup.POST("/presubmit", controller.OrderPreSubmitHandler)
		// 提交订单
		orderGroup.POST("/submit", controller.OrderSubmitHandler)
		// 获取用户所有的订单
		orderGroup.GET("/all", controller.OrderGetAllHandler)
		// 获取订单明细
		orderGroup.GET("/one/:num", controller.OrderGetOneOrderItemHandler)
		// 删除订单
		orderGroup.DELETE("/del/:num", controller.OrderDelOrderHandler)
		// 支付接口
		orderGroup.POST("/pay", controller.AlipayHandler)
	}
	// 支付订单回调接口
	commonGroup.GET("/oms/order/pay/notify", controller.AlipayNotifyHandler)

	// 收货地址路由组，需要鉴权
	receiverAddressGroup := commonGroup.Group("/user/receiveraddress").Use(middleware.JWTAuthMiddleware())
	{
		// 获取数据库中所有的地址，用于用户选择
		receiverAddressGroup.GET("/list", controller.UserReceiverAddressListHandler)
		// 增加一条收货地址
		receiverAddressGroup.POST("/add", controller.UserReceiverAddressAddHandler)
		// 删除一条收货地址
		receiverAddressGroup.DELETE("/delete/:id", controller.UserReceiverAddressDeleteHandler)
		// 获取用户的收货地址列表
		receiverAddressGroup.GET("/my", controller.UserReceiverAddressPersonHandler)
		// 修改一条收货地址
		receiverAddressGroup.PUT("/update", controller.UserReceiverAddressUpdateHandler)
	}

	// 秒杀商品路由组
	secKillGroup := commonGroup.Group("/seckill").Use(middleware.RateLimitMiddleware(), middleware.JWTAuthMiddleware())
	{
		secKillGroup.GET("/sku/list", controller.SecKillAllSkuHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		controller.ResponseErrorWithMsg(c, http.StatusBadRequest, gin.H{"msg": "404"})
	})
	return r
}
