package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"shop-backend/models/dto"
	"strconv"
)

// OrderCartListHandler 用户购物车商品列表
// @Summary 获取购物车商品列表
// @Description 前端需要携带Token用来鉴权，鉴权通过后返回用户购物车商品列表
// @Tags 购物车相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /oms/cart/list [get]
func OrderCartListHandler(c *gin.Context) {
	data, err := logic.GetCarProductList(c.GetInt64("uid"))
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}

// OrderCartListCountHandler 用户购物车商品数量，用于展示购物车logo上的数字
// @Summary 用户购物车商品数量
// @Description 前端需要携带Token用来鉴权，鉴权通过后返回用户购物车商品数量，用于展示购物车logo上的数字
// @Tags 购物车相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /oms/cart/list/count [get]
func OrderCartListCountHandler(c *gin.Context) {
	count, err := logic.GetCarProductListCount(c.GetInt64("uid"))
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, count)
}

// OrderCartAddHandler 添加商品到购物车
// @Summary 添加商品到购物车
// @Description 前端传递JSON类型对象，后端完成将商品添加到用户购物车中。
// @Tags 购物车相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param CartProduct body dto.CartProduct true "购物车商品结构体"
// @Router /oms/cart/add [post]
func OrderCartAddHandler(c *gin.Context) {
	cartProduct := new(dto.CartProduct)
	if err := c.ShouldBindJSON(cartProduct); err != nil {
		zap.L().Error("添加商品到购物车接口，传递参数有误", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	skuID, err := strconv.ParseInt(cartProduct.SkuID, 10, 64)
	if err != nil {
		zap.L().Error("添加商品到购物车接口，转换skuID为int64失败", zap.Error(err), zap.String("skuID", cartProduct.SkuID))
		ResponseError(c, CodeInvalidParams)
		return
	}

	count, err := strconv.Atoi(cartProduct.Count)
	if err != nil || count <= 0 {
		zap.L().Error("添加商品到购物车接口，转换count为in失败", zap.Error(err), zap.String("count", cartProduct.Count))
		ResponseError(c, CodeInvalidParams)
		return
	}

	if err = logic.AddCartProduct(c.GetInt64("uid"), skuID, count); err != nil {
		zap.L().Error("添加商品到购物车接口，添加商品到购物车失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, "添加成功")
}

// OrderCartRemoveHandler 从购物车中删除指定商品
// @Summary 从购物车中删除指定商品
// @Description 前端传递商品skuID，后端完成删除。
// @Tags 购物车相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param skuID path string true "商品skuID"
// @Router /oms/cart/remove/{skuID} [delete]
func OrderCartRemoveHandler(c *gin.Context) {
	skuIDStr := c.Param("skuID")
	skuID, err := strconv.ParseInt(skuIDStr, 10, 64)
	if err != nil {
		zap.L().Error("删除购物车商品接口，skuIDStr不能转换为int64类型", zap.String("skuIDStr", skuIDStr))
		ResponseError(c, CodeInvalidParams)
		return
	}

	if err = logic.DelCartProduct(c.GetInt64("uid"), skuID); err != nil {
		zap.L().Error("删除购物车商品接口，删除商品失败", zap.Int64("skuID", skuID), zap.Error(err))
		ResponseError(c, CodeDeleteCartProductFailed)
		return
	}
	ResponseSuccess(c, "删除成功")
}

// OrderCartUpdateSelectedHandler 修改购物车中商品勾选状态
// @Summary 修改购物车中商品勾选状态
// @Description 前端以json格式传递品skuID和Selected，后端完成勾选状态的改变。
// @Tags 购物车相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param CartProduct body dto.CartProductSelected true "购物车商品状态"
// @Router /oms/cart/product/status [put]
func OrderCartUpdateSelectedHandler(c *gin.Context) {
	status := new(dto.CartProductSelected)
	if err := c.ShouldBindJSON(status); err != nil {
		zap.L().Error("修改购物车中商品勾选状态接口，传递参数有误", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	skuID, err := strconv.ParseInt(status.SkuID, 10, 64)
	if err != nil {
		zap.L().Error("修改购物车中商品勾选状态接口，转换skuID为int64失败", zap.Error(err), zap.String("skuID", status.SkuID))
		ResponseError(c, CodeInvalidParams)
		return
	}

	selected, err := strconv.Atoi(status.Selected)
	if err != nil || !(selected == 0 || selected == 1) {
		zap.L().Error("修改购物车中商品勾选状态接口，转换selected为int8失败或勾选状态值不正确", zap.Error(err), zap.String("count", status.Selected))
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err = logic.UpdateCartProductSelected(c.GetInt64("uid"), skuID, selected); err != nil {
		zap.L().Error("修改购物车商品状态，修改状态", zap.Int64("skuID", skuID), zap.Int("selected", selected), zap.Error(err))
		ResponseError(c, CodeUpdateCartProductStatusFailed)
		return
	}
	var info string
	if selected == 0 {
		info = "取消勾选成功🥑"
	} else {
		info = "勾选成功🥑"
	}
	ResponseSuccess(c, info)
}
