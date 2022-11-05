package logic

import (
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/dao/redis"
	"shop-backend/models/dto"
	"shop-backend/models/pojo"
	"shop-backend/models/vo"
	"strconv"
)

// GetAllAddress 获取数据库中所有的地址信息
func GetAllAddress() ([]*vo.PCDDicVO, error) {
	// 先从缓存中获取地址信息
	list, err := redis.GetAllAddress()
	if err == nil && len(list) > 0 {
		// 缓存中的数据可用
		zap.L().Info("成功使用缓存中的地址信息")
		return list, nil
	}

	zap.L().Info("未成功使用缓存中的地址信息")
	// 所有地址信息
	addresses, err := mysql.SelectAllAddress()
	if err != nil {
		return nil, err
	}
	// 最后返回的VO对象集合
	data := make([]*vo.PCDDicVO, 0)
	for _, address := range addresses {
		if address.ParentID == 0 {
			// 省
			pcdVO := &vo.PCDDicVO{ID: address.ID, Name: address.Name}
			for _, reAddress := range addresses {
				if reAddress.ParentID == address.ID {
					// 市
					pcdVO2 := &vo.PCDDicVO{ID: reAddress.ID, Name: reAddress.Name}
					pcdVO.PicList = append(pcdVO.PicList, pcdVO2)
					for _, reeAddress := range addresses {
						if reeAddress.ParentID == reAddress.ID {
							// 区
							pcdVO3 := &vo.PCDDicVO{ID: reeAddress.ID, Name: reeAddress.Name}
							pcdVO2.PicList = append(pcdVO2.PicList, pcdVO3)
						}
					}
				}
			}
			data = append(data, pcdVO)
		}
	}

	// 存入缓存中
	_ = redis.SetAllAddress(data)
	return data, nil
}

// AddReceiverAddress 添加用户收货地址
func AddReceiverAddress(address *dto.ReceiverAddress, uid int64) error {
	return mysql.InsertReceiverAddress(createReceiverAddressPojo(address, uid))
}

// UpdateReceiverAddress 修改用户收货地址
func UpdateReceiverAddress(address *dto.ReceiverAddress, uid int64) error {
	addressPojo := createReceiverAddressPojo(address, uid)
	id, _ := strconv.ParseInt(address.ID, 10, 64)
	addressPojo.ID = id
	if address.DefaultStatus == 1 {
		// 用户要将该地址设置为默认地址
		// 此时应该将用户的其他地址的状态都设置为2(不是默认地址)
		return mysql.UpdateDefaultReceiverAddress(addressPojo)
	} else {
		// 修改用户收货地址信息
		return mysql.UpdateReceiverAddress(addressPojo)
	}
}

// GetPersonAllAddress 获取用户所有的收货地址
func GetPersonAllAddress(uid int64) ([]*vo.AddressVO, error) {
	// 此时的地址结构体只包含区县ID，需要找出所在的省和市
	addresses, err := mysql.SelectPersonAllAddress(uid)
	if err != nil {
		return nil, err
	}

	data := make([]*vo.AddressVO, 0)
	for _, address := range addresses {
		// 封装VO对象
		addressVO := &vo.AddressVO{
			ID:            address.ID,
			CountyID:      address.CountyID,
			UserName:      address.UserName,
			PhoneNumber:   address.PhoneNumber,
			DefaultStatus: address.DefaultStatus,
			DetailAddress: address.DetailAddress,
		}
		// 区
		region, err := mysql.SelectPCDByID(address.CountyID)
		if err != nil {
			return nil, err
		}
		addressVO.Region = region.Name

		// 市
		city, err := mysql.SelectPCDByID(region.ParentID)
		if err != nil {
			return nil, err
		}
		addressVO.City = city.Name

		// 省
		province, err := mysql.SelectPCDByID(city.ParentID)
		if err != nil {
			return nil, err
		}
		addressVO.Province = province.Name

		data = append(data, addressVO)
	}
	return data, nil
}

// DelReceiverAddress 删除用户的一条收货地址
func DelReceiverAddress(id int, uid int64) error {
	return mysql.DelReceiverAddress(id, uid)
}

// 创建ReceiverAddress POJO结构体
func createReceiverAddressPojo(address *dto.ReceiverAddress, uid int64) *pojo.ReceiverAddress {
	addressPojo := &pojo.ReceiverAddress{
		UserID:        uid,
		CountyID:      address.CountyID,
		UserName:      address.ReceiverName,
		PhoneNumber:   address.ReceiverPhone,
		DefaultStatus: uint8(address.DefaultStatus),
		DetailAddress: address.DetailAddress,
	}
	return addressPojo
}
