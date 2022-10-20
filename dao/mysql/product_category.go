package mysql

import (
	"go.uber.org/zap"
	"shop-backend/models/pojo"
)

// SelectAllCategory 查询所有分类信息
func SelectAllCategory() ([]*pojo.ProductCategory, error) {
	categories := make([]*pojo.ProductCategory, 0)
	result := db.Find(&categories)
	if result.Error != nil {
		zap.L().Error("SelectAllCategory 查询所有分类信息失败", zap.Error(result.Error))
		return nil, result.Error
	}
	return categories, nil
}
