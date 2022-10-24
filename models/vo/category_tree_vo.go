package vo

// FirstProductCategoryVO 商品一级分类
type FirstProductCategoryVO struct {
	ID                     int64                      `json:"id,string"`
	Name                   string                     `json:"name"`
	Level                  uint8                      `json:"level"`
	ShowStatus             uint8                      `json:"showStatus"`
	Icon                   string                     `json:"icon"`
	SecProductCategoryList []*SecondProductCategoryVO `json:"secProductCategoryList"`
}

// SecondProductCategoryVO 商品二级分类
type SecondProductCategoryVO struct {
	SecID         int64  `json:"secID,string"`
	SecName       string `json:"secName"`
	SecLevel      uint8  `json:"secLevel"`
	SecShowStatus uint8  `json:"secShowStatus"`
	SecIcon       string `json:"secIcon"`
}
