package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"shop-backend/settings"
)

var rabbitmqConn *amqp.Connection
var rabbitmqChannel *amqp.Channel
var rabbitmqChannel2 *amqp.Channel
var rabbitmqChannel3 *amqp.Channel
var rabbitmqChannel4 *amqp.Channel

var rabbitmqChannel5 *amqp.Channel
var err error

// Init 初始化RabbitMQ
func Init(cfg *settings.RabbitMQConfig) {
	// 构造RabbitMQ连接url
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port)
	// 连接RabbitMQ
	rabbitmqConn, err = amqp.Dial(url)
	if err != nil {
		panic("RabbitMQConfig 初始化失败: " + err.Error())
	}

	// 初始化channel
	rabbitmqChannel, err = rabbitmqConn.Channel()
	if err != nil {
		panic("打开Channel失败: " + err.Error())
	}

	rabbitmqChannel2, err = rabbitmqConn.Channel()
	if err != nil {
		panic("打开Channel失败: " + err.Error())
	}

	rabbitmqChannel3, err = rabbitmqConn.Channel()
	if err != nil {
		panic("打开Channel失败: " + err.Error())
	}

	rabbitmqChannel4, err = rabbitmqConn.Channel()
	if err != nil {
		panic("打开Channel失败: " + err.Error())
	}

	rabbitmqChannel5, err = rabbitmqConn.Channel()
	if err != nil {
		panic("打开Channel失败: " + err.Error())
	}

	// 初始化短信相关的RabbitMQ实体对象
	smsMQ := NewSmsMQ()
	// 将管道绑定到MQ对象上
	smsMQ.channel = rabbitmqChannel
	// 创建接受短信信息的接收者
	smsReceiver := NewSmsReceiver(SmsQueueName, SmsRoutingKey)
	// 将接收者绑定到RabbitMQ实体对象
	smsMQ.RegisterReceiver(smsReceiver)
	// 启动
	go smsMQ.Start()

	// 初始化canal购物车相关的RabbitMQ实体对象
	canalCart := NewCanalCartMQ()
	// 将管道绑定到MQ对象上
	canalCart.channel = rabbitmqChannel2
	// 创建接受数据库变更信息的接收者
	cartReceiver1 := NewCanalCartReceiver(CanalCartInsertQueueName, CanalCartInsertRoutingKey)
	cartReceiver2 := NewCanalCartReceiver(CanalCartDeleteQueueName, CanalCartDeleteRoutingKey)
	cartReceiver3 := NewCanalCartReceiver(CanalCartSelectQueueName, CanalCartSelectRoutingKey)
	// 将接收者绑定到RabbitMQ实体对象
	canalCart.RegisterReceiver(cartReceiver1)
	canalCart.RegisterReceiver(cartReceiver2)
	canalCart.RegisterReceiver(cartReceiver3)
	// 启动
	go canalCart.Start()

	// 初始化异步删除购物车相关的RabbitMQ实体对象
	cartDel := NewCartDelMQ()
	// 将管道绑定到MQ对象上
	cartDel.channel = rabbitmqChannel3
	// 创建接受数据库变更信息的接收者
	cartDelReceiver := NewCartDelReceiver(CartDeleteQueueName, CartDeleteRoutingKey)
	// 将接收者绑定到RabbitMQ实体对象
	cartDel.RegisterReceiver(cartDelReceiver)
	// 启动
	go cartDel.Start()

	// 初始化订单相关的RabbitMQ实体对象
	order := NewOrderMQ()
	// 将管道绑定到MQ对象上
	order.channel = rabbitmqChannel4
	// 暂时不需要消费普通订单队列中的消息
	// orderReceiver := NewOrderReceiver(OrderQueueName, OrderRoutingKey)
	// 将接收者绑定到RabbitMQ实体对象
	// order.RegisterReceiver(orderReceiver)
	go order.Start()

	// 初始化订单超时回滚相关的RabbitMQ实体对象
	delayOrder := NewDelayOrderMQ()
	// 将管道绑定到MQ对象上
	delayOrder.channel = rabbitmqChannel4
	// 创建接受数据库变更信息的接收者
	delayOrderReceiver := NewDelayOrderReceiver(DelayOrderQueueName, DelayOrderRoutingKey)
	// 将接收者绑定到RabbitMQ实体对象
	delayOrder.RegisterReceiver(delayOrderReceiver)
	// 启动
	go delayOrder.Start()

	// 初始化秒杀商品相关的RabbitMQ实体对象
	secKill := NewSecKillMQ()
	// 将管道绑定到MQ对象上
	secKill.channel = rabbitmqChannel5
	// 创建接受数据库变更信息的接收者
	secKillReceiver := NewSecKillReceiver(SecKillReqQueueName, SecKillReqRoutingKey)
	// 将接收者绑定到RabbitMQ实体对象
	secKill.RegisterReceiver(secKillReceiver)
	// 启动
	go secKill.Start()
}

// Destroy 销毁RabbitMQ连接和通道
func Destroy() {
	rabbitmqConn.Close()
	rabbitmqChannel.Close()
	rabbitmqChannel2.Close()
}
