package controller

// ResCode 错误码
type ResCode int64

const (
	CodeSuccess         ResCode = 200
	CodeServeBusy       ResCode = 500
	CodePhoneIsNotEmpty ResCode = 1000 + iota
	CodePhoneFormatError
	CodeEmailFormatError
	CodeInvalidParams
	CodeWrongVerifyCode
	CodePassIsWeak
	CodeRequestCodeFrequently
	CodeUserIsRegistered
	CodeTokenIsEmpty
	CodeTokenIsWrongFormat
	CodeTokenIsInvalid
	CodeUpdateInfosFailed
	CodeUsernameToLongOrToShort
	CodeExceedMaxTerminalNum
	CodeSignOutFailed
	CodeTokenRefreshFailed
	CodeAccessTokenIsLiving
	CodeTokenExpire
	CodeUsernameOrPassError
	CodeUserNotExist
	CodeUploadAvatarFailed
	CodeUploadAvatarToBigOrExtError
	CodeMustRequestCode
	CodeNeedReLogin
	CodeFrontEndNeedUseNewToken
	CodeRequestAllCategoryFailed
	CodeRequestAllAttributeFailed
	CodeSearchConditionIsNil
	CodeDeleteCartProductFailed
	CodeUpdateCartProductStatusFailed
	CodeCreatePreSubmitOrderSuccess
	CodeAddReceiverAddressFailed
	CodeUpdateReceiverAddressFailed
	CodeOrderNumISNotExistOrExpired
	CodeCreateSubmitOrderSuccess
	CodeToManyRequest
	CodeSecKillFinished
)

// map字典 K: 错误码	V: 错误信息
var codeMsgMap = map[ResCode]string{

	CodeSuccess:                       "success",
	CodeServeBusy:                     "服务器繁忙，等会再试试吧~🧸",
	CodePhoneIsNotEmpty:               "手机号未输入或为空❌",
	CodePhoneFormatError:              "手机号格式错误❌",
	CodeEmailFormatError:              "手机号格式错误❌",
	CodeInvalidParams:                 "请求参数有误❌",
	CodeWrongVerifyCode:               "验证码错误或已过期❌",
	CodePassIsWeak:                    "密码强度太弱啦~🤗",
	CodeRequestCodeFrequently:         "验证码已发送，请注意查收~🐹",
	CodeUserIsRegistered:              "用户已注册，请直接登录👻",
	CodeTokenIsEmpty:                  "请求未携带Token❌",
	CodeTokenIsWrongFormat:            "携带Token的格式有误❌",
	CodeTokenIsInvalid:                "非法Token❌",
	CodeUpdateInfosFailed:             "更新个人资料失败，请稍后再试😪",
	CodeUsernameToLongOrToShort:       "用户名太长或太短😥",
	CodeExceedMaxTerminalNum:          "超过最大登录终端数量",
	CodeSignOutFailed:                 "退出失败，等会再试试吧😪",
	CodeTokenRefreshFailed:            "刷新Token失败",
	CodeAccessTokenIsLiving:           "AccessToken未过期，刷新Token失败",
	CodeTokenExpire:                   "token已过期",
	CodeUsernameOrPassError:           "账户或密码错误🥵",
	CodeUserNotExist:                  "用户不存在",
	CodeUploadAvatarFailed:            "上传头像失败🫥",
	CodeUploadAvatarToBigOrExtError:   "图片过大或格式不正确😣",
	CodeMustRequestCode:               "请先获取验证码",
	CodeNeedReLogin:                   "认证过期，请重新登录😎",
	CodeFrontEndNeedUseNewToken:       "请重置用户的AccessToken",
	CodeRequestAllCategoryFailed:      "获取商品分类信息失败",
	CodeRequestAllAttributeFailed:     "获取商品属性信息失败",
	CodeSearchConditionIsNil:          "搜索条件为空",
	CodeDeleteCartProductFailed:       "删除失败",
	CodeUpdateCartProductStatusFailed: "选择购物车商品失败",
	CodeCreatePreSubmitOrderSuccess:   "创建订单成功🧪",
	CodeAddReceiverAddressFailed:      "添加收货地址失败，等会再试试吧😴",
	CodeUpdateReceiverAddressFailed:   "修改收货地址失败，等会再试试吧😴",
	CodeOrderNumISNotExistOrExpired:   "请刷新预提交订单🪬",
	CodeCreateSubmitOrderSuccess:      "订单提交成功🐔",
	CodeToManyRequest:                 "当前活动太火爆啦，等会再试试吧🍻",
	CodeSecKillFinished:               "秒杀活动已结束，谢谢参与😮",
}

// Msg 为ResCode注册一个Msg方法，负责返回错误码对应的错误信息
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServeBusy]
	}
	return msg
}
