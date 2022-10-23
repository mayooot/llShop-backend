package logic

import (
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/models/dto"
	"shop-backend/models/vo"
	"strconv"
)

// Search 多条件搜索业务
func Search(condition *dto.SearchCondition) (*vo.Page[[]*vo.Product], error) {
	// 获取符合条件的sku集合
	products, _, err := mysql.BaseSearchCondition(condition, true)
	if err != nil {
		zap.L().Error("mysql层BaseSearchCondition(分页) 查询失败", zap.Error(err))
		return nil, err
	}
	// 获取总页数
	_, totalPage, err := mysql.BaseSearchCondition(condition, false)
	if err != nil {
		zap.L().Error("mysql层BaseSearchCondition(不分页) 查询失败", zap.Error(err))
		return nil, err
	}
	page := &vo.Page[[]*vo.Product]{
		PageNo:    condition.PageNo,
		PageSize:  condition.PageSize,
		TotalPage: strconv.Itoa(totalPage),
		Data:      products,
	}
	return page, nil
}
