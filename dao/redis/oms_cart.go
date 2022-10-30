package redis

import (
	"encoding/json"
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
	return nil
}
