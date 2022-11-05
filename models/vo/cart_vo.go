package vo

import (
	"sync"
	"time"
)

// CartProductVO  购物车中单个商品展示对象
type CartProductVO struct {
	// 用于控制协程
	WG sync.WaitGroup `json:"-"`
	// 错误信息，用于多协程中错误的返回，前端不用渲染
	Err error `json:"-"`
	// 商品skuID
	SkuID int64 `json:"skuID,string"`
	// 商品标题
	Title string `json:"title"`
	// 数量
	Count int `json:"count"`
	// 商品勾选状态
	Selected int8 `json:"selected"`
	// 商品sku规格，json格式
	ProductSkuSpecification string `json:"productSkuSpecification"`
	// 商品默认图片URL
	DefaultPicUrl string `json:"defaultPicUrl"`
	// 上架状态：0->下架；1->上架
	PublishStatus uint8 `json:"publishStatus"`
	// 价格
	Price float64 `json:"price,string"`
	// 创建时间
	CreatedTime time.Time `json:"createdTime"`
}

// UserCartProductVOList 用户购物车商品列表
type UserCartProductVOList struct {
	// 用户ID
	UserID int64
	// 用户购物车商品列表
	CartList []*CartProductVO
}
