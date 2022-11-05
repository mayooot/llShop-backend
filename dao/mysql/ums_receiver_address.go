package mysql

import (
	"go.uber.org/zap"
	"shop-backend/models/pojo"
)

// SelectAllAddress 获取数据库中的所有地址
func SelectAllAddress() ([]*pojo.PCDDic, error) {
	addresses := make([]*pojo.PCDDic, 0)
	if err := db.Model(&pojo.PCDDic{}).Find(&addresses).Error; err != nil {
		zap.L().Error("获取数据库中的所有地址失败", zap.Error(err))
		return nil, err
	}
	return addresses, nil
}

// InsertReceiverAddress 新增一条用户收货地址信息
func InsertReceiverAddress(address *pojo.ReceiverAddress) error {
	tx := db.Begin()

	// 如果用户的收货地址为空，也就是0个收货地址。那么新增的第一条就是默认的收货地址
	var count int64
	if err := tx.Model(&pojo.ReceiverAddress{}).Where("user_id = ?", address.UserID).Count(&count).Error; err != nil {
		tx.Rollback()
		zap.L().Error("获取用户收货地址数量失败", zap.Error(err))
		return err
	}
	if count == 0 {
		// 第一条收货地址就是默认地址
		address.DefaultStatus = 1
	}

	// 新增数据
	if err := tx.Create(address).Error; err != nil {
		tx.Rollback()
		zap.L().Error("新增一条用户收货地址信息到数据库失败", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

// UpdateDefaultReceiverAddress 修改用户默认收货地址状态和信息，并将之前的默认地址状态修改为0
func UpdateDefaultReceiverAddress(address *pojo.ReceiverAddress) error {
	tx := db.Begin()
	// 将用户的所有收货地址状态都设置为2
	err := tx.Model(&pojo.ReceiverAddress{}).Where("user_id = ?", address.UserID).Update("default_status", 2).Error
	if err != nil {
		tx.Rollback()
		zap.L().Error("将用户的所有收货地址状态都设置为2失败", zap.Error(err))
		return err
	}
	// 更新本行信息
	err = tx.Debug().Model(&pojo.ReceiverAddress{}).
		Where("id = ? and user_id = ?", address.ID, address.UserID).
		Updates(&pojo.ReceiverAddress{
			CountyID:      address.CountyID,
			UserName:      address.UserName,
			PhoneNumber:   address.PhoneNumber,
			DefaultStatus: address.DefaultStatus,
			DetailAddress: address.DetailAddress,
		}).Error
	if err != nil {
		tx.Rollback()
		zap.L().Error("修改用户收货地址信息失败", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

// UpdateReceiverAddress 使用主键ID和用户ID修改用户收货地址。如果只使用主键ID，那么用户如果传递错误的主键ID，那么就会修改到其他用户的收货地址
func UpdateReceiverAddress(address *pojo.ReceiverAddress) error {
	err := db.Debug().Model(&pojo.ReceiverAddress{}).
		Where("id = ? and user_id = ?", address.ID, address.UserID).
		Updates(&pojo.ReceiverAddress{
			CountyID:      address.CountyID,
			UserName:      address.UserName,
			PhoneNumber:   address.PhoneNumber,
			DefaultStatus: address.DefaultStatus,
			DetailAddress: address.DetailAddress,
		}).Error
	if err != nil {
		zap.L().Error("修改用户收货地址信息失败", zap.Error(err))
		return err
	}
	return nil
}

// SelectPersonAllAddress 查询出用户所有的收货地址
func SelectPersonAllAddress(uid int64) ([]*pojo.ReceiverAddress, error) {
	data := make([]*pojo.ReceiverAddress, 0)
	if err := db.Model(&pojo.ReceiverAddress{}).Where("user_id = ?", uid).Find(&data).Error; err != nil {
		zap.L().Error("查询用户所有的收货地址失败", zap.Error(err))
		return nil, err
	}
	return data, nil
}

// SelectPCDByID 使用主键ID获取PCD信息
func SelectPCDByID(id int) (*pojo.PCDDic, error) {
	pcdPojo := new(pojo.PCDDic)
	if err := db.Model(&pojo.PCDDic{}).Where("id = ?", id).First(pcdPojo).Error; err != nil {
		zap.L().Error("使用主键ID获取PCD信息失败", zap.Error(err))
		return nil, err
	}
	return pcdPojo, nil
}

// DelReceiverAddress 使用主键ID和用户ID删除用户的收货地址
func DelReceiverAddress(id int, uid int64) error {
	if err := db.Where("id = ? and user_id = ?", id, uid).Delete(&pojo.ReceiverAddress{}).Error; err != nil {
		zap.L().Error("使用主键ID和用户ID删除用户的收货地址失败", zap.Error(err))
		return err
	}
	return nil
}
