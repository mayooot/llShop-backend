package po

import "time"

type UmsUser struct {
	ID         int64     `gorm:"column:user_id"`
	Username   string    `gorm:"column:username"`
	Password   string    `gorm:"column:password"`
	Email      string    `gorm:"column:email"`
	Phone      string    `gorm:"column:phone"`
	Avatar     string    `gorm:"column:avatar"`
	Gender     int8      `gorm:"column:gender"`
	CreateTime time.Time `gorm:"autoCreateTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime"`
}

func (u UmsUser) TableName() string {
	return "ums_user"
}
