package dto

// CartProduct 封装用户要添加进购物车的商品属性
type CartProduct struct {
	// 商品skuID
	SkuID string `json:"skuID"`
	// 添加数量
	Count string `json:"count"`
}
