package pojo

import "time"

type ReceiverAddress struct {
	// 主键ID
	ID int64 `gorm:"column:id"`
	// 用户ID
	UserID int64 `gorm:"column:user_id"`
	// 区县ID
	CountyID int `gorm:"column:county_id"`
	// 收货人
	UserName string `gorm:"column:user_name"`
	// 手机号
	PhoneNumber string `gorm:"column:phone_number"`
	// 是否是默认收货地址
	DefaultStatus uint8 `gorm:"column:default_status"`
	// 详细地址(街道、乡、村)
	DetailAddress string `gorm:"column:detail_address"`
	// 创建时间
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	// 修改时间
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (ReceiverAddress) TableName() string {
	return "ums_receiver_address"
}
