package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"reflect"
	"shop-backend/logic"
	"shop-backend/models/dto"
)

// ProductSearchHandler 支持多条件的商品搜索接口
// @Summary 商品搜索接口
// @Description 支持多条件的商品搜索接口。如未指定分页字段，则默认返回第1页的前20条数据。<br> case1: 主页搜索框搜索功能 <br> {"keyword": "犬","pageNo": "1","pageSize": "10"} <br> <br> case2: 使用二级分类ID搜索商品 <br> {"productCategoryId": "14","pageNo": "1","pageSize": "10"} <br> <br> case2: 使用商品属性列表搜索商品 <br> {"productAttributeIds": [14,15,16],"pageNo": "1","pageSize": "10"} <br> <br>
// @Tags 商品相关接口
// @Produce  json
// @Param searchCondition body dto.SearchCondition true "搜索条件"
// @Router /pms/product/search [post]
func ProductSearchHandler(c *gin.Context) {
	condition := dto.NewCondition()
	nilCondition := new(dto.SearchCondition)
	err := c.ShouldBindJSON(condition)
	if err != nil {
		zap.L().Error("商品搜索接口，前端传递的条件有误", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	if reflect.DeepEqual(condition, nilCondition) {
		// 如果condition全部条件为空，不应返回所有数据(防止用户获取sku表所有数据)
		zap.L().Error("商品搜索接口，前端传递的条件全部为空", zap.Error(err))
		ResponseError(c, CodeSearchConditionIsNil)
		return
	}
	// 调用logic层根据条件查询商品
	data, err := logic.Search(condition)
	if err != nil {
		zap.L().Error("商品搜索logic层错误", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}
