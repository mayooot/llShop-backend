package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"shop-backend/utils/sms"
)

// SmsReceiver 实现了Receiver接口
type SmsReceiver struct {
	queueName string
	routerKey string
	e         error
	body      []byte
}

// NewSmsReceiver 初始化一个消费短信信息的mq接收者
func NewSmsReceiver(queueName, routerKey string) *SmsReceiver {
	return &SmsReceiver{
		queueName: queueName,
		routerKey: routerKey,
	}
}

// QueueName 返回队列名称
func (r *SmsReceiver) QueueName() string {
	return r.queueName
}

// RoutingKey 返回RoutingKey
func (r *SmsReceiver) RoutingKey() string {
	return r.routerKey
}

// OnError 将执行过程中产生的异常赋值到接收者
func (r *SmsReceiver) OnError(e error) {
	r.e = e
}

// OnReceive 处理队列中的消息，处理成功返回true，则会应答。否则返回false，会反复处理消息，直到处理成功
// 这里的接收者是负责处理短信发送，所以不需要反复处理，如果一次发送失败，让用户再次获取验证码
func (r *SmsReceiver) OnReceive(body []byte) bool {
	if r.e != nil {
		zap.L().Error("SMS短信服务，SmsReceiver中有异常", zap.Error(r.e))
		return true
	}
	data := make(map[string]string)
	// 反序列json
	err := json.Unmarshal(body, &data)
	if err != nil {
		zap.L().Error("SMS短信服务，解析json失败", zap.Error(err))
		return true
	}
	err = sms.SendSms(data["phone"], data["code"])
	if err != nil {
		zap.L().Error("SMS短信服务，调用阿里云SMS发送短信失败", zap.Error(err),
			zap.String("phone", data["phone"]))
		return true
	}
	zap.L().Info("SMS短信服务，调用阿里云SMS发送短信成功", zap.Error(err),
		zap.String("phone", data["phone"]),
		zap.String("phone", data["phone"]),
	)
	return true
}

// SendSms 发送手机号和验证码到RabbitMQ中
func SendSms(phone, code string) error {
	// 使用map存储手机号和验证码
	data := map[string]string{"phone": phone, "code": code}
	// 转换为json数据
	dataJson, err := json.Marshal(data)
	if err != nil {
		zap.L().Error("SMS短信服务，手机号和验证码转换为json格式失败", zap.Error(err))
		return err
	}
	// 发送消息
	err = rabbitmqChannel.Publish(
		SmsExchangeName,
		SmsRoutingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: 2, // 2 表示消息持久化
			ContentType:  "application/json",
			Body:         dataJson,
		},
	)
	if err != nil {
		zap.L().Error("SMS短信服务，发送消息到RabbitMQ失败", zap.Error(err))
		return err
	}
	return nil
}
