CREATE TABLE `user`
(
    `id`          bigint(20)                                                   NOT NULL AUTO_INCREMENT COMMENT 'ID, 范围 -2^63 ~ 2^63-1',
    `user_id`     bigint(20)                                                   NOT NULL COMMENT '用户ID',
    `username`    varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '用户名',
    `password`    varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
    `phone`       char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci             default null COMMENT '手机号',
    `email`       varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '邮箱',
    `avatar`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci         DEFAULT NULL COMMENT '头像',
    `gender`      tinyint(4) UNSIGNED                                          NOT NULL DEFAULT '0' COMMENT '性别，0 -> 男, 1 -> 女, 2 -> 未知; 范围: -128~127',
    `create_time` timestamp                                                    NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp                                                    NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    unique key `idx_phone` (`phone`) USING BTREE,
    unique key `idx_email` (`email`) USING BTREE,
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;
