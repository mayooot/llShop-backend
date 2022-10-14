package controller

// ResCode 错误码
type ResCode int64

const (
	CodeSuccess         ResCode = 200
	CodeServeBusy       ResCode = 500
	CodePhoneIsNotEmpty ResCode = 1000 + iota
	CodePhoneFormatError
	CodeInvalidParams
	CodeWrongVerifyCode
	CodePassIsWeak
	CodeRequestCodeFrequently
	CodeUserIsRegistered
)

// map字典 K: 错误码	V: 错误信息
var codeMsgMap = map[ResCode]string{

	CodeSuccess:               "success",
	CodeServeBusy:             "服务器繁忙，等会再试试吧~🧸",
	CodePhoneIsNotEmpty:       "手机号未输入或为空❌",
	CodePhoneFormatError:      "手机号格式错误❌",
	CodeInvalidParams:         "请求参数格式错误❌",
	CodeWrongVerifyCode:       "验证码错误或已过期❌",
	CodePassIsWeak:            "密码强度太弱啦~🤗",
	CodeRequestCodeFrequently: "验证码已发送，请注意查收~🐹",
	CodeUserIsRegistered:      "用户已注册，请直接登录👻",
}

// Msg 为ResCode注册一个Msg方法，负责返回错误码对应的错误信息
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[c]
	}
	return msg
}
