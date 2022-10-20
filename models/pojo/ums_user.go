package pojo

import "time"

// UmsUser 用户信息表
type UmsUser struct {
	ID         int64     `gorm:"column:user_id"`
	Username   string    `gorm:"column:username"`
	Password   string    `gorm:"column:password"`
	Email      string    `gorm:"column:email"`
	Phone      string    `gorm:"column:phone"`
	Avatar     string    `gorm:"column:avatar"`
	Gender     int8      `gorm:"column:gender"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

func (UmsUser) TableName() string {
	return "ums_user"
}
