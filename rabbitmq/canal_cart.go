package rabbitmq

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"go.uber.org/zap"
	"shop-backend/dao/redis"
	"shop-backend/models/pojo"
	"shop-backend/models/vo"
	"shop-backend/utils/build"
	"strconv"
	"time"
)

// CanalCartReceiver 实现了Receiver接口
type CanalCartReceiver struct {
	queueName string
	routerKey string
	e         error
	body      []byte
}

// NewCanalCartReceiver 初始化一个消费canal信息的mq接收者
func NewCanalCartReceiver(queueName, routerKey string) *CanalCartReceiver {
	return &CanalCartReceiver{
		queueName: queueName,
		routerKey: routerKey,
	}
}

// QueueName 返回队列名称
func (r *CanalCartReceiver) QueueName() string {
	return r.queueName
}

// RoutingKey 返回RoutingKey
func (r *CanalCartReceiver) RoutingKey() string {
	return r.routerKey
}

// OnError 将执行过程中产生的异常赋值到接收者
func (r *CanalCartReceiver) OnError(e error) {
	r.e = e
}

// OnReceive 处理队列中的消息，处理成功返回true，则会应答。否则返回false，会反复处理消息，直到处理成功
func (r *CanalCartReceiver) OnReceive(body []byte) bool {
	if r.e != nil {
		zap.L().Error("消费canal购物车消息出现异常", zap.Error(r.e))
		return false
	}
	data := new(pojo.Cart)
	// 反序列json
	err := json.Unmarshal(body, &data)
	if err != nil {
		zap.L().Error("canal购物车服务，解析json失败", zap.Error(err))
		return false
	}

	if r.queueName == CanalCartInsertQueueName {
		// 新增或更新类型的变更消息
		// 创建一个存入购物车商品展示对象的通道，缓存区为1
		channel := make(chan *vo.CartProductVO, 1)
		defer close(channel)
		// 添加到Redis缓存中
		if err = redis.AddCartProduct(data.UserID, data.SkuID, build.CreateCartProductVO(data, channel)); err != nil {
			return false
		}
	} else if r.queueName == CanalCartDeleteQueueName {
		// 删除类型的变更消息
		if err = redis.DelCartProduct(data.UserID, data.SkuID); err != nil {
			return false
		}
	} else {
		//
	}
	return true
}

// SendDBInfo2MQ 将数据库变更信息发送到MQ中
func SendDBInfo2MQ(entrys []pbe.Entry) error {
	// 循环消费数据库变更消息
	for _, entry := range entrys {
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			// 如果信息类型为事务开始和事务停止，就跳过此条信息
			continue
		}
		// 解析变更信息
		rowChange := new(pbe.RowChange)
		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		if err != nil {
			zap.L().Error("反序列化canal消息失败", zap.Error(err))
			continue
		}

		// 行变更信息不为空
		if rowChange != nil {
			// 获取行变更数据的类型(增、删、改、查)
			eventType := rowChange.GetEventType()
			for _, rowData := range rowChange.GetRowDatas() {
				// 到此变更数据都为单条
				if eventType == pbe.EventType_UPDATE || eventType == pbe.EventType_INSERT {
					// 变更信息为更新/新增数据
					sendMess2Queue(rowData.GetAfterColumns(), CanalCartInsertRoutingKey)
				} else if eventType == pbe.EventType_DELETE {
					// 变更信息为删除数据
					sendMess2Queue(rowData.GetBeforeColumns(), CanalCartDeleteRoutingKey)
				} else {
					continue
				}
			}
		}
	}
	return nil
}

// 负责发送更新或新增类型的数据库变更数据
func sendMess2Queue(columns []*pbe.Column, routingKey string) {
	// 构建要发送到MQ的对象
	data := buildCartVO(columns)
	if data == nil {
		// 如果对象为空，则不发送
		return
	}
	// 转换为json数据
	dataJson, err := json.Marshal(data)
	if err != nil {
		zap.L().Error("canal购物车服务，购物车商品序列化为json失败", zap.Error(err))
		return
	}

	// 发送消息
	err = rabbitmqChannel2.Publish(
		CanalCartExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: 2, // 2 表示消息持久化
			ContentType:  "application/json",
			Body:         dataJson,
		},
	)
	if err != nil {
		zap.L().Error("canal购物车服务，发送消息到RabbitMQ失败", zap.Error(err))
		return
	}
	zap.L().Info("canal购物车服务，发送消息到RabbitMQ成功")
	return
}

// 构建购物车商品对象
func buildCartVO(columns []*pbe.Column) *pojo.Cart {
	product := new(pojo.Cart)
	for _, col := range columns {
		key := col.GetName()
		value := col.GetValue()

		if key == "user_id" {
			uid, _ := strconv.ParseInt(value, 10, 64)
			product.UserID = uid
		} else if key == "specification" {
			product.Specification = value
		} else if key == "sku_id" {
			skuID, _ := strconv.ParseInt(value, 10, 64)
			product.SkuID = skuID
		} else if key == "count" {
			count, _ := strconv.Atoi(value)
			product.Count = count
		} else if key == "selected" {
			selected, _ := strconv.Atoi(value)
			product.Selected = int8(selected)
		} else if key == "created_time" {
			local, _ := time.LoadLocation("Asia/Shanghai")
			showTime, _ := time.ParseInLocation("2006-01-02 15:04:05", value, local)
			product.CreatedTime = showTime
		} else {
			continue
		}
	}
	return product
}
