package redis

import (
	"encoding/json"
	"go.uber.org/zap"
	"shop-backend/models/vo"
	"time"
)

var (
	usersReceiverAddressPrefix     = "user:receiver:address:"
	usersReceiverAddressLivingTime = time.Hour * 24 * 30
)

// GetAllAddress 从Redis中获取所有地址
func GetAllAddress() ([]*vo.PCDDicVO, error) {
	str, err := rdb.Get(usersReceiverAddressPrefix).Result()
	if err != nil {
		zap.L().Error("从Redis中获取所有地址缓存失败", zap.Error(err))
		return nil, err
	}

	data := make([]*vo.PCDDicVO, 0)
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		zap.L().Error("反序列化地址信息失败", zap.Error(err))
		return nil, err
	}

	return data, nil
}

// SetAllAddress 将所有的地址添加到Redis缓存中
func SetAllAddress(data []*vo.PCDDicVO) error {
	dataJson, _ := json.Marshal(data)
	if err := rdb.Set(usersReceiverAddressPrefix, dataJson, usersReceiverAddressLivingTime).Err(); err != nil {
		zap.L().Error("将所有的地址添加到Redis缓存失败", zap.Error(err))
		return err
	}
	return nil
}
