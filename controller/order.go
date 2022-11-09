package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/dao/redis"
	"shop-backend/logic"
	"shop-backend/models/dto"
	"shop-backend/utils/check"
	"strconv"
)

// OrderPreSubmitHandler 生成预提交订单
// @Summary 生成预提交订单
// @Description 用户在购物车点击结算时，生成预提交订单。此时需要获取一个全局唯一的订单号，并在真正提交订单时传递给后端。
// @Tags 订单相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param cartProductList body dto.PreSubmitOrder true "预提交订单结构体"
// @Router /oms/order/presubmit [post]
func OrderPreSubmitHandler(c *gin.Context) {
	preSubmitOrder := new(dto.PreSubmitOrder)
	if err := c.ShouldBindJSON(preSubmitOrder); err != nil {
		zap.L().Error("生成预提交订单，前端传递购物车已勾选商品有误", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	orderVO, err := logic.CreatePreSubmitOrder(preSubmitOrder, c.GetInt64("uid"))
	if err != nil {
		zap.L().Error("生成预提交订单失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccessWithMsg(c, CodeCreatePreSubmitOrderSuccess.Msg(), orderVO)
}

// OrderSubmitHandler 提交订单
// @Summary 提交订单
// @Description 用户需要传递订单号、购买商品列表、收货人信息
// @Tags 订单相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param order body dto.Order true "提交订单结构体"
// @Router /oms/order/submit [post]
func OrderSubmitHandler(c *gin.Context) {
	order := new(dto.Order)
	if err := c.ShouldBindJSON(order); err != nil {
		zap.L().Error("提交订单接口，用户传递参数错误")
		ResponseError(c, CodeInvalidParams)
		return
	}

	orderNum, err := strconv.ParseInt(order.OrderNumber, 10, 64)
	// 参数校验
	if order.OrderNumber == "" || err != nil || len(order.CartProductList) == 0 || order.ReceiverName == "" || !check.VerifyMobileFormat(order.ReceiverPhone) || order.ReceiverAddress == "" {
		// 订单号为空 || 订单转int64失败 || 订单商品集合为空 || 收货人为空 || 手机号格式不正确 || 详细地址为空
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 判断用户传递的订单号是否是后端生成的
	exist := redis.GetOrderNumber(orderNum)
	if !exist {
		// 用户传递的订单号不存在或缓存在Redis中的订单号已过期
		ResponseError(c, CodeOrderNumISNotExistOrExpired)
		return
	}

	// 否则，用户传递的订单号存在；提交订单幂等性由数据库主键的唯一性保证
	if err = logic.CreateSubmitOrder(order, c.GetInt64("uid"), orderNum); err != nil {
		zap.L().Error("提交订单失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccessWithMsg(c, CodeCreateSubmitOrderSuccess.Msg(), nil)
}

// OrderGetAllHandler 获取用户所有的订单
// @Summary 获取用户所有的订单
// @Description 前端需要携带Token，鉴权后返回用户所有的订单主要信息
// @Tags 订单相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /oms/order/all [get]
func OrderGetAllHandler(c *gin.Context) {
	data, err := logic.GetAllOrder(c.GetInt64("uid"))
	if err != nil {
		zap.L().Error("获取用户所有的订单失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, data)
}

// OrderGetOneOrderItemHandler 获取一条订单记录的明细
// @Summary 获取一条订单记录的明细
// @Description 前端需要携带Token并传递订单号，后端返回订单明细
// @Tags 订单相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param num path string true "订单号"
// @Router /oms/order/one/{num} [get]
func OrderGetOneOrderItemHandler(c *gin.Context) {
	idStr := c.Param("num")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("获取一条订单记录的明细接口，idStr不能转换为int64类型", zap.String("idStr", idStr))
		ResponseError(c, CodeInvalidParams)
		return
	}

	data, err := logic.GetOneOrderItem(id)
	if err != nil {
		zap.L().Error("获取一条订单记录的明细失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, data)
}
