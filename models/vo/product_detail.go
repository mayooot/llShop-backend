package vo

// ProductDetailVO 商品详情信息
type ProductDetailVO struct {
	Spu        *SpuVO        `json:"spu"`
	Categories []*CategoryVO `json:"categories"`
	SkuList    []*SkuVO      `json:"skuList"`
}

type SpuVO struct {
	ID                   int64  `json:"id,string"`
	Sale                 int    `json:"sale"`
	SubTitle             string `json:"subTitle"`
	ProductSpecification string `json:"productSpecification"`
}

type CategoryVO struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type SkuVO struct {
	ID                      int64       `json:"id"`
	SpuID                   int64       `json:"spuID"`
	Title                   string      `json:"title"`
	Price                   float64     `json:"price"`
	Unit                    string      `json:"unit"`
	Stock                   int         `json:"stock"`
	ProductSkuSpecification string      `json:"productSkuSpecification"`
	SkuPicList              []*SkuPicVO `json:"skuPicList"`
}

type SkuPicVO struct {
	ID        int64  `json:"id"`
	SkuID     int64  `json:"skuID"`
	PicUrl    string `json:"picUrl"`
	IsDefault uint8  `json:"isDefault"`
}
