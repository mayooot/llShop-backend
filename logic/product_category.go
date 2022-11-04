package logic

import (
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/dao/redis"
	"shop-backend/models/pojo"
	"shop-backend/models/vo"
)

var errChannel = make(chan error, 2)

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

	// 获取商品一级分类信息
	firstChannel := make(chan []*pojo.ProductCategory, 1)
	// 获取商品二级分类信息
	secondChannel := make(chan []*pojo.ProductCategory, 1)

	go getFirstCategories(firstChannel)
	go getSecondCategories(secondChannel)

	firstCategories := <-firstChannel
	secondCategories := <-secondChannel

	if len(errChannel) >= 1 {
		err := <-errChannel
		zap.L().Error("多协程获取一级分类，二级分类出错了", zap.Error(err))
		return nil, err
	}

	// 返回的结果切片
	var result = make([]vo.FirstProductCategoryVO, 0)
	for _, category := range firstCategories {
		// 遍历所有一级分类信息
		firstCategoryVO := vo.FirstProductCategoryVO{
			ID:         category.ID,
			Name:       category.Name,
			Level:      category.Level,
			ShowStatus: category.ShowStatus,
			Icon:       category.Icon,
		}
		// 筛选出从属该一级分类的二级分类信息
		for _, reCategory := range secondCategories {
			if reCategory.ParentID == category.ID {
				firstCategoryVO.SecProductCategoryList = append(firstCategoryVO.SecProductCategoryList, &vo.SecondProductCategoryVO{
					SecID:         reCategory.ID,
					SecName:       reCategory.Name,
					SecLevel:      reCategory.Level,
					SecShowStatus: reCategory.ShowStatus,
					SecIcon:       reCategory.Icon,
				})
			}
		}
		// 加入到返回结果切片
		result = append(result, firstCategoryVO)
	}
	// 将商品分类信息缓存进Redis
	redis.SetCategoryList(result)
	return result, nil
}

func getFirstCategories(channel chan []*pojo.ProductCategory) {
	// 获取商品一级分类信息
	firstCategories, err := mysql.SelectFirstCategory()
	if err != nil {
		errorChannel <- err
	}
	channel <- firstCategories
	defer close(channel)
}

func getSecondCategories(channel chan []*pojo.ProductCategory) {
	// 获取商品一级分类信息
	secondCategories, err := mysql.SelectSecondCategory()
	if err != nil {
		errorChannel <- err
	}
	channel <- secondCategories
	defer close(channel)
}
