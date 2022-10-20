package pojo

import "time"

// ProductAttribute 商品属性表
type ProductAttribute struct {
	ID          int64     `gorm:"column:id"`
	Type        uint8     `gorm:"column:type"`
	ParentID    int64     `gorm:"column:parent_id"`
	Name        string    `gorm:"column:name"`
	Sort        uint8     `gorm:"column:sort"`
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (ProductAttribute) TableName() string {
	return "pms_product_attribute"
}
