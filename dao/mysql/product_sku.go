package mysql

import (
	"go.uber.org/zap"
	"shop-backend/models/pojo"
)

// SelectSkuBySkuID 使用skuID查询sku信息
func SelectSkuBySkuID(skuID int64) (*pojo.Sku, error) {
	sku := new(pojo.Sku)
	if err := db.Where("id = ?", skuID).First(&sku).Error; err != nil {
		zap.L().Error("使用skuID查询sku信息失败", zap.Int64("skuID", skuID))
		return nil, err
	}
	return sku, nil
}
