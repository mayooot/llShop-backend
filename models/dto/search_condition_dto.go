package dto

// SearchCondition 封装搜索条件的请求体
type SearchCondition struct {
	// 品牌ID
	BrandId string `json:"brandId"`
	// 二级分类ID
	ProductCategoryId string `json:"productCategoryId"`
	// 搜索条件
	Keyword string `json:"keyword"`
	// 商品属性ID数组
	ProductAttributeIds []int64 `json:"productAttributeIds"`
	// 排序: 1->默认; 2->销量
	Sort string `json:"sort"`
	// 页码(从1开始),默认为1
	PageNo string `json:"pageNo"`
	// 页长,默认为20
	PageSize string `json:"pageSize"`
}
