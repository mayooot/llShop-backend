package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"shop-backend/logic"
	"shop-backend/models"
	"shop-backend/utils/check"
	"strconv"
)

// SendVerifyCodeHandler 发送手机验证码，并返回
// @Summary 获取验证码
// @Description 前端传递11位手机号，后端随机生成4位验证码。存入Redis，并异步调用阿里云SMS服务发送验证码。
// @Tags 用户相关接口
// @Produce  json
// @Param phone query string true "手机号"
// @Router /phone [get]
func SendVerifyCodeHandler(c *gin.Context) {
	// 获取参数
	phone := c.Query("phone")
	if phone == "" {
		// phone字段为空
		ResponseError(c, CodePhoneIsNotEmpty)
		return
	}
	if ok := check.VerifyMobileFormat(phone); !ok {
		// 手机号格式不正确
		ResponseError(c, CodePhoneFormatError)
		return
	}
	// 发送验证码
	code, err := logic.SendVerifyCode(phone)
	if err != nil {
		if errors.Is(err, logic.ErrorRequestCodeFrequent) {
			// 用户频繁请求验证码
			ResponseError(c, CodeRequestCodeFrequently)
			return
		}
		// 生成、发送验证码失败
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccessWithMsg(c, "发送成功，验证码五分钟内有效", gin.H{
		"code": code,
	})
}

// SignUpHandler 用户注册
// @Summary 注册用户
// @Description 前端传递JSON类型对象，后端完成校验后注册新用户。
// @Tags 用户相关接口
// @Produce json
// @Param ParamSignUp body models.ParamSignUp true "用户注册结构体"
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	// 获取参数并校验
	p := new(models.ParamSignUp)
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
// @Param ParamLogin body models.ParamLogin true "用户登录结构体"
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	// 获取参数并校验
	p := new(models.ParamLogin)
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
		}
		ResponseError(c, CodeServeBusy)
		return
	}
	fmt.Println(math.MaxInt64)
	ResponseSuccessWithMsg(c, "登录成功", gin.H{
		// 前端json能接受的整数范围为 - (2^53 -1) ~ 2^53 - 1，而我们要传递的是int64，所以要转化成字符串
		"userId":       strconv.FormatInt(uid, 10),
		"accessToken":  aToken,
		"refreshToken": rToken,
	})
}

// SomeInfoHandler 获取用户头像、用户名和购物车数量
// @Summary 获取用户头像、用户名和购物车数量。
// @Description 前端传递用户ID，后端返回用户头像、用户名和购物车数量。
// @Tags 用户相关接口
// @Produce json
// @Security x-token
// @param Authorization header string true "Bearer token"
// @Param id path string true "用户ID"
// @Router /someinfo/{id} [get]
func SomeInfoHandler(c *gin.Context) {
	// 获取用户id
	idStr := c.Param("id")
	// 将字符串转成int64
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取用户简略信息
	infos, err := logic.GetSomeInfo(uid)
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, infos)
}

// UserInfosHandler 获取用户详细的个人信息
// @Summary 获取用户个人信息
// @Description 前端传递用户ID，后端返回用户个人信息。
// @Tags 用户相关接口
// @Produce json
// @Security x-token
// @param Authorization header string true "Bearer token"
// @Param id path string true "用户ID"
// @Router /infos/{id} [get]
func UserInfosHandler(c *gin.Context) {
	// 获取用户id
	idStr := c.Param("id")
	// 字符串转为int64
	uid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	infos, err := logic.GetUserInfos(uid)
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
// @param Authorization header string true "Bearer token"
// @Param param body models.ParamInfos true "用户个人信息结构体"
// @Router /infos/update [put]
func UserInfosUpdateHandler(c *gin.Context) {
	infos := new(models.ParamInfos)
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
	if infos.Id != c.GetString("uid") {
		// 如果用户传递的uid和上一步校验jwt中间件中的uid不同
		// 请求参数有误
		ResponseError(c, CodeServeBusy)
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
// @Description 前端传递UserID，后端清空Redis中AccessToken。
// @Tags 用户相关接口
// @Produce  json
// @Security x-token
// @param Authorization header string true "Bearer token"
// @Param id path string true "用户ID"
// @Router /exit/{id} [delete]
func SignOutHandler(c *gin.Context) {
	idStr := c.Param("id")
	err := logic.SignOut(idStr)
	if err != nil {
		ResponseError(c, CodeSignOutFailed)
		return
	}
	ResponseSuccessWithMsg(c, "退出成功", nil)
}

// RefreshToken 刷新AccessToken
// @Summary 刷新AccessToken
// @Description 前端收到code为1019的状态码后，应重新发起一个put请求访问该接口。
// @Tags 用户相关接口
// @Produce  json
// @Security x-token
// @param Authorization header string true "Bearer AToken&RToken"
// @Param id path string true "用户ID"
// @Router /refreshToken/{id} [put]
func RefreshToken(c *gin.Context) {

}
