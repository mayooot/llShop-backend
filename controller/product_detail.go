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
// @Param spuID path string true "spuID"
// @Router /pms/product/detail/{spuID} [get]
func ProductDetailHandler(c *gin.Context) {
	spuIDStr := c.Param("spuID")
	spuID, err := strconv.ParseInt(spuIDStr, 10, 64)
	if err != nil {
		zap.L().Error("商品详情接口，spuID不能转换为int64类型", zap.String("spuIDStr", spuIDStr))
		ResponseError(c, CodeServeBusy)
		return
	}
	data, err := logic.GetProductDetail(spuID)
	if err != nil {
		zap.L().Error("商品详情接口，获取商品详情信息失败", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}
