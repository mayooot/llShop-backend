/*
 Navicat Premium Data Transfer

 Source Server         : mysql8.0
 Source Server Type    : MySQL
 Source Server Version : 80019
 Source Host           : localhost:3306
 Source Schema         : shop

 Target Server Type    : MySQL
 Target Server Version : 80019
 File Encoding         : 65001

 Date: 17/10/2022 18:16:51
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for oms_cart
-- ----------------------------
DROP TABLE IF EXISTS `oms_cart`;
CREATE TABLE `oms_cart`  (
                             `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                             `userId` bigint NOT NULL COMMENT '用户ID\r\n',
                             `skuId` bigint NOT NULL COMMENT '商品ID\r\n',
                             `num` int NOT NULL COMMENT '购买数量',
                             PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for oms_cart_mess
-- ----------------------------
DROP TABLE IF EXISTS `oms_cart_mess`;
CREATE TABLE `oms_cart_mess`  (
                                  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                  `userId` bigint NOT NULL COMMENT '用户ID\r\n',
                                  `skuId` bigint NOT NULL COMMENT '商品ID\r\n',
                                  `status` int NOT NULL COMMENT '消息状态，0=>已经发出但mq没有确认，1=>生产者已经收到确认',
                                  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for oms_order
-- ----------------------------
DROP TABLE IF EXISTS `oms_order`;
CREATE TABLE `oms_order`  (
                              `id` bigint NOT NULL AUTO_INCREMENT,
                              `user_id` bigint NULL DEFAULT NULL COMMENT '用户ID(对应用户表主键ID)',
                              `order_no` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '订单编号',
                              `total_money` decimal(10, 2) NULL DEFAULT NULL COMMENT '订单总金额合计',
                              `pay_money` decimal(10, 2) NULL DEFAULT NULL COMMENT '实付金额合计',
                              `total_num` int UNSIGNED NULL DEFAULT NULL COMMENT '数量合计',
                              `pay_type` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '支付方式：0->在线支付；1->货到付款',
                              `order_status` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '订单状态：0->待付款；1->待发货；2->已发货；3->已完成；4->已关闭；5->超时',
                              `pay_status` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '支付状态：0->未支付；1->支付成功；2->支付失败',
                              `pay_time` datetime NULL DEFAULT NULL COMMENT '支付时间',
                              `receiver_name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '收件人名称',
                              `receiver_phone` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '收件人电话',
                              `receiver_address` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '收件人地址',
                              `expiration_time` datetime NULL DEFAULT NULL COMMENT '订单过期时间',
                              `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                              PRIMARY KEY (`id`) USING BTREE,
                              INDEX `order_no`(`order_no`) USING BTREE,
                              INDEX `user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '订单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for oms_order_item
-- ----------------------------
DROP TABLE IF EXISTS `oms_order_item`;
CREATE TABLE `oms_order_item`  (
                                   `id` bigint NOT NULL AUTO_INCREMENT,
                                   `order_id` bigint NULL DEFAULT NULL COMMENT '订单ID(对应订单表主键ID)',
                                   `spu_id` bigint NULL DEFAULT NULL COMMENT '商品spuID(对应商品spu表主键ID)',
                                   `product_pic` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '商品图片',
                                   `product_name` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '商品名称',
                                   `product_price` decimal(10, 2) NULL DEFAULT NULL COMMENT '销售价格',
                                   `product_quantity` int NULL DEFAULT NULL COMMENT '购买数量',
                                   `sku_id` bigint NULL DEFAULT NULL COMMENT '商品skuID(对应商品sku表主键ID)',
                                   `product_attr` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '商品销售属性:[{\"key\":\"颜色\",\"value\":\"颜色\"},{\"key\":\"容量\",\"value\":\"4G\"}]',
                                   `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                   PRIMARY KEY (`id`) USING BTREE,
                                   INDEX `order_id`(`order_id`) USING BTREE,
                                   INDEX `spu_id`(`spu_id`) USING BTREE,
                                   INDEX `sku_id`(`sku_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 322 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '订单商品明细表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for oms_pay_log
-- ----------------------------
DROP TABLE IF EXISTS `oms_pay_log`;
CREATE TABLE `oms_pay_log`  (
                                `id` bigint NOT NULL AUTO_INCREMENT,
                                `user_id` bigint NULL DEFAULT NULL COMMENT '用户ID(对应用户表主键ID)',
                                `order_id` bigint NULL DEFAULT NULL COMMENT '订单ID(对应订单表主键ID)',
                                `order_no` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '订单编号',
                                `pay_way` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '支付方式：1->支付宝支付；2->微信支付',
                                `pay_trade_no` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '第三方支付订单交易号',
                                `pay_status` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '支付状态：1->支付成功；2->支付失败',
                                `pay_amount` decimal(10, 2) NULL DEFAULT NULL COMMENT '支付金额',
                                `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 366 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '支付记录表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_brand
-- ----------------------------
DROP TABLE IF EXISTS `pms_brand`;
CREATE TABLE `pms_brand`  (
                              `id` bigint NOT NULL AUTO_INCREMENT,
                              `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '品牌名称',
                              `first_letter` char(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '品牌的首字母',
                              `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                              PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '品牌表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_product_attribute
-- ----------------------------
DROP TABLE IF EXISTS `pms_product_attribute`;
CREATE TABLE `pms_product_attribute`  (
                                          `id` bigint NOT NULL AUTO_INCREMENT,
                                          `type` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '类型：0->属性key；1->属性value',
                                          `parent_id` bigint NULL DEFAULT NULL COMMENT '\"属性key\"ID：0->属性key',
                                          `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '属性名称',
                                          `sort` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '排序',
                                          `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                          `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                          PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 45 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品属性表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_product_attribute_rel
-- ----------------------------
DROP TABLE IF EXISTS `pms_product_attribute_rel`;
CREATE TABLE `pms_product_attribute_rel`  (
                                              `id` bigint NOT NULL AUTO_INCREMENT,
                                              `spu_id` bigint NULL DEFAULT NULL COMMENT '商品spuID(对应商品spu表主键ID)',
                                              `product_attribute_id` bigint NULL DEFAULT NULL COMMENT '商品属性ID(对应商品属性表主键ID)',
                                              `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                              `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                              PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品和属性关联表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_product_category
-- ----------------------------
DROP TABLE IF EXISTS `pms_product_category`;
CREATE TABLE `pms_product_category`  (
                                         `id` bigint NOT NULL AUTO_INCREMENT,
                                         `parent_id` bigint NULL DEFAULT NULL COMMENT '上级分类的编号：0->一级分类',
                                         `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '分类名称',
                                         `abbreviation` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '分类名称简称',
                                         `level` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '分类级别：0->1级；1->2级',
                                         `show_status` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '显示状态：0->不显示；1->显示',
                                         `icon` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '图标',
                                         `sort` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '排序',
                                         `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                         `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                         PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 46 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品分类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_product_category_attribute_rel
-- ----------------------------
DROP TABLE IF EXISTS `pms_product_category_attribute_rel`;
CREATE TABLE `pms_product_category_attribute_rel`  (
                                                       `id` bigint NOT NULL AUTO_INCREMENT,
                                                       `product_category_id` bigint NULL DEFAULT NULL COMMENT '商品分类ID(对应商品分类表主键ID)',
                                                       `product_attribute_id` bigint NULL DEFAULT NULL COMMENT '商品属性ID(对应商品属性表主键ID)',
                                                       `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                                       `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                                       PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品分类和属性关联表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_product_category_brand_rel
-- ----------------------------
DROP TABLE IF EXISTS `pms_product_category_brand_rel`;
CREATE TABLE `pms_product_category_brand_rel`  (
                                                   `id` bigint NOT NULL AUTO_INCREMENT,
                                                   `product_category_id` bigint NULL DEFAULT NULL COMMENT '商品分类ID(对应商品分类表主键ID)',
                                                   `brand_id` bigint NULL DEFAULT NULL COMMENT '商品品牌ID(对应商品品牌表主键ID)',
                                                   `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                                   `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                                   PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品分类和品牌关联表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_product_detail_pic
-- ----------------------------
DROP TABLE IF EXISTS `pms_product_detail_pic`;
CREATE TABLE `pms_product_detail_pic`  (
                                           `id` bigint NOT NULL AUTO_INCREMENT,
                                           `spu_id` bigint NULL DEFAULT NULL COMMENT '商品spuID(对应商品spu表主键ID)',
                                           `pic_url` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '图片URL',
                                           `sort` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '排序',
                                           `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                           `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                           PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 115 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品详情图片表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_sku
-- ----------------------------
DROP TABLE IF EXISTS `pms_sku`;
CREATE TABLE `pms_sku`  (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `spu_id` bigint NULL DEFAULT NULL COMMENT '商品spuID(对应商品spu表主键ID)',
                            `title` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '商品标题',
                            `price` decimal(10, 2) NULL DEFAULT NULL COMMENT '价格',
                            `unit` varchar(16) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '商品单位',
                            `stock` int NULL DEFAULT 0 COMMENT '库存',
                            `sale` int NULL DEFAULT NULL COMMENT '销量',
                            `indexes` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT 'spu中商品规格的对应下标组合',
                            `product_sku_specification` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '商品sku规格(json格式，反序列化时请使用linkedHashMap，保证有序)',
                            `is_default` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '默认规格：0->不是；1->是',
                            `valid` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否有效，0->无效；1->有效',
                            `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1000026 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品sku表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_sku_pic
-- ----------------------------
DROP TABLE IF EXISTS `pms_sku_pic`;
CREATE TABLE `pms_sku_pic`  (
                                `id` bigint NOT NULL AUTO_INCREMENT,
                                `sku_id` bigint NULL DEFAULT NULL COMMENT '商品skuID(对应商品sku表主键ID)',
                                `pic_url` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '图片URL',
                                `is_default` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '默认展示：0->否；1->是',
                                `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 86 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品sku图片表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_spec_param
-- ----------------------------
DROP TABLE IF EXISTS `pms_spec_param`;
CREATE TABLE `pms_spec_param`  (
                                   `id` bigint NOT NULL AUTO_INCREMENT,
                                   `product_category_id` bigint NULL DEFAULT NULL COMMENT '商品分类ID(对应商品分类表主键ID)',
                                   `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '规格key名称',
                                   `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                   PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品规格key表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for pms_spu
-- ----------------------------
DROP TABLE IF EXISTS `pms_spu`;
CREATE TABLE `pms_spu`  (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `brand_id` bigint NULL DEFAULT NULL COMMENT '品牌ID(对应品牌表主键ID)',
                            `cid1` bigint NULL DEFAULT NULL COMMENT '一级分类ID(对应商品分类表主键ID)',
                            `cid2` bigint NULL DEFAULT NULL COMMENT '二级分类ID(对应商品分类表主键ID)',
                            `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '商品名称',
                            `sub_title` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '副标题',
                            `sale` int NULL DEFAULT NULL COMMENT '商品总销量',
                            `comment_total_score` int NULL DEFAULT NULL COMMENT '评价总评分',
                            `comment_amount` int NULL DEFAULT NULL COMMENT '商品评价数量',
                            `product_specification` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '商品规格(json格式，用于商品详情页展示商品所有规格)',
                            `default_price` decimal(10, 2) NULL DEFAULT NULL COMMENT '商品默认价格',
                            `default_pic_url` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '商品默认图片URL',
                            `publish_status` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '上架状态：0->下架；1->上架',
                            `verify_status` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '审核状态：0->未审核；1->审核通过',
                            `valid` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否有效，0->已删除；1->有效',
                            `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '商品spu表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ums_pcd_dic
-- ----------------------------
DROP TABLE IF EXISTS `ums_pcd_dic`;
CREATE TABLE `ums_pcd_dic`  (
                                `id` int NOT NULL,
                                `name` varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '名称',
                                `parent_id` int NULL DEFAULT NULL COMMENT '上级编号：parentId=id -> 省/直辖市',
                                PRIMARY KEY (`id`) USING BTREE,
                                INDEX `FK_CHINA_REFERENCE_CHINA`(`parent_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '省市区字典表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ums_receiver_address
-- ----------------------------
DROP TABLE IF EXISTS `ums_receiver_address`;
CREATE TABLE `ums_receiver_address`  (
                                         `id` bigint NOT NULL AUTO_INCREMENT,
                                         `user_id` bigint NULL DEFAULT NULL COMMENT '用户ID(对应用户表主键ID)',
                                         `user_name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '收件人名称',
                                         `phone_number` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '收件人电话',
                                         `default_status` tinyint UNSIGNED NULL DEFAULT NULL COMMENT '是否为默认收货地址：0->不是；1->是',
                                         `county_id` bigint NULL DEFAULT NULL COMMENT '区县id',
                                         `detail_address` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '详细地址(街道)',
                                         `created_time` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                         `updated_time` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                         `is_delete` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除（0：否，1：是）',
                                         PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 35 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户收货地址表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ums_user
-- ----------------------------
DROP TABLE IF EXISTS `ums_user`;
CREATE TABLE `ums_user`  (
                             `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID, 范围 -2^63 ~ 2^63-1',
                             `user_id` bigint NOT NULL COMMENT '用户ID',
                             `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'user' COMMENT '用户名',
                             `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
                             `phone` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '手机号',
                             `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '邮箱',
                             `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'https://pic1.zhimg.com/v2-abed1a8c04700ba7d72b45195223e0ff_is.jpeg' COMMENT '头像',
                             `gender` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '性别，0 -> 男, 1 -> 女, 2 -> 未知; 范围: -128~127',
                             `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                             `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             PRIMARY KEY (`id`) USING BTREE,
                             UNIQUE INDEX `idx_user_id`(`user_id`) USING BTREE,
                             UNIQUE INDEX `idx_phone`(`phone`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 27 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
