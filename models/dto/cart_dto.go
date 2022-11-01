package dto

// CartProduct 封装用户要添加进购物车的商品属性对象
type CartProduct struct {
	// 商品skuID
	SkuID string `json:"skuID" binding:"required"`
	// 商品规格
	Specification string `json:"specification" binding:"required"`
	// 添加数量
	Count int `json:"count" binding:"required"`
}

// CartProductDel 封装用户删除购物车中商品时要发送的数据对象
type CartProductDel struct {
	// 商品skuID
	SkuID string `json:"skuID" binding:"required"`
	// 商品规格
	Specification string `json:"specification" binding:"required"`
}

// CartProductSelected 封装购物车中商品勾选状态对象
type CartProductSelected struct {
	// 商品skuID
	SkuID string `json:"skuID" binding:"required"`
	// 勾选状态；1 -> 勾选；2 -> 未勾选
	Selected string `json:"selected" binding:"required"`
	// 商品规格
	Specification string `json:"specification" binding:"required"`
}
