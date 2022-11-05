package logic

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/dao/redis"
	"shop-backend/models/vo"
	"shop-backend/rabbitmq"
	"shop-backend/utils/build"
	"strings"
)

// AddCartProduct 添加商品到用户购物车
// 缓存设计：无论Redis中是否有该商品的缓存，都应该被覆盖。所以每次新增商品时，无需判断缓存是否存在。商品入库后，回写到缓存即可
func AddCartProduct(userID, skuID int64, count int, specification string) error {
	// 查询该商品sku是否存在，是否还是上架状态
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
		totalCount := oldCart.Count + count

		if totalCount < 0 {
			// 如果用户购买数量小于0
			zap.L().Error("用户添加商品到购物车的数量小于0", zap.Int("totalCount", totalCount), zap.Int("stock", sku.Stock))
			return errors.New("用户添加商品到购物车的数量小于0")
		}
		if totalCount > sku.Stock {
			// 如果用户购买数量大于库存
			zap.L().Error("用户添加商品到购物车的数量大于该商品库存", zap.Int("totalCount", totalCount), zap.Int("stock", sku.Stock))
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
	return nil
}

// DelCartProduct 删除用户购物车中的某个商品
// 缓存设计：从数据库中删除该商品后，也应该删除缓存中的商品数据
func DelCartProduct(userID, skuID int64, specification string) error {
	err := mysql.DelCartProductBySkuIDAndUID(userID, skuID, specification)
	if err != nil {
		zap.L().Error("删除用户购物车中的某个商品失败", zap.Error(err))
		return err
	}
	return nil
}

// UpdateCartProductSelected 修改购物车中商品的勾选状态
// 缓存设计：在数据库中修改完单个商品的勾选状态后，需要覆盖缓存中该商品的缓存
func UpdateCartProductSelected(userID, skuID int64, selected int, specification string) error {
	// cart, err := mysql.UpdateCartProductSelected(userID, skuID, selected, specification)
	_, err := mysql.UpdateCartProductSelected(userID, skuID, selected, specification)
	if err != nil {
		zap.L().Error("修改购物车中商品的勾选状态失败", zap.Error(err), zap.Int64("skuID", skuID))
		return err
	}
	return nil
}

// GetCarProductListCount  返回用户购物车中的商品数量
func GetCarProductListCount(userID int64) (int, error) {
	// 获取用户购物车信息集合
	list, err := GetCarProductList(userID)
	if err != nil {
		zap.L().Error("获取用户购物车中的商品数量失败", zap.Error(err))
		return 0, err
	}
	return len(list), nil
}

// GetCarProductList 返回用户购物车中的商品集合
// 缓存设计：先查看缓存中是否有用户购物车列表；如果有，直接返回。如果没有，从数据库中查询出来后，再回写到缓存中
func GetCarProductList(userID int64) ([]*vo.CartProductVO, error) {
	// 查看缓存中是否有该用户购物车列表数据
	list, err := redis.GetCartProductList(userID)
	if list != nil && err == nil {
		zap.L().Info("使用缓存获取用户购物车中的商品集合成功")
		return list, nil
	}
	zap.L().Error("使用缓存获取用户购物车中的商品集合失败", zap.Error(err))

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
		cartVO := build.CreateCartProductVO(product, channel)
		data = append(data, cartVO)
	}

	// 将要加入Redis缓存的用户购物车列表异步发送到MQ
	rabbitmq.SendListMess2Queue(&vo.UserCartProductVOList{
		UserID:   userID,
		CartList: data,
	})

	return data, nil
}

// CheckSpecificationExist 检查用户输入的商品规则是否存在
func CheckSpecificationExist(skuID int64, specification string) (error, bool) {
	var correctSpec string
	ret, ok := redis.GetSpuSpecification(skuID)
	if ok {
		// Redis缓存中该spu规格信息存在
		correctSpec = ret
		zap.L().Info("成功使用Redis缓存中的商品规格", zap.Int64("skuID", skuID), zap.String("specification", correctSpec))
	} else {
		// Redis缓存中不存在，从数据库中获取spu规格信息
		spu, err := mysql.SelectSpuBySkuID(skuID)
		if err != nil {
			zap.L().Error("用户添加到购物车的商品对应的spu不存在", zap.Error(err), zap.Int64("skuID", skuID))
			return errors.New("用户添加到购物车的商品对应的spu不存在"), false
		}
		correctSpec = spu.ProductSpecification
		// 回写到Redis中
		if err := redis.SetSpuSpecification(skuID, correctSpec); err != nil {
			zap.L().Error("回写spu商品规格到Redis中失败", zap.Error(err))
		}
	}

	// 解析商品规格
	specMap := make(map[string][]string)
	err := json.Unmarshal([]byte(correctSpec), &specMap)
	if err != nil {
		zap.L().Error("解析商品规格失败", zap.Error(err))
		return errors.New("解析商品规格失败"), false
	}

	//  字符串解析匹配
	specList := specMap["规格"]
	// 去除前端传递的商品规格中的空格
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
