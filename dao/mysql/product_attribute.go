package mysql

import (
	"go.uber.org/zap"
	"shop-backend/models/pojo"
)

// SelectAllAttribute 返回所有商品属性
func SelectAllAttribute() ([]*pojo.ProductAttribute, error) {
	attrs := make([]*pojo.ProductAttribute, 0)
	result := db.Find(&attrs)
	if result.Error != nil {
		zap.L().Error("查询所有商品属性", zap.Error(result.Error))
		return nil, result.Error
	}
	return attrs, nil
}

// SelectAttrIDeByCategoryID 查询商品二级分类下对应的所有商品属性ID
func SelectAttrIDeByCategoryID(categoryID int64) ([]*pojo.ProductCategoryAttributeRel, error) {
	caRels := make([]*pojo.ProductCategoryAttributeRel, 0)
	result := db.Where("product_category_id = ?", categoryID).Find(&caRels)
	if result.Error != nil {
		zap.L().Error("查询商品二级分类下对应的所有商品属性ID出错了", zap.Error(result.Error))
		return nil, result.Error
	}
	return caRels, nil
}
