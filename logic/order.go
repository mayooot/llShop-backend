package logic

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/dao/redis"
	"shop-backend/models/dto"
	"shop-backend/models/vo"
	"shop-backend/rabbitmq"
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

	// 将订单编号设置进Redis，并设置5分钟的失效时间。实现提交订单幂等性和限流
	err := redis.SetOrderNumber(orderVO.OrderNumber)
	if err != nil {
		return nil, err
	}
	return orderVO, nil
}

// CreateSubmitOrder 创建订单
// 1. 生成订单
// 2. 校验库存
// 3. 扣减库存
// 4. 生成订单明细
// 5. 清空购物车
// 6. 失败后或者未支付回滚库存
func CreateSubmitOrder(orderDTO *dto.Order, uid, orderNum int64) error {
	// 生成订单 && 校验库存和商品状态 && 生成订单明细
	err := mysql.CreateOrderAndOrderItem(orderDTO, uid, orderNum)
	if err != nil {
		return err
	}

	// 异步清除购物车
	go rabbitmq.SendCartDelMess2MQ(orderDTO.CartProductList, uid)
	return nil
}
