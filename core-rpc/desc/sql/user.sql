CREATE TABLE `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
    `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码',
    `phone` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号',
    `avatar` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像URL',
    `email` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '邮箱',
    `role` VARCHAR(32) NOT NULL DEFAULT 'user' COMMENT '权限：admin/user/guest',
    `sex` VARCHAR(8) NOT NULL DEFAULT 'unknown' COMMENT '性别：male/female/unknown',
    `age` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '年龄',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`) COMMENT '用户名唯一索引',
    KEY `idx_phone` (`phone`) COMMENT '手机号索引',
    KEY `idx_email` (`email`) COMMENT '邮箱索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';