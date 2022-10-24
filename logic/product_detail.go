package logic

import (
	"shop-backend/dao/mysql"
	"shop-backend/models/vo"
)

// GetProductDetail 获取商品详情
func GetProductDetail(spuID int64) (*vo.ProductDetailVO, error) {
	detail := new(vo.ProductDetailVO)
	// 获取spu
	spu, err := mysql.SelectSpuByID(spuID)
	if err != nil {
		return nil, err
	}
	// 封装spuVO
	spuVO := &vo.SpuVO{
		ID:                   spu.ID,
		Sale:                 spu.Sale,
		SubTitle:             spu.SubTitle,
		ProductSpecification: spu.ProductSpecification,
	}
	detail.Spu = spuVO

	// 获取categories
	categories, err := mysql.SelectSpuCategoryByCID(spu.CID1, spu.CID2)
	if err != nil {
		return nil, err
	}
	// 封装categoriesVO
	categoriesVO := make([]*vo.CategoryVO, 0)
	for _, cate := range categories {
		categoriesVO = append(categoriesVO, &vo.CategoryVO{
			ID:   cate.ID,
			Name: cate.Name,
		})
	}
	detail.Categories = categoriesVO

	// 获取skuList
	skuList, err := mysql.SelectSkuListBySpuID(spuID)
	// 封装skuListVO
	skuListVO := make([]*vo.SkuVO, 0)
	for _, sku := range skuList {
		// 获取sku的所有商品图片
		skuPicList, err := mysql.SelectSkuPicBySkuID(sku.ID)
		if err != nil {
			return nil, err
		}
		// 封装skuPicListVO
		skuPicListVO := make([]*vo.SkuPicVO, 0)
		for _, pic := range skuPicList {
			skuPicListVO = append(skuPicListVO, &vo.SkuPicVO{
				ID:        pic.ID,
				SkuID:     pic.SkuID,
				PicUrl:    pic.PicUrl,
				IsDefault: pic.IsDefault,
			})
		}
		skuListVO = append(skuListVO, &vo.SkuVO{
			ID:                      sku.ID,
			SpuID:                   sku.SpuID,
			Title:                   sku.Title,
			Price:                   sku.Price,
			Unit:                    sku.Unit,
			Stock:                   sku.Stock,
			ProductSkuSpecification: sku.ProductSkuSpecification,
			SkuPicList:              skuPicListVO,
		})
	}
	detail.SkuList = skuListVO
	return detail, nil
}
