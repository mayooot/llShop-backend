package redis

import (
	"go.uber.org/zap"
	"shop-backend/utils/concatstr"
	"strconv"
	"time"
)

var (
	productSpuSpecificationPrefix = "product:spu:specification:"
	specificationLivingTime       = time.Hour * 24 * 15
)

// SetSpuSpecification 缓存spu规格信息
func SetSpuSpecification(skuID int64, specifications string) error {
	key := concatstr.ConcatString(productSpuSpecificationPrefix, strconv.FormatInt(skuID, 10))
	if err := rdb.Set(key, specifications, specificationLivingTime).Err(); err != nil {
		zap.L().Error("将spu规格信息缓存进Redis失败", zap.Error(err))
		return err
	}
	return nil
}

// GetSpuSpecification 获取缓存中的商品分类信息
func GetSpuSpecification(skuID int64) (string, bool) {
	key := concatstr.ConcatString(productSpuSpecificationPrefix, strconv.FormatInt(skuID, 10))
	result, err := rdb.Get(key).Result()
	if err != nil {
		zap.L().Info("获取缓存中的spu规格信息失败")
		return "", false
	}
	return result, true
}
