package logic

import (
	"shop-backend/dao/mysql"
	"shop-backend/models/pojo"
)

// GetAllSecKillSku 获取所有正在秒杀的商品
func GetAllSecKillSku() ([]*pojo.SecKillSku, error) {
	return mysql.SelectAllSecKillSku()
}
