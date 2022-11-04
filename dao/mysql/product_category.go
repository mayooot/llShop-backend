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

// SelectFirstCategory 查询一级分类信息(level = 0)
func SelectFirstCategory() ([]*pojo.ProductCategory, error) {
	categories := make([]*pojo.ProductCategory, 0)
	result := db.Where("level = 0").Find(&categories)
	if result.Error != nil {
		zap.L().Error("查询一级分类信息失败", zap.Error(result.Error))
		return nil, result.Error
	}
	return categories, nil
}

// SelectSecondCategory 查询二级分类信息(level = 0)
func SelectSecondCategory() ([]*pojo.ProductCategory, error) {
	categories := make([]*pojo.ProductCategory, 0)
	result := db.Where("level = 1").Find(&categories)
	if result.Error != nil {
		zap.L().Error("查询二级分类信息失败", zap.Error(result.Error))
		return nil, result.Error
	}
	return categories, nil
}
