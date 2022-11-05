package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"shop-backend/models/dto"
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
	ResponseSuccessWithMsg(c, CodeCreatePreSubmitOrderSuccess, orderVO)
}
