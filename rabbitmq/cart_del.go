package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"shop-backend/dao/mysql"
	"shop-backend/models/dto"
	"strconv"
)

// CartDelReceiver 实现了Receiver接口，负责异步删除用户购物车
type CartDelReceiver struct {
	queueName string
	routerKey string
	e         error
	body      []byte
}

// NewCartDelReceiver 初始化一个消费canal信息的mq接收者
func NewCartDelReceiver(queueName, routerKey string) *CartDelReceiver {
	return &CartDelReceiver{
		queueName: queueName,
		routerKey: routerKey,
	}
}

// QueueName 返回队列名称
func (r *CartDelReceiver) QueueName() string {
	return r.queueName
}

// RoutingKey 返回RoutingKey
func (r *CartDelReceiver) RoutingKey() string {
	return r.routerKey
}

// OnError 将执行过程中产生的异常赋值到接收者
func (r *CartDelReceiver) OnError(e error) {
	r.e = e
}

// OnReceive 处理队列中的消息，处理成功返回true，则会应答。否则返回false，会反复处理消息，直到处理成功
func (r *CartDelReceiver) OnReceive(body []byte) bool {
	if r.e != nil {
		zap.L().Error("消费异步删除购物车消息出现异常", zap.Error(r.e))
		return false
	}
	data := new(dto.CartProductListDTO)
	_ = json.Unmarshal(body, &data)
	uid := data.UserID
	for _, cart := range data.CartProductList {
		skuId, _ := strconv.ParseInt(cart.SkuID, 10, 64)
		err := mysql.DelCartProductBySkuIDAndUID(uid, skuId, cart.Specification)
		if err != nil {
			return false
		}
	}
	return true
}

// SendCartDelMess2MQ 负责发送要删除的购物车商品信息到RabbitMQ
func SendCartDelMess2MQ(data *dto.CartProductListDTO) {
	sendSucc := false
	for !sendSucc {
		// 转换为json数据
		dataJson, _ := json.Marshal(data)
		// 发送消息
		err = rabbitmqChannel3.Publish(
			CartDelExchangeName,
			CartDeleteRoutingKey,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: 2, // 2 表示消息持久化
				ContentType:  "application/json",
				Body:         dataJson,
			},
		)
		if err != nil {
			zap.L().Error("异步删除购物车服务，发送消息到RabbitMQ失败", zap.Error(err))
		}
		sendSucc = true
		zap.L().Info("异步删除购物车服务，发送消息到RabbitMQ成功")
	}
}
