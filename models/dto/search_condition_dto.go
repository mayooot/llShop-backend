package dto

import "strconv"

// SearchCondition 封装搜索条件的请求体
type SearchCondition struct {
	// 品牌ID
	// todo 未使用，可以通过品牌ID查询商品
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

// NewCondition 初始化搜索条件，并指定分页默认值
func NewCondition() *SearchCondition {
	return &SearchCondition{
		PageNo:   strconv.Itoa(1),
		PageSize: strconv.Itoa(20),
	}
}
