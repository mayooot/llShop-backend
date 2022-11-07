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
)
