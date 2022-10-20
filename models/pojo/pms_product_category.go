package pojo

import "time"

// ProductCategory 商品分类表
type ProductCategory struct {
	ID           int64     `gorm:"column:id"`
	ParentID     int64     `gorm:"column:parent_id"`
	Name         string    `gorm:"column:name"`
	Abbreviation string    `gorm:"column:abbreviation"`
	Level        uint8     `gorm:"column:level"`
	ShowStatus   uint8     `gorm:"column:show_status"`
	Icon         string    `gorm:"column:icon"`
	Sort         uint8     `gorm:"column:sort"`
	CreatedTime  time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime  time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (ProductCategory) TableName() string {
	return "pms_product_category"
}
