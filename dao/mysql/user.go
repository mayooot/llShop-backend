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
	sqlStr := `insert into user(user_id, phone, password) values (?, ?, ?)`
	// 入库
	_, err = db.Exec(sqlStr, u.UserID, u.Phone, u.Password)
	if err != nil {
		zap.L().Error("InsertUser failed", zap.Error(err))
	}
	return
}

// QueryOneUser 通过手机号查询用户是否已经注册
func QueryOneUser(phone string) bool {
	strStr := `select user_id from user where phone = ?`
	_, err := db.Exec(strStr, phone)
	if err != nil {
		// 如果查询不出来
		return false
	}
	return true
}

// encryptPass 使用秘钥采用md5算法加密用户密码
func encryptPass(oPass string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPass)))
}
