CREATE TABLE `interaction_comment` (
       `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
       `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
       `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
       `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
       `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
       `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
       `content` TEXT NOT NULL COMMENT '评论内容',
       PRIMARY KEY (`id`),
       KEY `idx_article_id` (`article_id`) COMMENT '文章索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论互动表';