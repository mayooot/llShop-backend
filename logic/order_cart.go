package logic

import (
	"errors"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/models/vo"
)

// AddCartProduct 添加商品到用户购物车
func AddCartProduct(userID, skuID int64, count int) error {
	sku, err := mysql.SelectSkuBySkuID(skuID)
	if err != nil || sku.Valid == 0 {
		zap.L().Error("用户添加到购物车的商品不存在", zap.Error(err), zap.Int64("skuID", skuID))
		return errors.New("用户添加到购物车的商品不存在")
	}

	// 根据用户ID和商品skuID查询用户购物车中是否已经有该商品的记录
	oldCart, exist := mysql.SelectOneCartProductByUIDAndSkuId(userID, skuID)
	if exist {
		// 如果该商品已经存在于该用户购物车下，更新数量
		count += oldCart.Count
		if count > sku.Stock {
			// 如果用户购买数量大于库存
			zap.L().Error("用户添加商品到购物车的数量大于该商品库存", zap.Int("count", count), zap.Int("stock", sku.Stock))
			return errors.New("用户添加商品到购物车的数量大于该商品库存")
		}
		err = mysql.UpdateCartProductByUIDAndSkuId(userID, skuID, count)
		if err != nil {
			// 更新失败
			return err
		}
		// 更新成功
		return nil
	} else {
		if count > sku.Stock {
			// 如果用户购买数量大于库存
			zap.L().Error("用户添加商品到购物车的数量大于该商品库存", zap.Int("count", count), zap.Int("stock", sku.Stock))
			return errors.New("用户添加商品到购物车的数量大于该商品库存")
		}
		return mysql.InsertCartProduct(userID, skuID, count)
	}

}

// DelCartProduct 删除用户购物车中的某个商品
func DelCartProduct(userID, skuID int64) error {
	return mysql.DelCartProductBySkuIDAndUID(userID, skuID)
}

// GetCarProductListCount  返回用户购物车中的商品数量
func GetCarProductListCount(userID int64) (int, error) {
	// 获取用户购物车信息集合
	cartList, err := mysql.SelectCartList(userID)
	if err != nil {
		return 0, err
	}
	return len(cartList), nil
}

// GetCarProductList 返回用户购物车中的商品集合
func GetCarProductList(userID int64) ([]*vo.CartProductVO, error) {
	// 获取用户购物车信息集合
	cartList, err := mysql.SelectCartList(userID)
	if err != nil {
		return nil, err
	}

	// 封装CartProductVO集合
	data := make([]*vo.CartProductVO, 0)
	// 创建一个存入购物车商品展示对象的通道，缓存区为1
	channel := make(chan *vo.CartProductVO, 1)
	defer close(channel)

	for _, product := range cartList {
		// 遍历用户购物车信息集合
		cartVO := new(vo.CartProductVO)
		cartVO.WG.Add(2)
		// 将此对象放入到通道中
		channel <- cartVO
		// 开启两个协程，并发构建购物车商品展示对象
		go setSkuInfo(channel, product.SkuID)
		go setSpuInfo(channel, product.SkuID)
		cartVO.WG.Wait()

		// 从通道中获取已经构建完spu、sku信息的对象
		cartVO = <-channel
		cartVO.CreatedTime = product.CreatedTime
		cartVO.Count = product.Count
		if cartVO.Err != nil {
			zap.L().Error("获取用户购物车商品列表错误", zap.Int64("uid", userID), zap.Int64("skuID", product.SkuID))
			continue
		}
		data = append(data, cartVO)
	}

	return data, nil
}

func setSkuInfo(channel chan *vo.CartProductVO, skuId int64) {
	sku, err := mysql.SelectSkuBySkuID(skuId)
	// 从通道中获取对象
	cartVO := <-channel
	if err != nil {
		cartVO.Err = err
		return
	}

	// 赋值
	cartVO.Price = sku.Price
	cartVO.ProductSkuSpecification = sku.ProductSkuSpecification
	cartVO.Title = sku.Title
	channel <- cartVO
	defer cartVO.WG.Done()
}

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
