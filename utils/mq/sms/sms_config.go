package sms

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"shop-backend/utils/mq"
)

const (
	SmsDirectQueueName = "sms_direct_queue"
	SmsExchange        = "sms_exchange"
	SmsRoutingKey      = "sms_routing_key"
)

var smsChannel *amqp.Channel

// SmsConfig 初始化SMS使用到的消息队列配置
func SmsConfig() (err error) {
	// 创建一个消息队列通道
	smsChannel, err = mq.Conn.Channel()
	if err != nil {
		zap.L().Error("SMS短信服务，创建RabbitMQ通道失败", zap.Error(err))
		return err
	}

	// 声明一个直接交换机
	err = smsChannel.ExchangeDeclare(
		SmsExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Error("SMS短信服务，创建RabbitMQ直接交换机失败", zap.Error(err), zap.String("exchange", SmsExchange))
		return err
	}

	// 声明一个队列
	var q amqp.Queue
	q, err = smsChannel.QueueDeclare(
		SmsDirectQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Error("SMS短信服务，创建RabbitMQ队列失败", zap.Error(err), zap.String("queue", SmsDirectQueueName))
		return err
	}

	// 绑定直接交换机和队列
	err = smsChannel.QueueBind(
		q.Name,
		SmsRoutingKey,
		SmsExchange,
		false,
		nil,
	)
	if err != nil {
		zap.L().Error("SMS短信服务，绑定直接交换机和队列失败", zap.Error(err),
			zap.String("exchange", SmsExchange),
			zap.String("queue", SmsDirectQueueName),
			zap.String("routing-key", SmsRoutingKey),
		)
		return err
	}
	return nil
}
