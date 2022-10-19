package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-backend/logic"
	"shop-backend/models/dto"
	"shop-backend/utils/check"
	"shop-backend/utils/oss"
	"strconv"
)

// SendVerifyCodeHandler 发送手机验证码，并返回
// @Summary 获取验证码
// @Description 前端传递11位手机号，后端随机生成4位验证码。存入Redis，并异步调用阿里云SMS服务发送验证码。
// @Tags 用户相关接口
// @Produce  json
// @Param phone query string true "手机号"
// @Router /user/phone [get]
func SendVerifyCodeHandler(c *gin.Context) {
	// 获取参数
	phone := c.Query("phone")
	if phone == "" {
		zap.L().Error("The mobile phone number field is empty.")
		// phone字段为空
		ResponseError(c, CodePhoneIsNotEmpty)
		return
	}
	if ok := check.VerifyMobileFormat(phone); !ok {
		// 手机号格式不正确
		zap.L().Error("The format of the mobile phone number is wrong.")
		ResponseError(c, CodePhoneFormatError)
		return
	}
	// 发送验证码
	code, err := logic.SendVerifyCode(phone)
	zap.L().Info("Send verification code.", zap.String("code", code))
	if err != nil {
		if errors.Is(err, logic.ErrorRequestCodeFrequent) {
			// 用户频繁请求验证码
			zap.L().Warn("The verification code is frequently obtained.", zap.String("phone", phone))
			ResponseError(c, CodeRequestCodeFrequently)
			return
		}
		// 生成、发送验证码失败
		zap.L().Error("Failed to send verification code.")
		ResponseError(c, CodeServeBusy)
		return
	}
	zap.L().Info("Verification code sent successfully.", zap.String("code", code))
	ResponseSuccessWithMsg(c, "发送成功，验证码五分钟内有效", gin.H{
		"code": code,
	})
}

// SignUpHandler 用户注册
// @Summary 注册用户
// @Description 前端传递JSON类型对象，后端完成校验后注册新用户。
// @Tags 用户相关接口
// @Produce json
// @Param ParamSignUp body dto.ParamSignUp true "用户注册结构体"
// @Router /user/signup [post]
func SignUpHandler(c *gin.Context) {
	// 获取参数并校验
	p := new(dto.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		ResponseError(c, CodeInvalidParams)
		return
	}

	if ok := check.VerifyMobileFormat(p.Phone); !ok {
		// 手机号格式不正确
		ResponseError(c, CodePhoneFormatError)
		return
	}

	// 业务处理
	err := logic.SignUp(p)
	if err != nil {
		if errors.Is(err, logic.ErrorUserIsRegistered) {
			// 用户已注册
			ResponseError(c, CodeUserIsRegistered)
			return
		} else if errors.Is(err, logic.ErrorWrongVerifyCode) {
			// 如果验证码错误或过期
			ResponseError(c, CodeWrongVerifyCode)
			return
		} else if errors.Is(err, logic.ErrorPassWeak) {
			// 密码强度太低
			ResponseError(c, CodePassIsWeak)
			return
		} else if errors.Is(err, logic.ErrorMustRequestCode) {
			// 用户未获取验证码
			ResponseError(c, CodeMustRequestCode)
			return
		}
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccessWithMsg(c, "注册成功，请登录", nil)
}

// LoginHandler 用户登录
// @Summary 用户登录
// @Description 前端传递JSON类型对象，后端完成校验后登录，返回AccessToken和RefreshToken、UserID。
// @Tags 用户相关接口
// @Produce  json
// @Param ParamLogin body dto.ParamLogin true "用户登录结构体"
// @Router /user/login [post]
func LoginHandler(c *gin.Context) {
	// 获取参数并校验
	p := new(dto.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		ResponseError(c, CodeInvalidParams)
		return
	}
	uid, aToken, rToken, err := logic.Login(p)
	if err != nil {
		if errors.Is(err, logic.ErrorWrongPass) {
			ResponseError(c, CodeUsernameOrPassError)
			return
		} else if errors.Is(err, logic.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		} else {
			ResponseError(c, CodeServeBusy)
			return
		}
	}
	ResponseSuccessWithMsg(c, "登录成功", gin.H{
		// 前端json能接受的整数范围为 - (2^53 -1) ~ 2^53 - 1，而我们要传递的是int64，所以要转化成字符串
		"userId":       strconv.FormatInt(uid, 10),
		"accessToken":  aToken,
		"refreshToken": rToken,
	})
}

// SomeInfoHandler 获取用户头像、用户名和购物车数量
// @Summary 获取用户头像、用户名和购物车数量。
// @Description 后端返回用户头像、用户名和购物车数量。
// @Tags 用户相关接口
// @Produce json
// @Security x-token
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /user/someinfo [get]
func SomeInfoHandler(c *gin.Context) {
	// 获取用户简略信息
	infos, err := logic.GetSomeInfo(c.GetInt64("uid"))
	if err != nil {
		zap.L().Error("logic.GetSomeInfo(c.GetInt64(\"uid\") failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, infos)
}

// UserInfosHandler 获取用户详细的个人信息
// @Summary 获取用户个人信息
// @Description 后端返回用户个人信息。
// @Tags 用户相关接口
// @Produce json
// @Security x-token
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /user/infos [get]
func UserInfosHandler(c *gin.Context) {
	infos, err := logic.GetUserInfos(c.GetInt64("uid"))
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, infos)
}

// UserInfosUpdateHandler 用户修改个人信息
// @Summary 修改个人信息
// @Description 前端传递JSON类型对象，后端完成更新用户个人信息。
// @Tags 用户相关接口
// @Produce json
// @Security x-token
// @param Authorization header string true "Bearer AToken&RToken"
// @Param param body dto.ParamInfos true "用户个人信息结构体"
// @Router /user/infos/update [put]
func UserInfosUpdateHandler(c *gin.Context) {
	infos := new(dto.ParamInfos)
	infos.ID = strconv.FormatInt(c.GetInt64("uid"), 10)
	if err := c.ShouldBindJSON(infos); err != nil {
		// 请求参数有误
		ResponseError(c, CodeInvalidParams)
		return
	}
	if _, err := strconv.ParseInt(infos.Gender, 10, 8); err != nil {
		// 校验性别字符串
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 参数格式校验
	if !check.VerifyUsernameFormat(infos.Username) {
		ResponseError(c, CodeUsernameToLongOrToShort)
		return
	}
	if !check.VerifyMobileFormat(infos.Phone) {
		ResponseError(c, CodePhoneFormatError)
		return
	}
	if !check.VerifyEmailFormat(infos.Email) {
		ResponseError(c, CodeEmailFormatError)
		return
	}

	err := logic.UpdateInfos(infos)
	if err != nil {
		ResponseError(c, CodeUpdateInfosFailed)
	}
	ResponseSuccessWithMsg(c, "更新成功", infos)
}

// SignOutHandler 用户退出
// @Summary 用户退出
// @Description 后端清空Redis中AccessToken
// @Tags 用户相关接口
// @Produce  json
// @Security x-token
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /user/exit [delete]
func SignOutHandler(c *gin.Context) {
	err := logic.SignOut(c.GetString("uid"))
	if err != nil {
		ResponseError(c, CodeSignOutFailed)
		return
	}
	ResponseSuccessWithMsg(c, "退出成功", nil)
}

// UserInfoUpdateAvatarHandler 更新用户头像
// @Summary 用户更新头像
// @Description 前端上传图片，后端将头像上传到阿里云OSS后，修改用户头像。
// @Tags 用户相关接口
// @Produce  json
// @Security x-token
// @param Authorization header string true "Bearer AToken&RToken"
// @Router /user/infos/update/avatar [post]
func UserInfoUpdateAvatarHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		// 读取文件失败
		ResponseError(c, CodeUploadAvatarFailed)
		return
	}
	// 检查图片格式
	if err = check.CheckPic(fileHeader); err != nil {
		ResponseError(c, CodeUploadAvatarFailed)
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		// 打开文件失败
		ResponseError(c, CodeUploadAvatarFailed)
		return
	}

	// 上传头像到阿里云OSS
	path, err := oss.UploadPic(file)
	if err != nil || path == "" {
		// 上传失败
		ResponseError(c, CodeUploadAvatarFailed)
		return
	}

	// 获取用户ID
	idStr := c.GetInt64("uid")
	err = logic.UpdateUserAvatar(idStr, path)
	if err != nil {
		// 上传失败
		ResponseError(c, CodeUploadAvatarFailed)
		return
	}

	ResponseSuccessWithMsg(c, "更新成功✅", nil)
}
