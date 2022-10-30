package pojo

import "time"

type Cart struct {
	UserID        int64     `gorm:"column:user_id"`
	SkuID         int64     `gorm:"column:sku_id"`
	Specification string    `gorm:"column:specification"`
	Count         int       `gorm:"column:count"`
	Selected      int8      `gorm:"column:selected"`
	CreatedTime   time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime   time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (Cart) TableName() string {
	return "oms_cart"
}
