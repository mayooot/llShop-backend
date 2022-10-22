package vo

// Product 商品信息
type Product struct {
	// 主键ID
	ID int64 `json:"id" gorm:"column:id"`
	// 商品销量
	Sale int `json:"sale" gorm:"column:sale"`
	// 商品默认价格
	DefaultPrice float64 `json:"defaultPrice,string" gorm:"column:defaultPrice"`
	// 商品名称
	Name string `json:"name" gorm:"column:name"`
	// 商品默认图片URL
	DefaultPicUrl string `json:"defaultPicUrl" gorm:"column:defaultPicUrl"`
}

func (Product) TableName() string {
	return "pms_sku"
}
