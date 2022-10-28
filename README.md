## llShop-backend

🍨llShop是一个前后端分离的商城项目。主要包括以下模块:

* 🧍🏻用户模块
* 🗒商品信息模块
* 📱首页展示模式
* 🔦搜索模块
* 🪗购物车模块
* 💸️订单模块

🪄技术栈：

* 前端: Vue2 + VueX + ElementUI
* 后端: Gin + Gorm + MySQL + Canal + Redis + RabbitMQ + Elasticsearch + MongoDB

🎨第三方库:

* viper: 配置文件的读取，修改时自动加载新的配置
* zap: 日志记录
* swaggo: 生成Swagger接口文档
* validator: 参数校验
* air: 程序热启动
* snowflake: 生成分布式唯一ID
* jwt-go: 前后端身份认证
* fsnotify: 监听文件或目录，配合viper使用
* lumberjack: 配置zap，实现日志的滚动记录
* ratelimit: 令牌桶限流

🎏目录结构:

~~~text
├───controller              控制层
├───dao                     持久层
│   ├───mysql               MySQL和Gorm初始化
│   └───redis               Redis初始化
├───docs                    Swagger文档
├───logger                  zap日志库配置
├───logic                   业务逻辑层
├───middleware              中间件层
├───models                  结构体模型层
├───router                  路由管理层
├───settings                配置信息层
├───utils                   常用工具类层
└───main.go                 主启动类
~~~
#### 🛖数据库表

| 表名                               | 描述                 |
| ---------------------------------- | -------------------- |
| ums_user                           | 用户信息表           |
| usm_receiver_address               | 用户收货地址         |
| usm_pcd_dic                        | 省市区字典表         |
| pms_spu                            | 商品spu表            |
| pms_spec_param                     | 商品规格key表        |
| pms_sku_pic                        | 商品sku图片表        |
| pms_sku                            | 商品sku表            |
| pms_product_detail_pic             | 商品详情图片表       |
| pms_product_category_brand_rel     | 商品分类和品牌关联表 |
| pms_product_category_attribute_rel | 商品分类和属性关联表 |
| pms_product_category               | 商品分类表           |
| pms_product_attribute_rel          | 商品和属性关联表     |
| pms_product_attribute              | 商品属性表           |
| pms_brand                          | 品牌表               |
| oms_pay_log                        | 支付记录表           |
| oms_order_item                     | 订单商品明细表       |
| oms_order                          | 订单表               |
| oms_cart_mess                      | 购物车信息表         |
| oms_cart                           | 购物车表             |




#### 🦉用户模块
* 注册:用户需要先获取验证码，然后校验用户手机号对应的用户是否存在，如果不存在就注册功能。
  1. 用户输入手机号，后端对手机号进行正则校验。如果不通过，返回错误信息。校验通过后，检查Redis中是否有该手机号的验证码缓存，如果有，返回信息，提示用户已经获取过验证码。以实现**获取验证码接口幂等性**。
  2. 如果Redis缓存中也没有，那么随机生成四位验证码，并将手机号和验证码缓存到Redis，设置5分钟的失效时间。并将手机号和验证码发送到RabbitMQ中，然后返回，本次请求结束。
  3. RabbitMQ消费端监听到队列中的消息后，调用阿里云SMS服务，发送短信。只有成功消费后，才会返回应答。然后这条消息才会从消息队列中删除。（启动MQ服务时，会有for循环不断的连接RabbitMQ，同时只有消费端成功消费一条消息后才会手动应答，同时处理消息也用for循环，直到处理方法返回true。这样就保证了即使连接异常关闭或者TCP连接丢失也会马上重建连接，因为没有手动应答，RabbitMQ也不会删除消息。同时也保证了一定处理成功这条数据）
  4. 用户输入验证码和密码后，首先校验密码强度，然后校验验证码是否和Redis中的一致，如果过期，本次注册失败。校验全部通过后，使用雪花算法生成全局唯一ID，同时使用MD5算法对密码加密，然后入库。
* 登录: 用户传递的密码加密后和数据库存储的密码比对
* 双Token刷新&鉴权: 用户登录后返回AccessToken和RefreshToken。前者短期存活，存储用户ID，用来鉴权。后者长期存活，不携带任何信息，用来辅助刷新前者
  	1. 返回给用户的Token不能长期存在，否则一旦用户的Token泄漏，那么意味着用户的账号将变得不安全。所以使用双Token刷新机制，每次用户登录成功后，返回短期存活的AccessToken，该Token中包含用户ID。RefreshToken只是存活时间更长，不包含任何信息。
   	2. 当用户访问个人资料、订单时，需要同时携带两个Token。我们使用中间件先校验AccessToken。
       * 如果Token不合法，用户鉴权未通过。
       * 如果Token未过期且合法，用户鉴权通过。
       * 如果Token已经过期且合法，判断RefershToken是否存在、合法、过期。如果都满足，则返回新的AccessToken。前端携带新的AccessToken重发本次请求。
   	3. 发起请求时都使用的HTTPS请求，所以泄漏Token的概率很低。
* 限制同一时间只有一台设备访问: 用户登录后缓存AccessToken到Redis中，访问特定接口时，需要比对Redis中、用户携带的AccessToken是否一致
* 用户个人信息/更新用户个人信息
* 更新头像: 接收图片后上传到阿里云OSS

#### 🦦商品模块
* 商品分类信息

  * 分类信息展示在商品主页，不会频繁变更。所以在数据库中查询完分类信息后，会缓存在Redis服务器中并设置失效时间。
  * 当更改商品分类信息会删除缓存或缓存过期，当用户再次成功获取分类信息后，缓存到Redis中。

* 商品属性

* 商品搜索功能

* 商品详情接口

  * 商品详情对象由三部分组成：商品的spu信息、商品的sku集合、商品的分类属性集合

  * 而后两个对象依赖于前一个spu信息，这里使用Golang中的协程+通道机制，同时查询两者并合并到商品详情对象中。单次查询商品详情所耗时降低一半。

    ![](https://richarli.oss-cn-beijing.aliyuncs.com/images/075316b978447d62027e0f41b3998d8.jpg)

