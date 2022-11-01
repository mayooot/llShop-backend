package build

import (
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/models/pojo"
	"shop-backend/models/vo"
)

// CreateCartProductVO 多协程构建购物车商品展示对象
func CreateCartProductVO(product *pojo.Cart, channel chan *vo.CartProductVO) *vo.CartProductVO {
	cartVO := new(vo.CartProductVO)
	cartVO.WG.Add(2)
	// 将此对象放入到通道中
	channel <- cartVO
	// 开启两个协程，并发构建购物车商品展示对象
	go setSkuInfo(channel, product.SkuID)
	go setSpuInfo(channel, product.SkuID)
	// 阻塞在此，直到VO对象完成sku、spu信息的填充
	cartVO.WG.Wait()

	// 从通道中获取已经构建完spu、sku信息的对象
	cartVO = <-channel
	// 除sku、spu外的属性
	cartVO.Count = product.Count
	cartVO.Selected = product.Selected
	cartVO.CreatedTime = product.CreatedTime
	cartVO.ProductSkuSpecification = product.Specification

	if cartVO.Err != nil {
		zap.L().Error("生成购物车商品展示对象失败")
	}
	return cartVO
}

// 完成对购物车商品对象中sku属性的赋值
func setSkuInfo(channel chan *vo.CartProductVO, skuId int64) {
	sku, err := mysql.SelectSkuBySkuID(skuId)
	// 从通道中获取对象
	cartVO := <-channel
	if err != nil {
		cartVO.Err = err
		return
	}

	// 赋值
	cartVO.SkuID = sku.ID
	cartVO.Title = sku.Title
	cartVO.Price = sku.Price
	channel <- cartVO
	defer cartVO.WG.Done()
}

// 完成对购物车商品对象中spu属性的赋值
func setSpuInfo(channel chan *vo.CartProductVO, skuId int64) {
	spu, err := mysql.SelectSpuBySkuID(skuId)

	// 从通道中获取对象
	cartVO := <-channel
	if err != nil {
		cartVO.Err = err
		return
	}

	// 赋值
	cartVO.DefaultPicUrl = spu.DefaultPicUrl
	cartVO.PublishStatus = spu.PublishStatus
	channel <- cartVO
	defer cartVO.WG.Done()
}
