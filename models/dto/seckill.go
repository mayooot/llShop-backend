package dto

type SecKillProduct struct {
	// 秒杀商品ID
	SkuID string `json:"skuID" binding:"required"`
}

// SecKillMQ 发送到MQ中的结构体
type SecKillMQ struct {
	SkuID int64
	UID   int64
}
