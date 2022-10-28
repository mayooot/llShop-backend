package logic

import (
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/dao/redis"
	"shop-backend/models/vo"
)

// GetAllCategory 获取所有商品分类信息
func GetAllCategory() ([]vo.FirstProductCategoryVO, error) {
	// 先从缓存中获取商品分类信息
	data, exist := redis.GetCategoryList()
	if exist {
		zap.L().Info("使用商品分类信息缓存成功")
		// 如果商品分类信息存在，直接返回
		return data, nil
	} // 缓存不存在，从数据库中查询，并放入缓存
	zap.L().Info("使用商品分类信息缓存失败，查询数据库")

	// 获取所有商品分类信息
	categories, err := mysql.SelectAllCategory()
	if err != nil {
		return nil, err
	}

	// 返回的结果切片
	var result = make([]vo.FirstProductCategoryVO, 0)
	for _, category := range categories {
		// 遍历所有商品信息
		if category.ParentID == 0 {
			// 如果父级分类ID为0,说明为一级分类
			firstCategory := vo.FirstProductCategoryVO{
				ID:         category.ID,
				Name:       category.Name,
				Level:      category.Level,
				ShowStatus: category.ShowStatus,
				Icon:       category.Icon,
			}
			// 筛选出从属该一级分类的二级分类信息
			for _, reCategory := range categories {
				if reCategory.ParentID == category.ID {
					firstCategory.SecProductCategoryList = append(firstCategory.SecProductCategoryList, &vo.SecondProductCategoryVO{
						SecID:         reCategory.ID,
						SecName:       reCategory.Name,
						SecLevel:      reCategory.Level,
						SecShowStatus: reCategory.ShowStatus,
						SecIcon:       reCategory.Icon,
					})
				}
			}
			// 加入到返回结果切片
			result = append(result, firstCategory)
		}
	}

	// 将商品分类信息缓存进Redis
	redis.SetCategoryList(result)
	return result, nil
}
