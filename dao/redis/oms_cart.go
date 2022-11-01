package redis

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"shop-backend/models/vo"
	"shop-backend/utils/concatstr"
	"strconv"
	"time"
)

var (
	cartPrefix     = "order:cart:"
	cartLivingTime = time.Hour * 24 * 7
)

// AddCartProduct 添加购物车商品展示对象到Redis缓存中
func AddCartProduct(userID, skuID int64, product *vo.CartProductVO) error {
	// 将购物车商品展示对象转换为json格式
	data, err := json.Marshal(product)
	if err != nil {
		zap.L().Error("购物车商品序列化为json失败", zap.Error(err))
		return err
	}

	// 开启一个带有事务的管道(同时执行多条命令)
	key := concatstr.ConcatString(cartPrefix, strconv.FormatInt(userID, 10))
	pipe := rdb.TxPipeline()
	// 添加一条购物车商品数据到Redis缓存
	_, err = rdb.HSet(key,
		strconv.FormatInt(skuID, 10),
		string(data)).Result()
	if err != nil {
		zap.L().Error("添加一条购物车商品数据到Redis缓存失败", zap.Error(err))
		return err
	}
	// 设置失效时间
	_, err = rdb.Expire(key, cartLivingTime).Result()
	if err != nil {
		zap.L().Error("设置购物车商品TTL失败", zap.Error(err))
		return err
	}
	// 执行
	_, err = pipe.Exec()
	if err != nil {
		zap.L().Error("添加购物车商品展示对象到Redis缓存中失败", zap.Error(err))
		return err
	}
	zap.L().Info("添加购物车商品展示对象到Redis缓存中成功")
	return nil
}

// DelCartProduct 删除用户购物车缓存中的单个商品
func DelCartProduct(userID, skuID int64) error {
	key := concatstr.ConcatString(cartPrefix, strconv.FormatInt(userID, 10))
	err := rdb.HDel(key, strconv.FormatInt(skuID, 10)).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 如果异常Nil。说明Redis中不存在该key。可能是已经过期
			err = nil
		} else {
			zap.L().Error("删除用户购物车缓存中的某个商品失败", zap.Error(err))
		}
		return err
	}
	zap.L().Info("删除用户购物车缓存中的单个商品成功", zap.Int64("skuID", skuID))
	return nil
}

// AddCartProductList  添加用户购物车列表到Redis缓存中
func AddCartProductList(userID int64, cartList []*vo.CartProductVO) error {
	for _, cart := range cartList {
		if err := AddCartProduct(userID, cart.SkuID, cart); err != nil {
			zap.L().Error("添加用户购物车列表到Redis缓存中失败", zap.Int64("skuID", cart.SkuID))
			return err
		}
	}
	return nil
}

// GetCartProductList  从Redis缓存中获取用户购物车列表
func GetCartProductList(userID int64) ([]*vo.CartProductVO, error) {
	key := concatstr.ConcatString(cartPrefix, strconv.FormatInt(userID, 10))
	result, err := rdb.HGetAll(key).Result()
	if err != nil || result == nil {
		zap.L().Error("从Redis缓存中获取用户购物车列表失败", zap.Error(err))
		return nil, err
	}

	data := make([]*vo.CartProductVO, 0)
	for _, ret := range result {
		cart := new(vo.CartProductVO)
		_ = json.Unmarshal([]byte(ret), cart)
		data = append(data, cart)
	}
	// if err := json.Unmarshal([]byte(str), &data); err != nil {
	// 	zap.L().Error("反序列缓存中的用户购物车列表失败", zap.Error(err))
	// 	return nil, err
	// }

	return data, nil
}
