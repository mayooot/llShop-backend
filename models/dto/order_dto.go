package dto

// PreSubmitOrder 封装用户预提交订单的属性
type PreSubmitOrder struct {
	// 勾选的商品集合
	CartProductList []*CartProduct `json:"cartProductList" binding:"required"`
}

// Order 封装用户提交的订单属性
type Order struct {
	// 订单号
	OrderNumber string `json:"orderNumber" binding:"required"`
	// 订单商品集合
	CartProductList []*CartProduct `json:"cartProductList" binding:"required"`
	// 收货人姓名
	ReceiverName string `json:"receiverName" binding:"required"`
	// 收货人手机号
	ReceiverPhone string `json:"receiverPhone" binding:"required"`
	// 收货人地址
	ReceiverAddress string `json:"receiverAddress" binding:"required"`
}
