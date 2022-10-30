package dto

// CartProduct 封装用户要添加进购物车的商品属性对象
type CartProduct struct {
	// 商品skuID
	SkuID string `json:"skuID"`
	// 添加数量
	Count string `json:"count"`
}

// CartProductSelected 封装购物车中商品勾选状态对象
type CartProductSelected struct {
	// 商品skuID
	SkuID string `json:"skuID"`
	// 勾选状态；0 -> 未勾选；1 -> 勾选
	Selected string `json:"selected"`
}
