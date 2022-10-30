package sms

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"shop-backend/settings"
	"shop-backend/utils/concatstr"
)

var client *dysmsapi20170525.Client

// Init 初始化阿里云SMS服务
func Init(cfg *settings.AliyunConfig) (err error) {
	config := &openapi.Config{
		AccessKeyId:     &cfg.AccessKeyId2,
		AccessKeySecret: &cfg.AccessKeySecret2,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	// 初始化短信发送客户端
	client, err = dysmsapi20170525.NewClient(config)
	return err
}

// SendSms 发送短信
func SendSms(phone, code string) (err error) {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:     tea.String("阿里云短信测试"),
		TemplateCode: tea.String("SMS_154950909"),
		// 手机号
		PhoneNumbers: tea.String(phone),
		// 验证码
		TemplateParam: tea.String(concatstr.ConcatString("{\"code\":\"", code, "\"}")),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				e = r
			}
		}()
		// 发送短信
		_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		_, err = util.AssertAsString(error.Message)
		if err != nil {
			return err
		}
	}
	return err
}
