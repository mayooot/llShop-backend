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
