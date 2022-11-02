package redis

import (
	"go.uber.org/zap"
	"time"
)

var (
	productSpuSpecification = "product:spu:specification"
	specificationLivingTime = time.Hour * 24 * 15
)

// SetSpuSpecification 缓存spu规格信息
func SetSpuSpecification(specifications string) error {
	if err := rdb.Set(productSpuSpecification, specifications, specificationLivingTime).Err(); err != nil {
		zap.L().Error("将spu规格信息缓存进Redis失败", zap.Error(err))
		return err
	}
	return nil
}

// GetSpuSpecification 获取缓存中的商品分类信息
func GetSpuSpecification() (string, bool) {
	result, err := rdb.Get(productSpuSpecification).Result()
	if err != nil {
		zap.L().Info("获取缓存中的spu规格信息失败")
		return "", false
	}
	return result, true
}
