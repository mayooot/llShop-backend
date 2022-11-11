package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"shop-backend/models/dto"
	"strconv"
)

// AlipayNotifyHandler 支付回调接口
// @Summary 支付回调接口
// @Description 支付完成后，支付宝会调到前端页面。此时前端页面需要请求本接口
// @Tags 订单相关接口
// @Produce json
// @Router /oms/order/pay/notify [get]
func AlipayNotifyHandler(c *gin.Context) {
	ResponseSuccessWithMsg(c, "支付成功，请在我的订单查看详情🙉", nil)
}

// AlipayHandler 支付接口
// @Summary 支付接口
// @Description 用户携带Token和订单号请求本接口，将会返回支付宝支付的url
// @Tags 订单相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param Alipay body dto.AliPay true "支付宝支付结构体"
// @Router /oms/order/pay [post]
func AlipayHandler(c *gin.Context) {
	aliPay := new(dto.AliPay)
	err := c.ShouldBindJSON(aliPay)
	if err != nil {
		zap.L().Error("支付宝支付接口，传递参数错误")
		ResponseError(c, CodeServeBusy)
		return
	}

	orderNum, err := strconv.ParseInt(aliPay.OrderNum, 10, 64)
	if err != nil {
		zap.L().Error("支付宝支付接口，订单号转为int64错误")
		ResponseError(c, CodeServeBusy)
		return
	}

	payUrl, err := logic.CreateAlipayOrder(c.GetInt64("uid"), orderNum)
	if err != nil || payUrl == "" {
		zap.L().Error("支付宝支付接口，拉起支付宝支付失败")
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, payUrl)
}
