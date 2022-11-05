package pojo

// PCDDic 省市区字典表
type PCDDic struct {
	// ID
	ID int `gorm:"column:id"`
	// 上级编号
	ParentID int `gorm:"column:parent_id"`
	// 名称
	Name string `gorm:"column:name"`
}

func (PCDDic) TableName() string {
	return "ums_pcd_dic"
}
