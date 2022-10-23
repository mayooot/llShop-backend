package vo

type Pageable interface {
	[]*Product
}
type Page[T Pageable] struct {
	// 起始页
	PageNo string `json:"pageNo"`
	// 每次大小
	PageSize string `json:"pageSize"`
	// 根据条件查询出来的总记录数(不分页)
	TotalPage string `json:"totalPage"`
	// 数据集合
	Data T `json:"data"`
}
