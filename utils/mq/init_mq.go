package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"shop-backend/settings"
	"shop-backend/utils/mq/sms"
)

var Conn *amqp.Connection

// Init 初始化RabbitMQ连接
func Init(cfg *settings.RabbitMQ) (err error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port)
	Conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	err = sms.SmsConfig()
	if err != nil {
		zap.L().Error("初始化SMS RabbitMQ短信服务失败")
		return err
	}
	return nil
}
