package canal

import (
	"github.com/withlin/canal-go/client"
	"go.uber.org/zap"
	"shop-backend/rabbitmq"
	"shop-backend/settings"
	"time"
)

// Init 初始化并连接canal服务
func Init(cfg *settings.CanalConfig) {
	// 构建canal连接URL
	connector := client.NewSimpleCanalConnector(cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Destination, 60000, 60*60*1000)
	err := connector.Connect()
	if err != nil {
		panic("CanalConfig 初始化失败: " + err.Error())
	}
	zap.L().Info("初始化canal服务成功")

	// 监听shop库下的所有表
	err = connector.Subscribe("shop\\..*")
	if err != nil {
		panic("CanalConfig 监听shop库失败: " + err.Error())
	}
	zap.L().Info("开始监听shop库数据变更")

	// 监听数据库变更信息
	for {
		// 每次获取一百条变更数据
		message, err := connector.Get(100, nil, nil)
		if err != nil {
			zap.L().Error("canal监听shop库，获取数据库变更信息失败", zap.Error(err))
		}
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			// 如果获取到的消息集合为空，1秒后再次获取
			time.Sleep(1 * time.Second)
			continue
		}
		// 将消息发送到RabbitMQ中
		_ = rabbitmq.SendDBInfo2MQ(message.Entries)
	}
}
