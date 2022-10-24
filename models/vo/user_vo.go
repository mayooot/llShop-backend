package vo

import "time"

// SomeInfoVO 用户简略信息
type SomeInfoVO struct {
	Avatar   string `json:"avatar" `
	Username string `json:"username" `
	// todo 添加购物车数量对应的数据库tag
	CartNum int `json:"cartNum"`
}

// UserInfosVO 用户详细信息
type UserInfosVO struct {
	Id          int64     `json:"id,string" `
	Username    string    `json:"username" `
	Phone       string    `json:"phone" `
	Email       string    `json:"email" `
	Avatar      string    `json:"avatar" `
	Gender      uint8     `json:"gender" `
	CreatedTime time.Time `json:"createdTime" `
}
