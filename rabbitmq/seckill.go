package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/models/dto"
)

// SecKillReceiver 实现了Receiver接口，负责消费存入秒杀请求的队列
type SecKillReceiver struct {
	queueName string
	routerKey string
	e         error
	body      []byte
}

// NewSecKillReceiver 初始化一个消费canal信息的mq接收者
func NewSecKillReceiver(queueName, routerKey string) *SecKillReceiver {
	return &SecKillReceiver{
		queueName: queueName,
		routerKey: routerKey,
	}
}

// QueueName 返回队列名称
func (r *SecKillReceiver) QueueName() string {
	return r.queueName
}

// RoutingKey 返回RoutingKey
func (r *SecKillReceiver) RoutingKey() string {
	return r.routerKey
}

// OnError 将执行过程中产生的异常赋值到接收者
func (r *SecKillReceiver) OnError(e error) {
	r.e = e
}

// OnReceive 处理队列中的消息，处理成功返回true，则会应答。否则返回false，会反复处理消息，直到处理成功
func (r *SecKillReceiver) OnReceive(body []byte) bool {
	if r.e != nil {
		zap.L().Error("消费秒杀队列中的消息出现异常", zap.Error(r.e))
		return true
	}
	data := new(dto.SecKillMQ)
	_ = json.Unmarshal(body, &data)
	_ = mysql.UpdateSecKillProductStock(data.SkuID)
	return true
}

// SendSecKillReqMess2MQ 负责发送用户秒杀请求到RabbitMQ
func SendSecKillReqMess2MQ(data *dto.SecKillMQ) error {
	// 转换为json数据
	dataJson, _ := json.Marshal(data)
	// 发送消息
	err = rabbitmqChannel5.Publish(
		SecKillReqExchangeName,
		SecKillReqRoutingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: 2, // 2 表示消息持久化
			ContentType:  "application/json",
			Body:         dataJson,
		},
	)
	if err != nil {
		zap.L().Error("秒杀服务，发送消息到RabbitMQ失败", zap.Error(err))
		return err
	}
	zap.L().Info("秒杀服务，发送消息到RabbitMQ成功")
	return nil
}
