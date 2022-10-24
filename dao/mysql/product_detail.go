package mysql

import (
	"go.uber.org/zap"
	"shop-backend/models/pojo"
)

// SelectSpuByID 使用spuID获取spu信息
func SelectSpuByID(spuID int64) (*pojo.Spu, error) {
	spu := new(pojo.Spu)
	if err := db.Model(&pojo.Spu{}).Where("id = ?", spuID).Debug().First(spu).Error; err != nil {
		zap.L().Error("使用spuID获取spu信息失败", zap.Error(err), zap.Int64("spuID", spuID))
		return nil, err
	}
	return spu, nil
}

// SelectSpuCategoryByCID 获取spu所有的分类信息(
func SelectSpuCategoryByCID(cid1, cid2 int64) ([]*pojo.ProductCategory, error) {
	categories := make([]*pojo.ProductCategory, 0)
	result := db.Model(&pojo.ProductCategory{}).Where("id in ?", []int64{cid1, cid2}).Debug().Find(&categories)
	if result.Error != nil {
		zap.L().Error("通过cid1和cid2获取spu所有的分类信息失败", zap.Error(result.Error), zap.Int64("cid1", cid1), zap.Int64("cid2", cid2))
		return nil, result.Error
	}
	return categories, nil
}

// SelectSkuListBySpuID 根据spuID获取skuList
func SelectSkuListBySpuID(spuID int64) ([]*pojo.Sku, error) {
	skuList := make([]*pojo.Sku, 0)
	result := db.Model(&pojo.Sku{}).Where("spu_id = ?", spuID).Where("is_default = 1").Debug().Find(&skuList)
	if result.Error != nil {
		zap.L().Error("根据spuID获取skuList", zap.Error(result.Error), zap.Int64("spuID", spuID))
		return nil, result.Error
	}
	return skuList, nil
}

// SelectSkuPicBySkuID 使用skuID获取sku商品图片
func SelectSkuPicBySkuID(skuID int64) ([]*pojo.SkuPic, error) {
	skuPicList := make([]*pojo.SkuPic, 0)
	if err := db.Model(&pojo.SkuPic{}).Where("sku_id = ?", skuID).Debug().Find(&skuPicList).Error; err != nil {
		zap.L().Error("使用skuID获取sku商品图片", zap.Error(err), zap.Int64("skuID", skuID))
		return nil, err
	}
	return skuPicList, nil
}
