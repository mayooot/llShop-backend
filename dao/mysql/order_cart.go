package mysql

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"shop-backend/models/pojo"
)

// InsertCartProduct 添加一条商品sku信息到用户购物车下
func InsertCartProduct(userID, skuID int64, count int, specification string) error {
	err := db.Create(&pojo.Cart{
		UserID:        userID,
		SkuID:         skuID,
		Specification: specification,
		Count:         count,
		Selected:      2, // 代表未被勾选
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// SelectOneCartProductByUIDAndSkuId 根据用户ID和商品skuID和规格查询用户购物车中是否已经有该商品的记录
func SelectOneCartProductByUIDAndSkuId(userID, skuID int64, specification string) (*pojo.Cart, bool) {
	cart := new(pojo.Cart)
	result := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? and sku_id = ? and specification = ?", userID, skuID, specification).First(cart)
	if result.Error != nil || result.RowsAffected <= 0 {
		// 商品不存在
		return nil, false
	}
	return cart, true
}

// UpdateCartProductByUIDAndSkuId 根据用户ID和商品skuID更新用户购物车下该商品数量
func UpdateCartProductByUIDAndSkuId(userID, skuID int64, count int) error {
	for time := 1; time <= 10; time++ {
		// 开启事务
		tx := db.Begin()
		// 先查询
		var cart pojo.Cart
		tx.Where("user_id = ? and sku_id = ?", userID, skuID).First(&cart)
		// 更新
		result := tx.Debug().Model(&cart).Update("count", gorm.Expr("count + ?", count))
		if result.Error != nil || result.RowsAffected <= 0 {
			tx.Rollback()
			zap.L().Error("--------更新用户购物车商品数量失败----------", zap.Error(result.Error), zap.Int64("RowsAffected", result.RowsAffected))
			zap.L().Error("--------尝试自旋等待更新", zap.Int("times", time))
		} else {
			tx.Commit()
			break
		}
	}

	return nil
}

// DelCartProductBySkuIDAndUID 根据用户ID和skuID删除购物车商品记录
func DelCartProductBySkuIDAndUID(userID, skuID int64, specification string) error {
	result := db.Where("user_id = ? and sku_id = ? and specification = ?", userID, skuID, specification).Delete(&pojo.Cart{})
	if result.Error != nil || result.RowsAffected == 0 {
		zap.L().Error("根据用户ID和skuID删除购物车商品记录失败", zap.Error(result.Error), zap.Int64("uid", userID), zap.Int64("skuID", skuID))
		return errors.New("根据用户ID和skuID删除购物车商品记录失败")
	}
	return nil
}

// UpdateCartProductSelected 根据用户ID和skuID修改购物车商品勾选状态
func UpdateCartProductSelected(userID, skuID int64, selected int, specification string) (*pojo.Cart, error) {
	cart := new(pojo.Cart)
	cart.Selected = int8(selected)

	result := db.Where("user_id = ? and sku_id = ? and specification = ?", userID, skuID, specification).Updates(cart)
	if result.Error != nil {
		zap.L().Error("根据用户ID和skuID修改购物车商品勾选状态失败", zap.Error(result.Error), zap.Int64("uid", userID), zap.Int64("skuID", skuID), zap.Int("selected", selected))
		return nil, errors.New("根据用户ID和skuID修改购物车商品勾选状态失败")
	}
	return cart, nil
}

// SelectCartList 获取用户购物车信息集合
func SelectCartList(userID int64) ([]*pojo.Cart, error) {
	cartList := make([]*pojo.Cart, 0)
	result := db.Where("user_id = ?", userID).Find(&cartList)
	if result.Error != nil {
		zap.L().Error("获取用户购物车信息失败", zap.Int64("uid", userID))
		return nil, result.Error
	}
	return cartList, nil
}
