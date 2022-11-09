package pojo

import (
	"gorm.io/plugin/optimisticlock"
	"time"
)

// Sku 商品sku表
type Sku struct {
	// 主键
	ID int64 `gorm:"column:id"`
	// 商品spuID(对应商品spu表主键ID)
	SpuID   int64                  `gorm:"column:spu_id"`
	Version optimisticlock.Version `gorm:"column:version"`
	// 销量
	Sale int `gorm:"column:sale"`
	// 库存
	Stock int `gorm:"column:stock"`
	// 默认规格：0->不是；1->是
	IsDefault uint8 `gorm:"column:is_default"`
	// 是否有效，0->无效；1->有效
	Valid uint8 `gorm:"column:valid"`
	// 商品标题
	Title string `gorm:"column:title"`
	// 商品单位
	Unit string `gorm:"column:unit"`
	// spu中商品规格的对应下标组合
	Indexes string `gorm:"column:indexes"`
	// 商品sku规格，json格式
	ProductSkuSpecification string `gorm:"column:product_sku_specification"`
	// 价格
	Price float64 `gorm:"column:price"`
	// 创建时间
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	// 修改时间
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (Sku) TableName() string {
	return "pms_sku"
}
