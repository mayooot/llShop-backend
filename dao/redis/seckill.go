package redis

import (
	"shop-backend/utils/concatstr"
	"strconv"
	"time"
)

var (
	SecKillUIDPrefix    = "seckill:uid:"
	SecKillUIDivingTime = time.Second * 5
)

// SetNXSecKillUID 使用Redis SETNX命名将用户ID设置进Redis
func SetNXSecKillUID(uid int64) bool {
	key := concatstr.ConcatString(SecKillUIDPrefix, strconv.FormatInt(uid, 10))
	return rdb.SetNX(key, nil, SecKillUIDivingTime).Val()
}
