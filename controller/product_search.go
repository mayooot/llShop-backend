package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"shop-backend/models/dto"
)

// ProductSearchHandler 商品搜索接口
func ProductSearchHandler(c *gin.Context) {
	condition := new(dto.SearchCondition)
	err := c.ShouldBindJSON(condition)
	if err != nil {
		zap.L().Error("商品搜索接口，前端传递的条件有误", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
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
