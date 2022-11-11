package dto

type AliPay struct {
	OrderNum string `json:"orderNum" binding:"required"`
}
