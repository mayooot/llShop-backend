package logic

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/models/dto"
	"shop-backend/models/vo"
	"shop-backend/utils/build"
	"shop-backend/utils/gen"
)

// CreatePreSubmitOrder 创建预提交订单
// 1. 生成全局唯一订单号
// 2. 判断预提交订单中的商品是否已经下架
// 3. 判断预提交订单中的商品购买数量是否大于库存
// 4. 计算出订单应付款 = 总金额 + 运费
func CreatePreSubmitOrder(preSubmitOrder *dto.PreSubmitOrder, uid int64) (*vo.OrderVO, error) {
	// 要返回的订单展示对象
	orderVO := new(vo.OrderVO)
	// 生成全局唯一订单号
	orderVO.OrderNumber = gen.GenSnowflakeID()
	// 初始化总金额
	totalMoney := decimal.NewFromFloat(0)
	channel := make(chan *vo.CartProductVO, 1)
	defer close(channel)
	// 遍历预提交订单中的商品
	for _, cartProduct := range preSubmitOrder.CartProductList {
		// 校验商品是否下架
		// 校验商品购买数量是否大于库存
		cartPojo, sku, err := mysql.CheckOrderProduct(cartProduct, uid)
		if err != nil {
			zap.L().Error("商品已下架或购买数量超过库存", zap.Error(err))
			return nil, err
		}

		// 构建订单(购物车)商品展示对象
		cartProductVO := build.CreateCartProductVO(cartPojo, channel)
		// 添加到订单展示对象中
		orderVO.CartProductVOList = append(orderVO.CartProductVOList, cartProductVO)

		// 计算总金额
		// 数量
		count := decimal.NewFromFloat(float64(cartProduct.Count))
		// 单价
		price := decimal.NewFromFloat(sku.Price)
		// 累加到预提交订单总价格
		totalMoney = totalMoney.Add(price.Mul(count))
	}

	zap.L().Info("totalMoney", zap.String("totalMoney", totalMoney.String()))
	// 运费为18元
	freight := decimal.NewFromFloat(18)

	// 计算出订单应付款 = 总金额 + 运费
	payMoney := totalMoney.Add(freight)

	orderVO.TotalMoney = totalMoney.String()
	orderVO.Freight = freight.String()
	orderVO.PayMoney = payMoney.String()
	return orderVO, nil
}