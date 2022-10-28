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
* 注册: 使用手机号注册，后端生成4位验证码，使用MQ异步调动阿里云SMS发送短信。注册成功后，密码使用MD5加密后入库
* 登录: 用户传递的密码加密后和数据库存储的密码比对
* 双Token刷新&鉴权: 用户登录后返回AccessToken和RefreshToken。前者短期存活，存储用户ID，用来鉴权。后者长期存活，不携带任何信息，用来辅助刷新前者
* 限制同一时间只有一台设备访问: 用户登录后缓存AccessToken到Redis中，访问特定接口时，需要比对Redis中、用户携带的AccessToken是否一致
* 用户个人信息/更新用户个人信息
* 更新头像: 接收图片后上传到阿里云OSS

#### 🦦商品模块
* 商品分类信息
* 商品属性
* 商品搜索功能
* 商品详情接口