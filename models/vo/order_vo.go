package vo

// OrderVO 预支付订单展示对象
type OrderVO struct {
	// 订单号
	OrderNumber int64 `json:"orderNumber,string"`
	// 商品总金额 = 所有商品的(商品单价 * 数量)
	TotalMoney string `json:"totalMoney"`
	// 订单总金额 = 商品总金额 + 运费
	PayMoney string `json:"payMoney"`
	// 运费
	Freight string `json:"freight"`
	// 预支付订单商品对象集合
	CartProductVOList []*CartProductVO `json:"cartProductVOList"`
}
