package mysql

import (
	"errors"
	"go.uber.org/zap"
	"shop-backend/models/dto"
	"shop-backend/models/pojo"
	"strconv"
)

// CheckOrderProduct 检查预提交订单中的商品是否还在上架，购买数量是否超过库存
func CheckOrderProduct(cartProduct *dto.CartProduct, uid int64) (*pojo.Cart, *pojo.Sku, error) {
	tx := db.Begin()

	// 查询出sku，并校验上架状态和库存
	// 类型转换
	skuID, _ := strconv.ParseInt(cartProduct.SkuID, 10, 64)
	sku := new(pojo.Sku)
	if err := tx.Where("id = ?", skuID).First(&sku).Error; err != nil {
		tx.Rollback()
		zap.L().Error("使用skuID查询sku信息失败", zap.Int64("skuID", skuID))
		return nil, nil, err
	}
	if sku.Valid == 0 || sku.Stock < cartProduct.Count {
		tx.Rollback()
		zap.L().Error("商品已下架或者购买数量大于库存")
		return nil, nil, errors.New("商品已下架或者购买数量大于库存")
	}

	// 返回购物车对象，用于构建购物车展示对象
	cartPojo := new(pojo.Cart)
	if err := tx.Where("user_id = ? and sku_id = ? and specification = ?", uid, skuID, cartProduct.Specification).First(cartPojo).Error; err != nil {
		tx.Rollback()
		zap.L().Error("根据用户ID、skuID、商品规格获取一条购物车数据失败")
		return nil, nil, errors.New("根据用户ID、skuID、商品规格获取一条购物车数据失败")
	}
	tx.Commit()
	return cartPojo, sku, nil
}
