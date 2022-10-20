package controller

import (
	"github.com/gin-gonic/gin"
	"shop-backend/logic"
)

// CategoryListHandler 获取商品分类
// @Summary 获取所有商品分类信息
// @Description 返回商品一级分类、二级分类列表
// @Tags 商品相关接口
// @Produce  json
// @Router /pms/product/category/list [get]
func CategoryListHandler(c *gin.Context) {
	categories, err := logic.GetAllCategory()
	if err != nil {
		ResponseError(c, CodeRequestAllCategoryFailed)
	}
	ResponseSuccess(c, categories)
}
