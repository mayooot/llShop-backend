package models

type User struct {
	UserID   int64  `db:"user_id"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}
