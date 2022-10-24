package pojo

import "time"

// Spu 商品spu表
type Spu struct {
	// 主键
	ID int64 `gorm:"column:id"`
	// 品牌ID(对应品牌表主键ID)
	BrandId int64 `gorm:"column:brand_id"`
	// 一级分类ID(对应商品分类表主键ID)
	CID1 int64 `gorm:"column:cid1"`
	// 二级分类ID(对应商品分类表主键ID)
	CID2 int64 `gorm:"column:cid2"`
	// 商品总销量
	Sale int `gorm:"column:sale"`
	// 上架状态：0->下架；1->上架
	PublishStatus uint8 `gorm:"column:publish_status"`
	// 审核状态：0->未审核；1->审核通过
	VerifyStatus uint8 `gorm:"column:verify_status"`
	// 是否有效，0->已删除；1->有效
	Valid uint8 `gorm:"column:valid"`
	// 商品名称
	Name string `gorm:"column:name"`
	// 副标题
	SubTitle string `gorm:"column:sub_title"`
	// 商品规格(json格式，用于商品详情页展示商品所有规格)
	ProductSpecification string `gorm:"column:product_specification"`
	// 商品默认图片URL
	DefaultPicUrl string `gorm:"column:default_pic_url"`
	// 商品默认价格
	DefaultPrice float64 `gorm:"column:default_price"`
	// 创建时间
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime"`
	// 修改时间
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime"`
}

func (Spu) TableName() string {
	return "pms_spu"
}
