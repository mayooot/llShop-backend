package redis

import (
	"go.uber.org/zap"
	"shop-backend/utils/concatstr"
	"strconv"
	"time"
)

var (
	orderOrderNumPrefix     = "order:num:"
	orderOrderNumLivingTime = time.Minute * 5
)

// SetOrderNumber 将订单编号设置进Redis
func SetOrderNumber(orderNum int64) error {
	key := concatstr.ConcatString(orderOrderNumPrefix, strconv.FormatInt(orderNum, 10))
	if err := rdb.Set(key, nil, orderOrderNumLivingTime).Err(); err != nil {
		zap.L().Error("将订单编号设置进Redis失败", zap.Error(err))
		return err
	}
	return nil
}

// GetOrderNumber 判断订单号是否存在于Redis中
func GetOrderNumber(orderNum int64) bool {
	key := concatstr.ConcatString(orderOrderNumPrefix, strconv.FormatInt(orderNum, 10))
	if err := rdb.Get(key).Err(); err != nil {
		return false
	}
	return true
}
