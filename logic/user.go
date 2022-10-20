package logic

import (
	"errors"
	"github.com/DanPlayer/randomname"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/dao/redis"
	"shop-backend/models/dto"
	"shop-backend/models/pojo"
	"shop-backend/models/vo"
	"shop-backend/utils/check"
	"shop-backend/utils/gen"
)

var (
	ErrorUserIsRegistered    = errors.New("用户已注册，请直接登录")
	ErrorMustRequestCode     = errors.New("请先获取验证码")
	ErrorRequestCodeFrequent = errors.New("验证码已发送，请注意查收")
	ErrorWrongVerifyCode     = errors.New("验证码错误")
	ErrorPassWeak            = errors.New("密码强度太弱啦~")
	ErrorUserNotExist        = errors.New("用户不存在~")
	ErrorWrongPass           = errors.New("密码错误")
	ErrorServeBusy           = errors.New("服务器繁忙")
)

// SendVerifyCode 生成验证码并缓存到Redis中
func SendVerifyCode(phone string) (code string, err error) {
	// 接口幂等性
	code, err = redis.GetVerifyCode(phone)
	if code != "" {
		err = ErrorRequestCodeFrequent
		zap.L().Error("users frequently get verification codes", zap.String("phone", phone))
		return
	}
	if code, err = gen.GenVerifyCode(); err != nil {
		// 生成验证码失败
		zap.L().Error("utils.GenVerifyCode failed", zap.Error(err))
		return
	}
	if err = redis.SetVerifyCode(phone, code); err != nil {
		// 缓存到Redis失败
		zap.L().Error("redis.SetVerifyCode failed", zap.Error(err))
		return
	}
	return
}

// SignUp 用户注册逻辑
func SignUp(u *dto.ParamSignUp) error {
	// 到此，用户手机号格式一定是正确的
	// 如果用户已经注册
	if ok := mysql.SelectUserByPhone(u.Phone); ok {
		return ErrorUserIsRegistered
	}
	// 通过手机号从Redis获取验证码
	niceCode, err := redis.GetVerifyCode(u.Phone)
	if err != nil {
		// 如果获取验证码失败
		zap.L().Error("redis.GetVerifyCode(u.Phone) failed", zap.Error(err))
		return err
	}
	if niceCode == "" {
		return ErrorMustRequestCode
	}

	if u.Code != niceCode {
		// 用户输入的验证码不正确
		return ErrorWrongVerifyCode
	}

	// 校验密码强度
	if err = check.CheckPass(u.Password); err != nil {
		return ErrorPassWeak
	}

	// 生成uid
	uid := gen.GenSnowflakeId()
	user := &pojo.UmsUser{
		ID:       uid,
		Username: randomname.GenerateName(),
		Password: u.Password,
		Phone:    u.Phone,
		Avatar:   "https://richarli.oss-cn-beijing.aliyuncs.com/images/20221018175133.png",
	}

	// 入库
	err = mysql.InsertUser(user)
	return err
}

// Login 登录逻辑
func Login(p *dto.ParamLogin) (uid int64, aToken, rToken string, err error) {
	var ok bool
	ok = mysql.SelectUserByPhone(p.Phone)
	if !ok {
		// 如果用户不存在
		err = ErrorUserNotExist
		return
	}

	// 构建User实例
	user := &pojo.UmsUser{
		Phone:    p.Phone,
		Password: p.Password,
	}

	if uid, ok = mysql.SelectUserByPhoneAndPass(user); !ok {
		// 用户输入密码错误
		err = ErrorWrongPass
		return
	}

	// 用户校验通过，生成AccessToken和RefreshToken
	aToken, rToken, err = gen.GenToken(uid)

	// 将AccessToken缓存到Redis中，用来完成同一时间只有一台设备可以登录
	err = redis.SetAccessToken(uid, aToken)
	if err != nil {
		// 插入AccessToken失败
		err = ErrorServeBusy
		return
	}
	return
}

// GetSomeInfo 返回简略信息
func GetSomeInfo(uid int64) (info *vo.SomeInfo, err error) {
	info, err = mysql.SelectSomeInfoByUID(uid)
	return
}

// GetUserInfos 获取用户详细信息
// todo 邮箱登录功能未实现
func GetUserInfos(uid int64) (infos *vo.UserInfos, err error) {
	return mysql.SelectInofsByUID(uid)
}

// UpdateInfos 更新用户信息
func UpdateInfos(infos *dto.ParamInfos) (err error) {
	return mysql.UpdateUserInfosByUID(infos)
}

// SignOut 用户退出操作，删除Redis中保存的AccessToken
func SignOut(idStr string) (err error) {
	err = redis.DelAccessTokenByUID(idStr)
	return
}

// UpdateUserAvatar 修改用户头像
func UpdateUserAvatar(uid int64, path string) error {
	return mysql.UpdateAvatarByUID(uid, path)
}
