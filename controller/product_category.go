package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
)

// ProductCategoryListHandler 获取商品分类
// @Summary 获取所有商品分类信息
// @Description 返回商品一级分类、二级分类列表
// @Tags 商品相关接口
// @Produce  json
// @Router /pms/product/category/list [get]
func ProductCategoryListHandler(c *gin.Context) {
	categories, err := logic.GetAllCategory()
	if err != nil {
		zap.L().Error("获取商品分类信息失败", zap.Error(err))
		ResponseError(c, CodeRequestAllCategoryFailed)
	}
	ResponseSuccess(c, categories)
}
