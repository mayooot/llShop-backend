package mysql

import (
	"errors"
	"go.uber.org/zap"
	"shop-backend/models/pojo"
)

// InsertCartProduct 添加一条商品sku信息到用户购物车下
func InsertCartProduct(userID, skuID int64, count int, specification string) error {
	err := db.Create(&pojo.Cart{
		UserID:        userID,
		SkuID:         skuID,
		Specification: specification,
		Count:         count,
		Selected:      0, // 代表未被勾选
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// SelectOneCartProductByUIDAndSkuId 根据用户ID和商品skuID和规格查询用户购物车中是否已经有该商品的记录
func SelectOneCartProductByUIDAndSkuId(userID, skuID int64, specification string) (*pojo.Cart, bool) {
	cart := new(pojo.Cart)
	result := db.Where("user_id = ? and sku_id = ? and specification = ?", userID, skuID, specification).First(cart)
	if result.Error != nil || result.RowsAffected == 0 {
		zap.L().Error("根据用户ID和skuID删除购物车商品记录失败", zap.Error(result.Error), zap.Int64("uid", userID), zap.Int64("skuID", skuID))
		return nil, false
	}
	return cart, true
}

// UpdateCartProductByUIDAndSkuId 根据用户ID和商品skuID更新用户购物车下该商品数量
func UpdateCartProductByUIDAndSkuId(userID, skuID int64, count int) error {
	result := db.Model(&pojo.Cart{}).Where("user_id = ? and sku_id = ?", userID, skuID).Update("count", count)
	if result.Error != nil || result.RowsAffected <= 0 {
		zap.L().Error("更新用户购物车商品数量失败", zap.Error(result.Error), zap.Int64("uid", userID), zap.Int64("skuID", skuID))
		return errors.New("更新用户购物车商品数量失败")
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
func UpdateCartProductSelected(userID, skuID int64, selected int, specification string) error {
	result := db.Where("user_id = ? and sku_id = ? and specification = ?", userID, skuID, specification).Updates(&pojo.Cart{
		Selected: int8(selected),
	})
	if result.Error != nil {
		zap.L().Error("根据用户ID和skuID修改购物车商品勾选状态失败", zap.Error(result.Error), zap.Int64("uid", userID), zap.Int64("skuID", skuID), zap.Int("selected", selected))
		return errors.New("根据用户ID和skuID修改购物车商品勾选状态失败")
	}
	return nil
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
