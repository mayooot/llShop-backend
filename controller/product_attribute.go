package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"strconv"
)

// AttributeByCategoryIDHandler 通过二级分类ID获取商品属性
// @Summary 使用商品二级分类ID获取商品属性
// @Description 前端以param的形式传递商品二级分类ID，后端返回该商品的所有属性
// @Tags 商品相关接口
// @Produce  json
// @Param category path string true "商品二级分类ID"
// @Router /pms/product/attribute/byCategoryID [get]
func AttributeByCategoryIDHandler(c *gin.Context) {
	categoryIDStr := c.Param("categoryID")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	attributes, err := logic.GetAllAttribute(categoryID)
	if err != nil {
		zap.L().Error("通过二级分类ID获取商品属性失败", zap.Error(err))
		ResponseError(c, CodeRequestAllAttributeFailed)
		return
	}
	ResponseSuccess(c, attributes)
}
