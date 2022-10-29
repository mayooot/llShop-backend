package mysql

import (
	"errors"
	"go.uber.org/zap"
	"shop-backend/models/dto"
	"shop-backend/models/vo"
	"shop-backend/utils/concatstr"
	"strconv"
	"strings"
)

var MAXRecord = 100
var ErrorExceedMaxRecord = errors.New("超过单次查询最大记录条数")

// BaseSearchCondition 根据条件查询商品，needLimit：是否要分页查询
func BaseSearchCondition(condition *dto.SearchCondition, needLimit bool) ([]*vo.ProductVO, int, error) {
	data := make([]*vo.ProductVO, 0)
	// 绑定db对应的表为pms_sku
	db := db.Model(&vo.ProductVO{})
	db.Select("pms_sku.id, pms_sku.title AS name, pms_sku.sale, pms_sku.price AS defaultPrice, pms_spu.default_pic_url AS defaultPicUrl")
	db.Joins("LEFT JOIN pms_sku_pic ON pms_sku_pic.sku_id = pms_sku.id")
	db.Joins("LEFT JOIN pms_spu ON pms_sku.spu_id = pms_spu.id")
	db.Joins("LEFT JOIN pms_product_attribute_rel ON pms_product_attribute_rel.spu_id = pms_spu.id")
	db.Joins("LEFT JOIN pms_product_attribute ON pms_product_attribute.id = pms_product_attribute_rel.product_attribute_id")
	// 每个sku在pms_sku中都有多个规格，is_default取值为0和1，1代表默认规格
	db.Where("pms_sku.is_default = ?", 1)
	if strings.TrimSpace(condition.Keyword) != "" {
		// sku表 搜索关键字不为空
		db.Where("title like ?", concatstr.ConcatString("%", strings.TrimSpace(condition.Keyword), "%"))
	}

	if strings.TrimSpace(condition.BrandId) != "" {
		// spu表 品牌ID不为空
		brandId, err := strconv.ParseInt(strings.TrimSpace(condition.BrandId), 10, 64)
		if err != nil {
			zap.L().Error("BrandId转换为整型失败", zap.Error(err))
			return nil, 0, err
		}
		db.Where("brand_id = ?", brandId)
	}

	if strings.TrimSpace(condition.ProductCategoryId) != "" {
		// spu表 商品二级分类ID不为空
		productCategoryId, err := strconv.ParseInt(strings.TrimSpace(condition.ProductCategoryId), 10, 64)
		if err != nil {
			zap.L().Error("ProductCategoryId转换为整型失败", zap.Error(err))
			return nil, 0, err
		}
		db.Where("cid2 = ?", productCategoryId)
	}

	if len(condition.ProductAttributeIds) != 0 {
		// pms_product_attribute_rel表 商品属性集合不为空
		db.Where("pms_product_attribute_rel.product_attribute_id IN ? ", condition.ProductAttributeIds)
	}

	// 排序ID不为空
	sort, err := strconv.ParseUint(condition.Sort, 10, 8)
	if err != nil {
		zap.L().Error("sort转换为整型失败", zap.Error(err))
		return nil, 0, err
	}
	if sort == 1 {
		// 按照创建时间排序
		db.Order("pms_product_attribute.created_time desc")
	} else if sort == 2 {
		// 按照销量排序
		db.Order("pms_spu.sale desc")
	}

	if needLimit {
		pageSize, err := strconv.Atoi(condition.PageSize)
		if err != nil {
			zap.L().Error("PageSize转换为整型失败", zap.Error(err))
			return nil, 0, err
		}
		if pageSize > MAXRecord {
			zap.L().Error("超过单次查询最大记录条数", zap.Error(ErrorExceedMaxRecord))
			return nil, 0, ErrorExceedMaxRecord
		}

		pageNo, err := strconv.Atoi(condition.PageNo)
		if err != nil {
			zap.L().Error("PageNo转换为整型失败", zap.Error(err))
			return nil, 0, err
		}
		// 分页
		db.Limit(pageSize).Offset((pageNo - 1) * pageSize)
	}

	// 分组去重
	db.Group("pms_sku.id")
	result := db.Debug().Find(&data)
	if result.Error != nil {
		zap.L().Error("使用搜索条件查询数据库失败", zap.Error(result.Error))
		return nil, 0, result.Error
	}

	return data, len(data), nil
}
