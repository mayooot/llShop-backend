package logic

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/dao/redis"
	"shop-backend/models/pojo"
	"shop-backend/models/vo"
	"strings"
	"time"
)

// AddCartProduct 添加商品到用户购物车
func AddCartProduct(userID, skuID int64, count int, specification string) error {
	// 查询该sku是否存在，是否还是上架状态
	sku, err := mysql.SelectSkuBySkuID(skuID)
	if err != nil || sku.Valid == 0 {
		zap.L().Error("用户添加到购物车的商品不存在或已下架", zap.Error(err), zap.Int64("skuID", skuID))
		return errors.New("用户添加到购物车的商品不存在或已下架")
	}

	// 先检查加入购物车的商品sku规格，是否存在于该商品spu的总规格中
	err, exist := CheckSpecificationExist(skuID, specification)
	if err != nil || !exist {
		zap.L().Error("用户添加到购物车的商品规格", zap.Error(err))
		return errors.New("用户添加到购物车的规格不存在")
	}

	// 根据用户ID和商品skuID、规格查询用户购物车中是否已经有该商品的记录
	oldCart, exist := mysql.SelectOneCartProductByUIDAndSkuId(userID, skuID, specification)
	if exist {
		// 如果该商品已经存在于该用户购物车下，更新商品购买数量
		// 如果用户本来的购买数量为10，现在传递的为-20。这样用户该商品的购买数量就为-10。这是错误的。
		count += oldCart.Count
		if count < 0 {
			// 如果用户购买数量小于0
			zap.L().Error("用户添加商品到购物车的数量小于0", zap.Int("count", count), zap.Int("stock", sku.Stock))
			return errors.New("用户添加商品到购物车的数量小于0")
		}
		if count > sku.Stock {
			// 如果用户购买数量大于库存
			zap.L().Error("用户添加商品到购物车的数量大于该商品库存", zap.Int("count", count), zap.Int("stock", sku.Stock))
			return errors.New("用户添加商品到购物车的数量大于该商品库存")
		}

		// 更新购买数量
		err = mysql.UpdateCartProductByUIDAndSkuId(userID, skuID, count)
		if err != nil {
			// 更新失败
			return err
		}
	} else {
		// 如果该商品不存在于用户购物车下
		if count <= 0 {
			// 并且用户传递的数量小于等于0
			zap.L().Error("用户添加商品到购物车的数量小于等于0", zap.Int("count", count), zap.Int("stock", sku.Stock))
			return errors.New("用户添加商品到购物车的数量小于等于0")
		}
		if count > sku.Stock {
			// 并且用户购买数量大于库存
			zap.L().Error("用户添加商品到购物车的数量大于该商品库存", zap.Int("count", count), zap.Int("stock", sku.Stock))
			return errors.New("用户添加商品到购物车的数量大于该商品库存")
		}

		// 新增一条记录
		err = mysql.InsertCartProduct(userID, skuID, count, specification)
		if err != nil {
			// 插入失败
			return err
		}
	}

	// 添加到Redis缓存中
	// 构建pojo.Cart对象
	product := &pojo.Cart{
		UserID:        userID,
		SkuID:         skuID,
		Specification: specification,
		Count:         count,
		Selected:      1,
		CreatedTime:   time.Now(),
	}
	// 创建一个存入购物车商品展示对象的通道，缓存区为1
	channel := make(chan *vo.CartProductVO, 1)
	defer close(channel)
	// 添加到Redis缓存中
	return redis.AddCartProduct(userID, skuID, createCartProductVO(product, channel))
}

// DelCartProduct 删除用户购物车中的某个商品
func DelCartProduct(userID, skuID int64, specification string) error {
	return mysql.DelCartProductBySkuIDAndUID(userID, skuID, specification)
}

// UpdateCartProductSelected 修改购物车中商品的勾选状态
func UpdateCartProductSelected(userID, skuID int64, selected int, specification string) error {
	return mysql.UpdateCartProductSelected(userID, skuID, selected, specification)
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
		zap.L().Error("获取用户购物车信息集合失败", zap.Error(err))
		return nil, err
	}

	// 封装CartProductVO集合
	data := make([]*vo.CartProductVO, 0)

	// 创建一个存入购物车商品展示对象的通道，缓存区为1
	channel := make(chan *vo.CartProductVO, 1)
	defer close(channel)

	for _, product := range cartList {
		// 遍历用户购物车信息集合
		cartVO := createCartProductVO(product, channel)
		data = append(data, cartVO)
	}
	return data, nil
}

// CheckSpecificationExist 检查用户输入的商品规则是否存在
func CheckSpecificationExist(skuID int64, specification string) (error, bool) {
	// 获取spu
	spu, err := mysql.SelectSpuBySkuID(skuID)
	if err != nil {
		zap.L().Error("用户添加到购物车的商品对应的spu不存在", zap.Error(err), zap.Int64("skuID", skuID))
		return errors.New("用户添加到购物车的商品对应的spu不存在"), false
	}

	// 解析商品规格
	specMap := make(map[string][]string)
	err = json.Unmarshal([]byte(spu.ProductSpecification), &specMap)
	if err != nil {
		zap.L().Error("解析商品规格失败", zap.Error(err))
		return errors.New("解析商品规格失败"), false
	}

	// // 字符串解析匹配
	specList := specMap["规格"]
	// // 去除前端传递的商品规格中的空格
	specification = strings.TrimSpace(specification)
	var exist bool
	for _, spec := range specList {
		if spec == specification {
			// 存在该规格
			exist = true
			break
		}
	}
	return nil, exist
}

// 多协程构建购物车商品展示对象
func createCartProductVO(product *pojo.Cart, channel chan *vo.CartProductVO) *vo.CartProductVO {
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
