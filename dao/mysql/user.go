package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"shop-backend/models/dto"
	"shop-backend/models/pojo"
	"shop-backend/models/vo"
	"strconv"
)

// 数据加密秘钥
const secret = "shop-backend"

// InsertUser 插入一条用户数据
func InsertUser(u *pojo.UmsUser) error {
	// 密码加密
	u.Password = encryptPass(u.Password)
	// 入库
	result := db.Create(u)
	if result.Error != nil || result.RowsAffected == 0 {
		// 如果有异常或者影响的行数为0
		zap.L().Error("插入一条用户数据失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// SelectUserByPhone 通过手机号查询用户是否已经注册。用户存在返回true，否则返回false
func SelectUserByPhone(phone string) (exist bool) {
	// 通过手机号查询
	result := db.Where("phone = ?", phone).First(&pojo.UmsUser{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果是error是记录不存在异常，说明用户不存在，返回false
			return false
		} // 反之，认为用户已经存在
		zap.L().Error("用户已存在", zap.Error(result.Error), zap.String("phone", phone))
		return true
	}
	return true
}

// SelectUserByPhoneAndPass 通过手机号和密码校验用户是否存在
func SelectUserByPhoneAndPass(u *pojo.UmsUser) (int64, bool) {
	// 用户未加密密码
	originPass := u.Password
	// 通过手机号查询
	result := db.Where("phone = ?", u.Phone).First(u)
	if result.Error != nil {
		zap.L().Error("通过手机号查询用户失败或记录为空", zap.Error(result.Error))
		return 0, false
	}
	// 判断用户密码是否正确
	pass := encryptPass(originPass)
	if pass != u.Password {
		// 密码错误
		zap.L().Error("用户输入密码错误", zap.Error(result.Error), zap.String("input pass", originPass))
		return 0, false
	}
	// 登录成功
	return u.ID, true
}

// SelectSomeInfoByUID 获取用户购物车数量、用户头像、用户名称
func SelectSomeInfoByUID(uid int64) (*vo.SomeInfoVO, error) {
	var user = new(pojo.UmsUser)
	result := db.Where("user_id", uid).First(user)
	if result.Error != nil {
		zap.L().Error("获取用户购物车数量、用户头像、用户名称失败", zap.Error(result.Error))
		return nil, result.Error
	}

	// 获取用户购物车数量
	var count int
	cartList, err := SelectCartList(uid)
	if err != nil {
		count = 0
	} else {
		count = len(cartList)
	}

	// 封装SomeInfo对象
	info := &vo.SomeInfoVO{
		Avatar:   user.Avatar,
		Username: user.Username,
		CartNum:  count,
	}
	return info, result.Error
}

// SelectInfosByUID 查询用户详细信息
func SelectInfosByUID(uid int64) (*vo.UserInfosVO, error) {
	var user = new(pojo.UmsUser)
	result := db.Where("user_id", uid).First(user)
	if result.Error != nil {
		zap.L().Error("查询用户详细信息失败", zap.Error(result.Error))
		return nil, result.Error
	}
	// 封装UserInfos对象
	infos := &vo.UserInfosVO{
		Id:          user.ID,
		Username:    user.Username,
		Phone:       user.Phone,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		CreatedTime: user.CreatedTime,
	}
	return infos, nil
}

// UpdateUserInfosByUID 修改用户个人信息
func UpdateUserInfosByUID(infos *dto.Infos) error {
	id, _ := strconv.ParseInt(infos.ID, 10, 64)
	result := db.Model(&pojo.UmsUser{ID: id}).Updates(pojo.UmsUser{
		Username: infos.Username,
		Password: encryptPass(infos.Password),
		Email:    infos.Email,
		Phone:    infos.Phone,
		Gender:   infos.Gender,
	})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		// 如果有异常为不存在该记录异常或者影响的行数为0，说明用户不存在
		zap.L().Error("修改用户信息失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// UpdateAvatarByUID 使用用户ID修改用户头像
func UpdateAvatarByUID(id int64, path string) error {
	result := db.Model(&pojo.UmsUser{ID: id}).Update("avatar", path)
	if result.Error != nil {
		zap.L().Error("修改用户头像数据库记录失败", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

// encryptPass 使用秘钥采用md5算法加密用户密码
func encryptPass(oPass string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPass)))
}
