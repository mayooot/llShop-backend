package pay

import (
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"shop-backend/settings"
)

var privateKey string
var appId string
var Client *alipay.Client
var err error

func Init(cfg *settings.AliPayConfig) {
	privateKey = cfg.PrivateKey
	appId = cfg.AppID
	Client, err = alipay.New(appId, privateKey, false)
	if err != nil {
		panic("初始化支付宝支付模块失败")
	}
	err = Client.LoadAliPayPublicKey(cfg.PublicKey)
	if err != nil {
		panic("初始化支付宝支付模块失败")
		return
	}
}

// AliPay 生成支付宝网站支付页面
func AliPay(orderNum, totalAmount string) string {
	// 电脑网站支付
	var p = alipay.TradePagePay{}
	p.ReturnURL = "http://172.20.10.4:8888/paysuccess"
	p.Subject = "支付宝在线支付"
	p.OutTradeNo = orderNum
	p.TotalAmount = totalAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	// 这里返回的url中会包含sign，直接返回给前端就ok
	url, err := Client.TradePagePay(p)
	if err != nil {
		zap.L().Error("生成支付宝的支付页面失败", zap.Error(err))
	}
	return url.String()
}
