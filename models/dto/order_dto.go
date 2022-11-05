package dto

// PreSubmitOrder 封装用户预提交订单的属性
type PreSubmitOrder struct {
	// 勾选的商品集合
	CartProductList []*CartProduct `json:"cartProductList" binding:"required"`
}
