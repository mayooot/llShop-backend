package pojo

import "time"

// SkuPic 商品sku图片表
type SkuPic struct {
	// 主键
	ID int64 `gorm:"column:id"`
	// 商品skuID(对应商品sku表主键ID)
	SkuID int64 `gorm:"column:sku_id"`
	// 图片URL
	PicUrl string `gorm:"column:pic_url"`
	// 默认展示：0->否；1->是
	IsDefault uint8 `gorm:"column:is_default"`
	// 创建时间
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	// 修改时间
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (SkuPic) TableName() string {
	return "pms_sku_pic"
}
