package pojo

import "time"

// UmsUser 用户信息表
type UmsUser struct {
	ID          int64     `gorm:"column:user_id"`
	Username    string    `gorm:"column:username"`
	Password    string    `gorm:"column:password"`
	Email       string    `gorm:"column:email"`
	Phone       string    `gorm:"column:phone"`
	Avatar      string    `gorm:"column:avatar"`
	Gender      uint8     `gorm:"column:gender"`
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (UmsUser) TableName() string {
	return "ums_user"
}
