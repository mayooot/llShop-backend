package pojo

import "time"

// ProductCategoryAttributeRel 商品分类和属性中间表
type ProductCategoryAttributeRel struct {
	ID                 int64     `gorm:"column:id"`
	ProductCategoryID  int64     `gorm:"column:product_category_id"`
	ProductAttributeID int64     `gorm:"column:product_attribute_id"`
	CreatedTime        time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime        time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (ProductCategoryAttributeRel) TableName() string {
	return "pms_product_category_attribute_rel"
}
