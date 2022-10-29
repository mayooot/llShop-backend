package pojo

import "time"

type Cart struct {
	UserID      int64     `gorm:"column:user_id"`
	SkuID       int64     `gorm:"column:sku_id"`
	Count       int       `gorm:"column:count"`
	Selected    int8      `gorm:"column:selected"`
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
}

func (Cart) TableName() string {
	return "oms_cart"
}