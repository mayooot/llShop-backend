package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// OrderReceiver 实现了Receiver接口，负责订单相关操作
type OrderReceiver struct {
	queueName string
	routerKey string
	e         error
	body      []byte
}

// NewOrderReceiver 初始化一个消费canal信息的mq接收者
func NewOrderReceiver(queueName, routerKey string) *OrderReceiver {
	return &OrderReceiver{
		queueName: queueName,
		routerKey: routerKey,
	}
}

// QueueName 返回队列名称
func (r *OrderReceiver) QueueName() string {
	return r.queueName
}

// RoutingKey 返回RoutingKey
func (r *OrderReceiver) RoutingKey() string {
	return r.routerKey
}

// OnError 将执行过程中产生的异常赋值到接收者
func (r *OrderReceiver) OnError(e error) {
	r.e = e
}

// OnReceive 处理队列中的消息，处理成功返回true，则会应答。否则返回false，会反复处理消息，直到处理成功
func (r *OrderReceiver) OnReceive(body []byte) bool {
	zap.L().Info("orderReceiver接收到了消息")
	return true
}

// SendDelayOrderMess2MQ 发送订单延时消息到RabbitMQ，超时未支付后回滚库存和修改订单状态为超时未支付
func SendDelayOrderMess2MQ(orderNum int64) {
	sendSucc := false
	for !sendSucc {
		// 转换为json数据
		dataJson, _ := json.Marshal(orderNum)
		// 发送消息
		err = rabbitmqChannel4.Publish(
			OrderExchangeName,
			OrderRoutingKey,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: 2, // 2 表示消息持久化
				ContentType:  "application/json",
				Body:         dataJson,
				Expiration:   DelayOrderTTL,
			},
		)
		if err != nil {
			zap.L().Error("订单超时回滚服务，发送消息到RabbitMQ失败", zap.Error(err))
		}
		sendSucc = true
		zap.L().Info("订单超时回滚服务，发送消息到RabbitMQ成功")
	}
}
