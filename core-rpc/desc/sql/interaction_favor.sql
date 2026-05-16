CREATE TABLE `interaction_favor` (
     `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
     `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
     `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
     `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
     `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
     `action_type` VARCHAR(16) NOT NULL COMMENT '操作类型：favor/unknown',
     PRIMARY KEY (`id`),
     UNIQUE KEY `uk_user_article_favor` (`user_id`, `article_id`, `action_type`) COMMENT '收藏唯一索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏互动表';