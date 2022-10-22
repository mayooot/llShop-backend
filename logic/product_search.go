package logic

import (
	"shop-backend/dao/mysql"
	"shop-backend/models/dto"
	"shop-backend/models/vo"
)

func Search(condition *dto.SearchCondition) ([]*vo.Product, error) {
	data, err := mysql.SelectProductSearchCondition(condition)
	if err != nil {
		return nil, err
	}
	return data, nil
}
