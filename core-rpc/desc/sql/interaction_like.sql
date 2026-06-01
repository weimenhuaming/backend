-- 点赞这一块内容目前是通过 userid - object_id - object_type 来做唯一索引的，
-- object_type 目前针对的对象有文章，评论这两个方面

CREATE TABLE `interaction_like` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `object_type` VARCHAR(16) NOT NULL COMMENT '对象类型：article/comment',
    `object_id` BIGINT UNSIGNED NOT NULL COMMENT '对象ID',
    `action_type` VARCHAR(16) NOT NULL COMMENT '操作类型：like/unknown',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_object_like` (`user_id`, `object_type`, `object_id`) COMMENT '用户对同一对象仅一条互动记录'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞互动表';