package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"shop-backend/logic"
	"shop-backend/models"
	"shop-backend/utils/check"
)

// SendVerifyCodeHandler 发送手机验证码，并返回
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

	ResponseSuccess(c, gin.H{"code": code})
}

// SignUpHandler 用户注册
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
	ResponseSuccess(c, nil)
}

// LoginHandler 用户登录
func LoginHandler(c *gin.Context) {
	// 获取参数并校验

}
