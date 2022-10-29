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

// OrderAddCartHandler 添加商品到购物车
// @Summary 添加商品到购物车
// @Description 前端传递JSON类型对象，后端完成将商品添加到用户购物车中。
// @Tags 购物车相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param CartProduct body dto.CartProduct true "购物车商品结构体"
// @Router /oms/cart/add [post]
func OrderAddCartHandler(c *gin.Context) {
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

// OrderRemoveCartHandler 从购物车中删除指定商品
// @Summary 从购物车中删除指定商品
// @Description 前端传递商品skuID，后端完成删除。
// @Tags 购物车相关接口
// @Produce json
// @param Authorization header string true "Bearer AToken&RToken"
// @Param skuID path string true "商品skuID"
// @Router /oms/cart/remove/{skuID} [delete]
func OrderRemoveCartHandler(c *gin.Context) {
	skuIDStr := c.Param("skuID")
	skuID, err := strconv.ParseInt(skuIDStr, 10, 64)
	if err != nil {
		zap.L().Error("删除商品到购物车接口，skuIDStr不能转换为int64类型", zap.String("skuIDStr", skuIDStr))
		ResponseError(c, CodeInvalidParams)
		return
	}

	if err = logic.DelCartProduct(c.GetInt64("uid"), skuID); err != nil {
		zap.L().Error("删除商品到购物车接口，删除商品失败", zap.Int64("skuID", skuID), zap.Error(err))
		ResponseError(c, CodeDeleteCartProductFailed)
		return
	}
	ResponseSuccess(c, "删除成功")
}
