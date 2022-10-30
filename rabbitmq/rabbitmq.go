package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
	"time"
)

// Receiver 观察者模式需要的接口(观察者)
// 观察者用于接收指定的queue到来的数据
type Receiver interface {
	QueueName() string     // 获取接收者需要监听的队列
	RoutingKey() string    // 这个队列绑定的路由
	OnError(error)         // 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
	OnReceive([]byte) bool // 处理收到的消息, 这里需要告知RabbitMQ对象消息是否处理成功
}

// RabbitMQ 用于管理和维护RabbitMQ的对象(被观察者)
type RabbitMQ struct {
	wg           sync.WaitGroup
	channel      *amqp.Channel
	exchangeName string // exchange的名称
	exchangeType string // exchange的类型
	receivers    []Receiver
}

// NewSmsMQ 创建一个用于发送短信业务的新的操作RabbitMQ的对象
func NewSmsMQ() *RabbitMQ {
	return &RabbitMQ{
		exchangeName: SmsExchangeName,
		exchangeType: SmsExchangeType,
	}
}

// 准备RabbitMQ的交换机
func (mq *RabbitMQ) prepareExchange() error {
	// 声明交换机
	err := mq.channel.ExchangeDeclare(
		mq.exchangeName, // exchange
		mq.exchangeType, // type
		true,            // durable
		false,           // autoDelete
		false,           // internal
		false,           // noWait
		nil,             // args
	)
	if err != nil {
		zap.L().Error("RabbitMQ直接交换机失败", zap.Error(err))
		return err
	}
	return nil
}

// run 开始获取连接并初始化相关操作
func (mq *RabbitMQ) run() {
	// 初始化Exchange
	mq.prepareExchange()

	for _, receiver := range mq.receivers {
		// 一个RabbitMQ对象可以对应多个消费者，每个消费者的加入都会使得WaitGroup+1
		mq.wg.Add(1)
		// 每个接收者单独启动一个goroutine用来初始化queue并接收消息
		// listen方法完成时会执行mq.wg.Done()
		zap.L().Info("开启一个新的协程处理RabbitMQ队列中的消息")
		go mq.listen(receiver)
	}

	// 一直在此等待，除非所有的消费者都意外退出
	mq.wg.Wait()

	zap.L().Error("所有处理queue的任务都意外退出了")

	// 理论上mq.run()在程序的执行过程中是不会结束的
	// 一旦结束就说明所有的接收者都退出了，那么意味着程序与RabbitMQ的连接断开
	// 那么则需要重新连接，这里尝试销毁当前连接
	Destroy()
}

// Start 启动RabbitMQ的客户端
func (mq *RabbitMQ) Start() {
	for {
		mq.run()
		// 一旦连接断开，那么需要隔一段时间去重连
		// 这里最好有一个时间间隔
		time.Sleep(3 * time.Second)
	}
}

// RegisterReceiver 注册一个用于接收指定队列指定路由的数据接收者
// 将若干个观察者聚集到被观察者中，实现类间解耦
func (mq *RabbitMQ) RegisterReceiver(receiver Receiver) {
	mq.receivers = append(mq.receivers, receiver)
}

// Listen 监听指定路由发来的消息
// 这里需要针对每一个接收者启动一个goroutine来执行listen
// 该方法负责从每一个接收者监听的队列中获取数据，并负责重试
func (mq *RabbitMQ) listen(receiver Receiver) {
	defer mq.wg.Done()

	// 这里获取每个接收者需要监听的队列和路由
	queueName := receiver.QueueName()
	routerKey := receiver.RoutingKey()

	// 声明Queue
	_, err := mq.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive(排他性队列)
		false,     // no-wait
		nil,       // arguments
	)
	if nil != err {
		// 当队列初始化失败的时候，需要告诉这个接收者相应的错误
		receiver.OnError(fmt.Errorf("初始化队列 %s 失败: %s", queueName, err.Error()))
	}

	// 将Queue绑定到Exchange上去
	err = mq.channel.QueueBind(
		queueName,       // queue name
		routerKey,       // routing key
		mq.exchangeName, // exchange
		false,           // no-wait
		nil,
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("绑定队列 [%s - %s] 到交换机失败: %s", queueName, routerKey, err.Error()))
	}

	// 消费者流控
	mq.channel.Qos(
		1,    // 当前消费者一次能接受的最大消息数量
		0,    // 服务器传递的最大容量(以八字节为单位)
		true) // 设置为true，对全局channel可用

	msgs, err := mq.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack 关闭自动应答
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("获取队列 %s 的消费通道失败: %s", queueName, err.Error()))
	}

	// 使用callback消费数据
	for msg := range msgs {
		// 当接收者消息处理失败的时候，
		// 比如网络问题导致的数据库连接失败，redis连接失败等等这种
		// 通过重试可以成功的操作，那么这个时候是需要重试的
		// 直到数据处理成功后再返回，然后才会回复rabbitmq ack
		for !receiver.OnReceive(msg.Body) {
			zap.L().Error("receiver 数据处理失败，将要重试")
			time.Sleep(1 * time.Second)
		}

		// 如果为true表示确认所有未确认的消息,一般用在批量消费中
		// 为false表示确认当前消息,rq就会删除
		msg.Ack(false)
	}
}
