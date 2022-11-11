package logic

import (
	"github.com/shopspring/decimal"
	"shop-backend/dao/mysql"
	"shop-backend/utils/pay"
	"strconv"
)

// CreateAlipayOrder 根据订单号和用户ID查询用户订单金额，并调用支付宝进行支付
func CreateAlipayOrder(uid, orderNum int64) (string, error) {
	// 获取订单信息
	order, err := mysql.SelectOneOrderByUIDAndOrderNum(uid, orderNum)
	if err != nil {
		return "", err
	}

	// 调用支付宝支付接口
	payMoney := decimal.NewFromFloat(order.TotalMoney)
	payUrl := pay.AliPay(strconv.FormatInt(orderNum, 10), payMoney.String())
	return payUrl, nil
}
