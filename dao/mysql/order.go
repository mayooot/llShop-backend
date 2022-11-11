package mysql

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"shop-backend/models/dto"
	"shop-backend/models/pojo"
	"strconv"
	"time"
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

// CreateOrderAndOrderItem 生成订单 && 校验库存和商品状态 && 生成订单明细
func CreateOrderAndOrderItem(orderDTO *dto.Order, uid, orderNum int64) error {
	// 生产订单对象
	order := new(pojo.Order)
	// 订单表记录主键自己生成，不需要数据库自增自动生成。目的是使用主键的唯一性来保证提交订单服务幂等性
	order.ID = orderNum
	order.UserID = uid
	order.ReceiverName = orderDTO.ReceiverName
	order.ReceiverPhone = orderDTO.ReceiverPhone
	order.ReceiverAddress = orderDTO.ReceiverAddress
	// 设置订单状态为代付款
	order.OrderStatus = 6
	// 设置支付状态为未支付
	order.PayStatus = 3
	// 设置订单过期时间为30分钟
	expire, _ := time.ParseDuration("30m")
	order.ExpirationTime = time.Now().Add(expire)

	// 初始化总金额和购买商品总数量
	totalMoney := decimal.NewFromFloat(0)
	totalNum := 0

	tx := db.Begin()

	// 校验库存 && 生成订单明细
	for _, product := range orderDTO.CartProductList {
		// 查询sku，并校验上架状态和库存
		// 类型转换
		skuID, _ := strconv.ParseInt(product.SkuID, 10, 64)
		sku := new(pojo.Sku)
		result := tx.Where("id = ?", skuID).First(&sku)
		if result.Error != nil || result.RowsAffected <= 0 {
			tx.Rollback()
			zap.L().Error("使用skuID查询sku信息失败", zap.Int64("skuID", skuID))
			return errors.New("使用skuID查询sku信息失败")
		}

		if sku.Valid == 0 || sku.Stock < product.Count {
			tx.Rollback()
			zap.L().Error("商品已下架或者购买数量大于库存")
			return errors.New("商品已下架或者购买数量大于库存")
		}

		// 扣减库存
		result = tx.Model(&pojo.Sku{}).Where("id = ?", skuID).Update("stock", gorm.Expr("stock - ?", product.Count))
		if result.Error != nil || result.RowsAffected <= 0 {
			tx.Rollback()
			zap.L().Error("扣减库存失败", zap.Int64("skuID", skuID))
			return errors.New(fmt.Sprintf("扣减ID为%s商品库存失败", product.SkuID))
		}

		// 累加总金额和购买商品的数量
		count := decimal.NewFromFloat(float64(product.Count))
		price := decimal.NewFromFloat(sku.Price)
		totalMoney = totalMoney.Add(price.Mul(count))
		totalNum += product.Count

		// 生成订单明细
		orderItem := new(pojo.OrderItem)
		orderItem.OrderID = orderNum
		orderItem.SkuID = sku.ID
		orderItem.SpuID = sku.SpuID
		orderItem.ProductName = sku.Title
		orderItem.ProductPrice = sku.Price
		orderItem.ProductQuantity = product.Count
		// 计算该商品的总金额
		itemPrice := decimal.NewFromFloat(sku.Price)
		itemCount := decimal.NewFromFloat(float64(product.Count))
		itemTotalMoney := itemPrice.Mul(itemCount)
		orderItem.ProductTotalMoney = itemTotalMoney.InexactFloat64()
		// 获取商品图片
		skuPic := &pojo.SkuPic{}
		result = tx.Model(skuPic).Where("sku_id = ? and is_default = 1", sku.ID).First(skuPic)
		if result.Error != nil || result.RowsAffected <= 0 {
			tx.Rollback()
			zap.L().Error("获取商品默认图片失败", zap.Int64("skuID", skuID))
			return errors.New(fmt.Sprintf("获取ID为%s商品图片失败", product.SkuID))
		}
		orderItem.ProductPic = skuPic.PicUrl
		// 订单明细入库
		result = tx.Create(orderItem)
		if result.Error != nil || result.RowsAffected <= 0 {
			tx.Rollback()
			zap.L().Error("订单明细入库失败", zap.Int64("skuID", skuID))
			return errors.New("订单明细入库")
		}
	}

	// 设置商品总价格
	order.TotalMoney = totalMoney.InexactFloat64()
	// 运费为18元
	freight := decimal.NewFromFloat(18)
	// 计算出订单应付款 = 总金额 + 运费
	payMoney := totalMoney.Add(freight)
	// 设置应付款
	order.PayMoney = payMoney.InexactFloat64()
	// 设置购买商品数量
	order.TotalNum = uint8(totalNum)
	// 订单入库
	result := tx.Create(order)
	if result.Error != nil || result.RowsAffected <= 0 {
		tx.Rollback()
		zap.L().Error("订单入库失败", zap.Int64("skuID", orderNum), zap.Error(result.Error))
		return errors.New("订单入库失败")
	}
	tx.Commit()
	return nil
}

// SelectAllOrder 返回用户所有订单主表信息
func SelectAllOrder(uid int64) ([]*pojo.Order, error) {
	data := make([]*pojo.Order, 0)
	if err := db.Model(&pojo.Order{}).Where("user_id = ?", uid).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// SelectOneOrderByUIDAndOrderNum 根据用户ID和订单号查询用户订单信息
func SelectOneOrderByUIDAndOrderNum(uid, orderNum int64) (*pojo.Order, error) {
	order := new(pojo.Order)
	result := db.Model(&pojo.Order{}).Where("id = ? and user_id = ?", orderNum, uid).First(order)
	if result.Error != nil || result.RowsAffected <= 0 {
		zap.L().Error("根据用户ID和订单号查询用户订单信息失败", zap.Error(result.Error), zap.Int64("rowAffected", result.RowsAffected))
		return nil, errors.New("根据用户ID和订单号查询用户订单信息失败")
	}
	return order, nil
}

// SelectOneOrderItem 返回一条订单的明细信息
func SelectOneOrderItem(id int64) ([]*pojo.OrderItem, error) {
	data := make([]*pojo.OrderItem, 0)
	result := db.Model(&pojo.OrderItem{}).Where("order_id = ?", id).Find(&data)
	if result.Error != nil || result.RowsAffected <= 0 {
		// 查询过程中出现异常或者查询到的行数为0
		zap.L().Error("获取订单明细失败", zap.Error(result.Error), zap.Int64("rowAffected", result.RowsAffected))
		return nil, errors.New("订单明细不存在")
	}
	return data, nil
}

// SelectOrderOrderStatus 返回订单状态
func SelectOrderOrderStatus(id int64) (uint8, error) {
	order := &pojo.Order{ID: id}
	result := db.Model(order).First(order)
	if result.Error != nil || result.RowsAffected <= 0 {
		zap.L().Error("获取订单状态失败", zap.Error(result.Error), zap.Int64("rowAffected", result.RowsAffected))
		return 0, errors.New("获取订单状态失败")
	}
	return order.OrderStatus, nil
}

// UpdateOrderOrderStatus 修改订单支付状态
func UpdateOrderOrderStatus(id int64, orderStatus uint8) error {
	result := db.Model(&pojo.Order{ID: id}).Update("order_status", orderStatus)
	if result.Error != nil || result.RowsAffected <= 0 {
		zap.L().Error("修改订单状态失败", zap.Error(result.Error), zap.Int64("rowAffected", result.RowsAffected))
		return errors.New("修改订单状态失败")
	}
	return nil
}

// RollbackOrderStock 回滚订单超时未支付的商品库存
func RollbackOrderStock(id int64) error {
	// 查询订单所包含的所有商品明细
	tx := db.Begin()
	items := make([]*pojo.OrderItem, 0)
	result := tx.Model(&pojo.OrderItem{}).Where("order_id = ?", id).Find(&items)
	if result.Error != nil || result.RowsAffected <= 0 {
		tx.Rollback()
		zap.L().Error("查询订单所包含的所有商品明细失败", zap.Error(result.Error), zap.Int64("rowAffected", result.RowsAffected))
		return errors.New("查询订单所包含的所有商品明细失败")
	}

	// 遍历
	for _, item := range items {
		// 判断商品是否已经下架，如果已经下架，无需回滚库存
		sku := new(pojo.Sku)
		err := tx.Model(&pojo.Sku{ID: item.SkuID}).First(sku).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 没有找到商品记录，可能是商品已经下架，并删除了数据库中的记录
				continue
			}
			tx.Rollback()
			return err
		} else {
			// 异常为空，但是商品已经下架，数据库记录未删除
			if sku.Valid == 0 {
				continue
			}
		}

		// 回滚库存(使用Version字段解决并发下的更新问题)
		result = tx.Debug().Model(&pojo.Sku{ID: item.SkuID}).Update("stock", gorm.Expr("stock + ?", item.ProductQuantity))
		if result.Error != nil || result.RowsAffected <= 0 {
			tx.Rollback()
			zap.L().Error("回滚商品库存失败", zap.Error(result.Error), zap.Int64("rowAffected", result.RowsAffected))
			return errors.New("回滚商品库存失败")
		}
	}
	// 库存全部都回滚成功
	tx.Commit()
	return nil
}

// DelOrderAndItems 删除订单主表信息和对应的订单明细信息
func DelOrderAndItems(id int64) error {
	tx := db.Begin()

	// 删除订单主表信息
	order := &pojo.Order{ID: id}
	result := tx.Delete(order)
	if result.Error != nil || result.RowsAffected <= 0 {
		tx.Rollback()
		zap.L().Error("删除订单主表信息失败", zap.Error(result.Error), zap.Int64("rowAffected", result.RowsAffected))
		return errors.New("删除订单主表信息失败")
	}

	// 删除订单明细
	result = tx.Where("order_id = ?", id).Delete(&pojo.OrderItem{})
	if result.Error != nil || result.RowsAffected <= 0 {
		tx.Rollback()
		zap.L().Error("删除订单明细失败", zap.Error(result.Error), zap.Int64("rowAffected", result.RowsAffected))
		return errors.New("删除订单明细失败")
	}

	tx.Commit()
	return nil
}
