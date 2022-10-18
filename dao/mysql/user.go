package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"shop-backend/models"
	"strconv"
)

// 数据加密秘钥
const secret = "shop-backend"

// InsertUser 插入一条用户数据
func InsertUser(u *models.User) (err error) {
	// 密码加密
	u.Password = encryptPass(u.Password)
	sqlStr := `insert into ums_user(user_id,username, phone, password) values (?, ?, ?, ?)`
	// 入库
	_, err = db.Exec(sqlStr, u.UserID, u.Username, u.Phone, u.Password)
	if err != nil {
		zap.L().Error("InsertUser failed", zap.Error(err))
	}
	return
}

// QueryOneUserByPhone 通过手机号查询用户是否已经注册。用户存在返回true，否则返回false
func QueryOneUserByPhone(phone string) (uid int64, exist bool) {
	strStr := `select user_id from ums_user where phone = ?`
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
	sqlStr := `select user_id, phone, password from ums_user where phone = ?`
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
	sqlStr := `select username, avatar from ums_user where user_id = ?`
	err = db.Get(info, sqlStr, uid)
	return
}

// QueryInfosByUID 查询用户详细信
func QueryInfosByUID(uid int64) (infos *models.UserInfos, err error) {
	infos = new(models.UserInfos)
	sqlStr := `select 
					user_id, username, phone, email, avatar, gender, create_time
				from ums_user
				where user_id = ?`
	err = db.Get(infos, sqlStr, uid)
	return
}

// UpdateUserInfosByUID 修改用户个人信息
func UpdateUserInfosByUID(infos *models.ParamInfos) (err error) {
	sqlStr := `update ums_user
					set username = ?, phone = ?, email = ?, avatar = ?, gender = ?
				where user_id = ? `
	// uid字符串转成int64，因为controller层已经判断前端传递的uid是否和JWT中间件存储的uid是否相同
	// 所有这里可以忽略类型转换的err
	uid, _ := strconv.ParseInt(infos.Id, 10, 64)
	gender, _ := strconv.ParseInt(infos.Gender, 10, 8)
	_, err = db.Exec(sqlStr, infos.Username, infos.Phone, infos.Email, infos.Avatar, gender, uid)
	return
}

// UpdateAvatarByUID 使用用户ID修改用户头像
func UpdateAvatarByUID(id int64, path string) error {
	sqlStr := `update ums_user set avatar = ? where user_id = ?`
	res, err := db.Exec(sqlStr, path, id)
	affectRow, _ := res.RowsAffected()
	if affectRow == 0 {
		zap.L().Error("dao user update avatar affect row is 0", zap.Error(err), zap.Int64("uid", id))
		return errors.New("修改失败")
	}
	return err
}

// encryptPass 使用秘钥采用md5算法加密用户密码
func encryptPass(oPass string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPass)))
}
