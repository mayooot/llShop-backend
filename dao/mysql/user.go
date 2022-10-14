package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"go.uber.org/zap"
	"shop-backend/models"
)

// 数据加密秘钥
const secret = "shop-backend"

// InsertUser 插入一条用户数据
func InsertUser(u *models.User) (err error) {
	// 密码加密
	u.Password = encryptPass(u.Password)
	sqlStr := `insert into user(user_id,username, phone, password) values (?, ?, ?, ?)`
	// 入库
	_, err = db.Exec(sqlStr, u.UserID, u.Username, u.Phone, u.Password)
	if err != nil {
		zap.L().Error("InsertUser failed", zap.Error(err))
	}
	return
}

// QueryOneUserByPhone 通过手机号查询用户是否已经注册。用户存在返回true，否则返回false
func QueryOneUserByPhone(phone string) (uid int64, exist bool) {
	strStr := `select user_id from user where phone = ?`
	err := db.Get(&uid, strStr, phone)
	if err != nil {
		exist = false
		return
	}
	exist = true
	return
}

// QueryOneUserByPhoneAndPass 通过手机号和密码校验用户是否存在
func QueryOneUserByPhoneAndPass(u *models.User) bool {
	// 用户未加密密码
	originPass := u.Password
	sqlStr := `select user_id, phone, password from user where phone = ?`
	err := db.Get(u, sqlStr, u.Phone)
	if err != nil {
		return false
	}
	// 判断用户密码是否正确
	pass := encryptPass(originPass)
	if pass != u.Password {
		// 密码错误
		return false
	}
	// 登录成功
	return true
}

// QuerySomeInfoByUID 获取用户购物车数量、用户头像、用户名称
func QuerySomeInfoByUID(uid int64) (info *models.SomeInfo, err error) {
	info = new(models.SomeInfo)
	sqlStr := `select username, avatar from user where user_id = ?`
	err = db.Get(info, sqlStr, uid)
	return
}

// encryptPass 使用秘钥采用md5算法加密用户密码
func encryptPass(oPass string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPass)))
}
