package redis

import (
	"encoding/json"
	"go.uber.org/zap"
	"shop-backend/models/vo"
	"time"
)

var (
	productCategoryPrefix = "product:category"
	categoryLivingTime    = time.Hour * 24 * 15
)

// SetCategoryList 缓存商品分类信息
func SetCategoryList(categories []vo.FirstProductCategoryVO) error {
	bytes, err := json.Marshal(categories)
	if err != nil {
		zap.L().Error("序列化商品分类信息失败")
		return err
	}
	if err := rdb.Set(productCategoryPrefix, string(bytes), categoryLivingTime).Err(); err != nil {
		zap.L().Error("将商品分类信息缓存进Redis失败", zap.Error(err))
		return err
	}
	return nil
}

// GetCategoryList 获取缓存中的商品分类信息
func GetCategoryList() ([]vo.FirstProductCategoryVO, bool) {
	result, err := rdb.Get(productCategoryPrefix).Result()
	if err != nil {
		zap.L().Info("获取缓存中的商品分类信息失败")
		return nil, false
	}

	var data = []byte(result)
	var categories = make([]vo.FirstProductCategoryVO, 0)
	if err := json.Unmarshal(data, &categories); err != nil {
		zap.L().Error("反序列化商品分类信息失败")
		return nil, false
	}

	return categories, true
}
