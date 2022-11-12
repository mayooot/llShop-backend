package pojo

import "gorm.io/plugin/optimisticlock"

type SecKillSku struct {
	ID            int64                  `gorm:"column:id" json:"id"`
	Price         float64                `gorm:"column:price" json:"price"`
	Stock         int                    `gorm:"column:stock" json:"stock"`
	Sale          int                    `gorm:"column:sale" json:"sale"`
	Title         string                 `gorm:"column:title" json:"title"`
	Specification string                 `gorm:"column:specification" json:"specification"`
	PicUrl        string                 `gorm:"column:pic_url" json:"picUrl"`
	Version       optimisticlock.Version `gorm:"column:version" json:"-"`
}

func (SecKillSku) TableName() string {
	return "pms_seckill_sku"
}
