package logic

import (
	"shop-backend/dao/mysql"
	"shop-backend/models/vo"
)

// GetAllAttribute 获取商品二级分类的所有属性
func GetAllAttribute(categoryID int64) ([]*vo.AttributeVO, error) {
	// category和attribute对应表集合
	caRels, err := mysql.SelectAttrIDeByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}
	// 商品属性ID Set集合
	attrIDSet := make(map[int64]bool)
	for _, caRel := range caRels {
		attrIDSet[caRel.ProductAttributeID] = true
	}

	result := make([]*vo.AttributeVO, 0)
	// 获取所有商品属性
	attrs, err := mysql.SelectAllAttribute()
	for _, attr := range attrs {
		if attrIDSet[attr.ID] {
			// 如果Set集合中有这一条记录
			attribute := &vo.AttributeVO{
				KeyID:   attr.ID,
				KeyName: attr.Name,
			}
			result = append(result, attribute)
			attrVals := make([]*vo.AttributeValueVO, 0)
			for _, reAttr := range attrs {
				if reAttr.ParentID == attr.ID {
					attrVals = append(attrVals, &vo.AttributeValueVO{
						ValueID:   reAttr.ID,
						ValueName: reAttr.Name,
					})
				}
			}
			attribute.AttributeValues = attrVals
		}
	}

	return result, nil
}
