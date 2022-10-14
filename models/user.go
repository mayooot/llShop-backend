package models

type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}

// SomeInfo 用户简略信息
type SomeInfo struct {
	Avatar   string `json:"avatar" db:"avatar"`
	Username string `json:"username" db:"username"`
	CartNum  int    `json:"cartNum"`
}
