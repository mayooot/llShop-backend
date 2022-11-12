package mysql

import (
	"errors"
	"gorm.io/gorm"
	"shop-backend/models/pojo"
)

// UpdateSecKillProductStock 使用乐观锁修改秒杀商品库存
func UpdateSecKillProductStock(skuID int64) error {
	tx := db.Begin()
	// 1. 根据主键ID查询出商品
	var product pojo.SecKillSku
	result := tx.First(&product, skuID)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("查询秒杀商品失败")
	}
	// 2. 更新
	result = tx.Model(&product).Updates(map[string]interface{}{
		"sale":  gorm.Expr("sale + ?", 1),
		"stock": gorm.Expr("stock - ?", 1),
	})
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		// 秒杀失败
		return errors.New("秒杀失败")
	}
	tx.Commit()
	return nil
}

// SelectAllSecKillSku 获取所有正在秒杀的商品
func SelectAllSecKillSku() ([]*pojo.SecKillSku, error) {
	data := make([]*pojo.SecKillSku, 0)
	if err := db.Model(&pojo.SecKillSku{}).Find(&data).Error; err != nil {
		return nil, errors.New("获取所有正在秒杀的商品失败")
	}
	return data, nil
}
