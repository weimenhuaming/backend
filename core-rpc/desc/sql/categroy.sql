-- 创建分类数据表
CREATE TABLE `category`
(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
    `name` VARCHAR(128) NOT NULL COMMENT '分类名称',
    `description` TEXT NULL COMMENT '分类简介描述',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间，自动赋值',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间，修改自动刷新',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间，为空代表正常数据',
    PRIMARY KEY (`id`) COMMENT '主键索引',
    UNIQUE KEY `uk_category_name` (`name`) COMMENT '分类名称唯一不可重复'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分类表';