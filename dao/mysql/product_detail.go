package mysql

import (
	"go.uber.org/zap"
	"shop-backend/models/pojo"
)

// SelectSpuBySkuID 使用skuID获取spu信息
func SelectSpuBySkuID(skuID int64) (*pojo.Spu, error) {
	spu := new(pojo.Spu)
	result := db.Model(&pojo.Sku{}).
		Select("pms_spu.id, pms_spu.brand_id, pms_spu.cid1, pms_spu.cid2, "+
			"pms_spu.sale, pms_spu.publish_status, pms_spu.verify_status, "+
			"pms_spu.valid, pms_spu.name, pms_spu.sub_title, "+
			"pms_spu.product_specification, pms_spu.default_pic_url, pms_spu.default_price, "+
			"pms_spu.created_time, pms_spu.updated_time").
		Joins("LEFT JOIN pms_spu ON pms_spu.id = pms_sku.spu_id").
		Where("pms_sku.id = ?", skuID).
		First(spu)
	if result.Error != nil {
		zap.L().Error("使用skuID获取spu信息失败", zap.Error(result.Error), zap.Int64("skuID", skuID))
		return nil, result.Error
	}
	return spu, nil
}

// SelectSpuCategoryByCID 获取spu所有的分类信息(
func SelectSpuCategoryByCID(cid1, cid2 int64) ([]*pojo.ProductCategory, error) {
	categories := make([]*pojo.ProductCategory, 0)
	result := db.Model(&pojo.ProductCategory{}).Where("id IN ?", []int64{cid1, cid2}).Find(&categories)
	if result.Error != nil {
		zap.L().Error("通过cid1和cid2获取spu所有的分类信息失败", zap.Error(result.Error), zap.Int64("cid1", cid1), zap.Int64("cid2", cid2))
		return nil, result.Error
	}
	return categories, nil
}

// SelectSkuListBySpuID 根据spuID获取skuList
func SelectSkuListBySpuID(spuID int64) ([]*pojo.Sku, error) {
	skuList := make([]*pojo.Sku, 0)
	result := db.Model(&pojo.Sku{}).Where("spu_id = ?", spuID).Where("is_default = 1").Find(&skuList)
	if result.Error != nil {
		zap.L().Error("根据spuID获取skuList", zap.Error(result.Error), zap.Int64("spuID", spuID))
		return nil, result.Error
	}
	return skuList, nil
}

// SelectSkuPicBySkuID 使用skuID获取sku商品图片
func SelectSkuPicBySkuID(skuID int64) ([]*pojo.SkuPic, error) {
	skuPicList := make([]*pojo.SkuPic, 0)
	if err := db.Model(&pojo.SkuPic{}).Where("sku_id = ?", skuID).Find(&skuPicList).Error; err != nil {
		zap.L().Error("使用skuID获取sku商品图片", zap.Error(err), zap.Int64("skuID", skuID))
		return nil, err
	}
	return skuPicList, nil
}
