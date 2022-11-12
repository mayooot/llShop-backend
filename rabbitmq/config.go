package rabbitmq

// SMS服务消息队列配置
const (
	SmsExchangeName = "sms_direct_exchange"
	SmsExchangeType = "direct"
	SmsQueueName    = "sms_queue"
	SmsRoutingKey   = "sms_routing_key"
)

// Canal购物车服务消息队列配置
const (
	CanalCartExchangeName = "canal_cart_direct_exchange"
	CanalCartExchangeType = "direct"

	CanalCartInsertQueueName  = "canal_cart_insert_queue"
	CanalCartInsertRoutingKey = "canal_cart_insert_routing_key"

	CanalCartDeleteQueueName  = "canal_cart_delete_queue"
	CanalCartDeleteRoutingKey = "canal_cart_delete_routing_key"

	CanalCartSelectQueueName  = "canal_cart_select_queue"
	CanalCartSelectRoutingKey = "canal_cart_select_routing_key"
)

// 提交订单后，异步删除购物车消息队列配置
const (
	CartDelExchangeName = "cart_del_direct_exchange"
	CartDelExchangeType = "direct"

	CartDeleteQueueName  = "cart_delete_queue"
	CartDeleteRoutingKey = "cart_delete_routing_key"
)

// 订单超时未支付回滚消息队列配置
const (
	OrderExchangeName = "order_exchange"
	OrderExchangeType = "direct"
	OrderQueueName    = "order_queue"
	OrderRoutingKey   = "order_routing_key"

	DelayOrderExchangeName = "delay_order_exchange"
	DelayOrderExchangeType = "direct"
	DelayOrderQueueName    = "delay_order_queue"
	DelayOrderRoutingKey   = "delay_order_routing_key"

	// DelayOrderTTL 订单在30分钟后未支付，就会进入死信队列(单位:毫秒)
	DelayOrderTTL = "1800000"
	// DelayOrderTTL = "60000"
)

// 存入用户秒杀商品请求的消息队列配置
const (
	SecKillReqExchangeName = "seckill_exchange"
	SecKillReqExchangeType = "direct"
	SecKillReqQueueName    = "seckill_queue"
	SecKillReqRoutingKey   = "seckill_routing_key"

	// SecKillStore 请求保存队列，队列长度为库存的两倍。防止前面预定失败，后面继续补上
	SecKillStore = 200
	// SecKillOverflow 队列满了之后的拒绝策略为直接拒绝发布新的消息
	SecKillOverflow = "reject-publish"
)
