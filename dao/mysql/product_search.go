package mysql

import (
	"go.uber.org/zap"
	"shop-backend/models/dto"
	"shop-backend/models/vo"
	"shop-backend/utils/concatstr"
	"strconv"
)

// SelectProductSearchCondition 根据条件查询商品
func SelectProductSearchCondition(condition *dto.SearchCondition) ([]*vo.Product, error) {
	data := make([]*vo.Product, 0)
	db := db.Model(&vo.Product{})
	db.Select("pms_sku.id, title AS name, pms_sku.sale, price AS defaultPrice, pms_spu.default_pic_url AS defaultPicUrl")
	db.Joins("LEFT JOIN pms_sku_pic ON pms_sku_pic.sku_id = pms_sku.id")
	db.Joins("LEFT JOIN pms_spu ON pms_sku.spu_id = pms_spu.id")
	db.Joins("LEFT JOIN pms_product_attribute_rel ON pms_product_attribute_rel.spu_id = pms_spu.id")
	db.Joins("LEFT JOIN pms_product_attribute ON pms_product_attribute.id = pms_product_attribute_rel.product_attribute_id")

	if condition.Keyword != "" {
		// sku表 搜索关键字不为空
		db.Where("title like ?", concatstr.ConcatString("%", condition.Keyword, "%"))
	}
	if condition.BrandId != "" {
		brandId, err := strconv.ParseInt(condition.BrandId, 10, 64)
		if err != nil {
			zap.L().Error("BrandId转换为整型失败", zap.Error(err))
			return nil, err
		}
		// spu表 品牌ID不为空
		db.Where("brand_id = ?", brandId)
	}

	if condition.ProductCategoryId != "" {
		productCategoryId, err := strconv.ParseInt(condition.ProductCategoryId, 10, 64)
		if err != nil {
			zap.L().Error("ProductCategoryId转换为整型失败", zap.Error(err))
			return nil, err
		}
		// spu表 二级分类ID不为空
		db.Where("cid2 = ?", productCategoryId)
	}

	if len(condition.ProductAttributeIds) != 0 {
		// pms_product_attribute_rel表 商品属性集合不为空
		db.Where("pms_product_attribute_rel.product_attribute_id IN ? ", condition.ProductAttributeIds)
	}

	if condition.Sort != "" {
		sort, err := strconv.ParseUint(condition.Sort, 10, 8)
		if err != nil {
			zap.L().Error("sort转换为整型失败", zap.Error(err))
			return nil, err
		}
		// 排序ID不为空
		if sort == 1 {
			// 按照创建时间排序
			db.Order("pms_product_attribute.created_time desc")
		} else if sort == 2 {
			// 按照销量排序
			db.Order("pms_spu.sale desc")
		}
	}

	if condition.PageSize != "" && condition.PageNo != "" {
		limit, err := strconv.Atoi(condition.PageSize)
		if err != nil {
			zap.L().Error("PageSize转换为整型失败", zap.Error(err))
			return nil, err
		}

		pageNo, err := strconv.Atoi(condition.PageNo)
		if err != nil {
			zap.L().Error("PageNo转换为整型失败", zap.Error(err))
			return nil, err
		}

		// 分页
		db.Limit(limit).Offset(pageNo)
	}
	result := db.Debug().Find(&data)
	if result.Error != nil {
		zap.L().Error("使用搜索条件查询数据库失败", zap.Error(result.Error))
		return nil, result.Error
	}

	return data, nil
}
