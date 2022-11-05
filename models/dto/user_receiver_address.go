package dto

// ReceiverAddress 封装用户收货地址参数
type ReceiverAddress struct {
	// 区县ID
	CountyID int `json:"countyID" binding:"required"`
	// 是否是默认地址
	DefaultStatus int `json:"defaultStatus" binding:"required"`
	// 主键ID。新增用户收货地址不需要携带；更新用户收货地址/更新默认收货地址时需要携带。
	ID string `json:"id"`
	// 收货人姓名
	ReceiverName string `json:"receiverName" binding:"required"`
	// 收货人手机号
	ReceiverPhone string `json:"receiverPhone" binding:"required"`
	// 详细地址
	DetailAddress string `json:"detailAddress" binding:"required"`
}
