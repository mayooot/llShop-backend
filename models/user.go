package models

import "time"

type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
	Avatar   string `db:"avatar"`
}

// SomeInfo 用户简略信息
type SomeInfo struct {
	Avatar   string `json:"avatar" db:"avatar"`
	Username string `json:"username" db:"username"`
	// todo 添加购物车数量对应的数据库tag
	CartNum int `json:"cartNum"`
}

// UserInfos 用户详细信息
type UserInfos struct {
	Id         int64     `json:"id" db:"user_id"`
	Username   string    `json:"username" db:"username"`
	Phone      string    `json:"phone" db:"phone"`
	Email      string    `json:"email" db:"email"`
	Avatar     string    `json:"avatar" db:"avatar"`
	Gender     int8      `json:"gender" db:"gender"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
}
