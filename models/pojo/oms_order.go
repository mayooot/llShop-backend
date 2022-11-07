package pojo

import (
	"gorm.io/plugin/optimisticlock"
	"time"
)

// Order 订单主表
type Order struct {
	// 雪花算法生成的主键ID
	ID int64 `gorm:"column:id"`
	// 用户ID
	UserID int64 `gorm:"column:user_id"`
	// 订单总金额合计
	TotalMoney float64 `gorm:"column:total_money"`
	// 实付金额合计
	PayMoney float64 `gorm:"column:pay_money"`
	// 购买商品数量
	TotalNum uint8 `gorm:"column:total_num"`
	// 支付方式：0->在线支付；1->货到付款
	PayType uint8 `gorm:"column:pay_type"`
	// 0->待付款；1->待发货；2->已发货；3->已完成；4->已关闭；5->超时
	OrderStatus uint8 `gorm:"column:order_status"`
	// 支付状态：0->未支付；1->支付成功；2->支付失败
	PayStatus uint8 `gorm:"column:pay_status"`
	// 收件人名称
	ReceiverName string `gorm:"column:receiver_name"`
	// 收件人电话
	ReceiverPhone string `gorm:"column:receiver_phone"`
	// 收件人地址
	ReceiverAddress string `gorm:"column:receiver_address"`
	// 支付时间
	PayTime time.Time `gorm:"column:pay_time"`
	// 订单过期时间
	ExpirationTime time.Time `gorm:"column:expiration_time"`
	// 创建时间
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	// 修改时间
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
	// 版本控制
	Version optimisticlock.Version `gorm:"column:version"`
}

func (Order) TableName() string {
	return "oms_order"
}

// OrderItem 订单明细表
type OrderItem struct {
	// 主键ID
	ID int64 `gorm:"column:id"`
	// 订单ID(对应订单表主键ID)
	OrderID int64 `gorm:"column:order_id"`
	// 商品spuID(对应商品spu表主键ID)
	SpuID int64 `gorm:"column:spu_id"`
	// 商品skuID(对应商品sku表主键ID)
	SkuID int64 `gorm:"column:sku_id"`
	// 商品图片
	ProductPic string `gorm:"column:product_pic"`
	// 商品名称
	ProductName string `gorm:"column:product_name"`
	// 销售价格
	ProductPrice float64 `gorm:"column:product_price"`
	// 购买数量
	ProductQuantity int `gorm:"column:product_quantity"`
	// 创建时间
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	// 修改时间
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (OrderItem) TableName() string {
	return "oms_order_item"
}

// OrderPayLog 订单支付记录表
type OrderPayLog struct {
	// 主键ID
	ID int64 `gorm:"column:id"`
	// 用户ID
	UserID int64 `gorm:"column:user_id"`
	// 订单ID(对应订单表主键ID)
	OrderID int64 `gorm:"column:order_id"`
	// 订单编号
	OrderNum int64 `gorm:"column:order_num"`
	// 第三方支付订单交易号
	PayTradeNum int64 `gorm:"column:pay_trade_num"`
	// 支付方式：1->支付宝支付；2->微信支付
	PayWay uint8 `gorm:"column:pay_way"`
	// 支付状态：1->支付成功；2->支付失败
	PayStatus uint8 `gorm:"column:pay_status"`
	// 支付金额
	PayAmount float64 `gorm:"column:pay_amount"`
	// 创建时间
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	// 修改时间
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (OrderItem) OrderPayLog() string {
	return "oms_pay_log"
}
