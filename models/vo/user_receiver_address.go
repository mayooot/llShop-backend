package vo

// AddressVO 展示用户收货地址信息
type AddressVO struct {
	// 主键ID
	ID int64 `json:"id"`
	// 区县ID
	CountyID int `json:"countyID"`
	// 省
	Province string `json:"province"`
	// 市
	City string `json:"city"`
	// 区
	Region string `json:"region"`
	// 收货人
	UserName string `gorm:"column:user_name"`
	// 手机号
	PhoneNumber string `gorm:"column:phone_number"`
	// 是否是默认收货地址
	DefaultStatus uint8 `gorm:"column:default_status"`
	// 详细地址(街道、乡、村)
	DetailAddress string `gorm:"column:detail_address"`
}
