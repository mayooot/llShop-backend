package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"strconv"
)

// ProductDetailHandler 商品详情接口
// @Summary 商品详情接口
// @Description 前端以path的形式传递spuID，后端返回该商品详情信息
// @Tags 商品相关接口
// @Produce  json
// @Param skuID path string true "skuID:1000002"
// @Router /pms/product/detail/{skuID} [get]
func ProductDetailHandler(c *gin.Context) {
	skuIDStr := c.Param("skuID")
	skuID, err := strconv.ParseInt(skuIDStr, 10, 64)
	if err != nil {
		zap.L().Error("商品详情接口，skuIDStr不能转换为int64类型", zap.String("skuIDStr", skuIDStr))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 多协程获取商品详情信息
	data, err := logic.GetProductDetailWithConcurrent(skuID)
	if err != nil {
		zap.L().Error("商品详情接口，获取商品详情信息失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}
