package rabbitmq

import (
	"encoding/json"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
)

// DelayOrderReceiver 实现了Receiver接口，负责订单超时未支付回滚
type DelayOrderReceiver struct {
	queueName string
	routerKey string
	e         error
	body      []byte
}

// NewDelayOrderReceiver 初始化一个消费canal信息的mq接收者
func NewDelayOrderReceiver(queueName, routerKey string) *DelayOrderReceiver {
	return &DelayOrderReceiver{
		queueName: queueName,
		routerKey: routerKey,
	}
}

// QueueName 返回队列名称
func (r *DelayOrderReceiver) QueueName() string {
	return r.queueName
}

// RoutingKey 返回RoutingKey
func (r *DelayOrderReceiver) RoutingKey() string {
	return r.routerKey
}

// OnError 将执行过程中产生的异常赋值到接收者
func (r *DelayOrderReceiver) OnError(e error) {
	r.e = e
}

// OnReceive 处理队列中的消息，处理成功返回true，则会应答。否则返回false，会反复处理消息，直到处理成功
func (r *DelayOrderReceiver) OnReceive(body []byte) bool {
	if r.e != nil {
		zap.L().Error("消费订单超时回滚消息出现异常", zap.Error(r.e))
		return false
	}
	var orderNum int64
	_ = json.Unmarshal(body, &orderNum)
	orderStatus, err := mysql.SelectOrderOrderStatus(orderNum)
	if err != nil || orderStatus == 0 {
		// 获取订单状态失败
		return false
	}
	if orderStatus == 5 {
		// 订单已经超时
		return true
	} else if orderStatus == 6 {
		// 订单超时后仍然未支付，将订单状态设置为超时未支付
		err := mysql.UpdateOrderOrderStatus(orderNum, 5)
		if err != nil {
			return false
		}
		// 回滚库存
		err = mysql.RollbackOrderStock(orderNum)
		if err != nil {
			return false
		}
	} else {
		// 订单已支付
		return true
	}
	return true
}
