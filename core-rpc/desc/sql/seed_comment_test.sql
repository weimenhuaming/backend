-- =============================================================================
-- 评论测试数据（执行前请先 USE 你的数据库）
--
-- 用法：
--   1. 打开下面「配置区」，把 @article_id / @user_* 改成你库里真实存在的 ID
--   2. 或直接执行「自动探测」段（有文章和用户时会自动取第一条）
--   3. 整文件执行一次即可
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 配置区（手动指定时取消注释并填写）
-- -----------------------------------------------------------------------------
-- SET @article_id = 1;
-- SET @user_a = 1;   -- 一级评论作者
-- SET @user_b = 2;   -- 回复者 A
-- SET @user_c = 3;   -- 回复者 B

-- -----------------------------------------------------------------------------
-- 自动探测（未手动配置时使用）
-- -----------------------------------------------------------------------------
SET @article_id = IFNULL(@article_id, (SELECT id FROM article WHERE deleted_at IS NULL ORDER BY id LIMIT 1));
SET @user_a = IFNULL(@user_a, (SELECT id FROM user WHERE deleted_at IS NULL ORDER BY id LIMIT 1 OFFSET 0));
SET @user_b = IFNULL(@user_b, (SELECT id FROM user WHERE deleted_at IS NULL ORDER BY id LIMIT 1 OFFSET 1));
SET @user_c = IFNULL(@user_c, (SELECT id FROM user WHERE deleted_at IS NULL ORDER BY id LIMIT 1 OFFSET 2));

-- 若只有一个用户，复用同一 user_id
SET @user_b = IFNULL(@user_b, @user_a);
SET @user_c = IFNULL(@user_c, @user_a);

SELECT @article_id AS article_id, @user_a AS user_a, @user_b AS user_b, @user_c AS user_c;

-- 没有文章或用户时直接停止（避免脏数据）
-- 若下面 INSERT 报错，请先插入 article / user 再执行本脚本

START TRANSACTION;

-- =============================================================================
-- 文章 1：3 条一级评论 + 若干回复
-- =============================================================================

-- 一级评论 1（热门）
INSERT INTO `comment` (`article_id`, `user_id`, `parent_id`, `root_id`, `reply_to_id`, `reply_to_name`, `content`, `like_count`, `child_count`, `created_at`)
VALUES (@article_id, @user_a, 0, 0, 0, '', '写得很清楚，尤其是中间那段原理说明，收藏了。', 28, 0, DATE_SUB(NOW(), INTERVAL 2 DAY));
SET @c1 = LAST_INSERT_ID();
UPDATE `comment` SET `root_id` = @c1 WHERE `id` = @c1;

-- 一级评论 2
INSERT INTO `comment` (`article_id`, `user_id`, `parent_id`, `root_id`, `reply_to_id`, `reply_to_name`, `content`, `like_count`, `child_count`, `created_at`)
VALUES (@article_id, @user_b, 0, 0, 0, '', '有个小问题：示例代码里变量命名是不是写反了？', 5, 0, DATE_SUB(NOW(), INTERVAL 36 HOUR));
SET @c2 = LAST_INSERT_ID();
UPDATE `comment` SET `root_id` = @c2 WHERE `id` = @c2;

-- 一级评论 3（最新）
INSERT INTO `comment` (`article_id`, `user_id`, `parent_id`, `root_id`, `reply_to_id`, `reply_to_name`, `content`, `like_count`, `child_count`, `created_at`)
VALUES (@article_id, @user_c, 0, 0, 0, '', '期待作者出续篇，讲一下生产环境怎么部署。', 2, 0, DATE_SUB(NOW(), INTERVAL 3 HOUR));
SET @c3 = LAST_INSERT_ID();
UPDATE `comment` SET `root_id` = @c3 WHERE `id` = @c3;

-- 回复：挂在 @c1 下（预览列表会显示前 3 条）
INSERT INTO `comment` (`article_id`, `user_id`, `parent_id`, `root_id`, `reply_to_id`, `reply_to_name`, `content`, `like_count`, `child_count`, `created_at`)
VALUES
(@article_id, @user_b, @c1, @c1, @user_a, (SELECT name FROM user WHERE id = @user_a LIMIT 1), '同感，我也照着做了一遍，成功了。', 3, 0, DATE_SUB(NOW(), INTERVAL 47 HOUR)),
(@article_id, @user_c, @c1, @c1, @user_b, (SELECT name FROM user WHERE id = @user_b LIMIT 1), '感谢分享，省了我不少排查时间。', 1, 0, DATE_SUB(NOW(), INTERVAL 40 HOUR)),
(@article_id, @user_a, @c1, @c1, @user_c, (SELECT name FROM user WHERE id = @user_c LIMIT 1), '不客气，有问题随时留言。', 0, 0, DATE_SUB(NOW(), INTERVAL 30 HOUR)),
(@article_id, @user_b, @c1, @c1, @user_a, (SELECT name FROM user WHERE id = @user_a LIMIT 1), '第四楼：测试「查看更多回复」分页用。', 0, 0, DATE_SUB(NOW(), INTERVAL 20 HOUR));

-- 回复：挂在 @c2 下
INSERT INTO `comment` (`article_id`, `user_id`, `parent_id`, `root_id`, `reply_to_id`, `reply_to_name`, `content`, `like_count`, `child_count`, `created_at`)
VALUES
(@article_id, @user_a, @c2, @c2, @user_b, (SELECT name FROM user WHERE id = @user_b LIMIT 1), '确实，我已在最新版里修正，感谢指出。', 4, 0, DATE_SUB(NOW(), INTERVAL 24 HOUR)),
(@article_id, @user_c, @c2, @c2, @user_a, (SELECT name FROM user WHERE id = @user_a LIMIT 1), '修了之后跑通了，赞。', 0, 0, DATE_SUB(NOW(), INTERVAL 18 HOUR));

-- 回复：挂在 @c3 下（仅 1 条）
INSERT INTO `comment` (`article_id`, `user_id`, `parent_id`, `root_id`, `reply_to_id`, `reply_to_name`, `content`, `like_count`, `child_count`, `created_at`)
VALUES (@article_id, @user_a, @c3, @c3, @user_c, (SELECT name FROM user WHERE id = @user_c LIMIT 1), '续篇在写了，大概下周发。', 0, 0, DATE_SUB(NOW(), INTERVAL 1 HOUR));

-- 更新一级评论 child_count（MySQL 同表子查询需包一层）
SET @cnt_c1 = (SELECT COUNT(*) FROM (SELECT id FROM `comment` WHERE root_id = @c1 AND parent_id > 0 AND deleted_at IS NULL) AS t);
SET @cnt_c2 = (SELECT COUNT(*) FROM (SELECT id FROM `comment` WHERE root_id = @c2 AND parent_id > 0 AND deleted_at IS NULL) AS t);
SET @cnt_c3 = (SELECT COUNT(*) FROM (SELECT id FROM `comment` WHERE root_id = @c3 AND parent_id > 0 AND deleted_at IS NULL) AS t);
UPDATE `comment` SET `child_count` = @cnt_c1 WHERE `id` = @c1;
UPDATE `comment` SET `child_count` = @cnt_c2 WHERE `id` = @c2;
UPDATE `comment` SET `child_count` = @cnt_c3 WHERE `id` = @c3;

-- 软删除一条回复（测试列表应过滤）
INSERT INTO `comment` (`article_id`, `user_id`, `parent_id`, `root_id`, `reply_to_id`, `reply_to_name`, `content`, `like_count`, `child_count`, `created_at`, `deleted_at`)
VALUES (@article_id, @user_b, @c1, @c1, @user_a, (SELECT name FROM user WHERE id = @user_a LIMIT 1), '这条已软删除，接口不应返回', 0, 0, DATE_SUB(NOW(), INTERVAL 10 HOUR), NOW());

-- 同步文章评论总数（仅统计未删除）
UPDATE `article` a
SET a.comment_count = (
    SELECT COUNT(*) FROM `comment` c
    WHERE c.article_id = a.id AND c.deleted_at IS NULL
)
WHERE a.id = @article_id;

COMMIT;

-- =============================================================================
-- 验证查询（可选）
-- =============================================================================
-- 一级评论（按时间）
-- SELECT id, user_id, content, like_count, child_count, created_at
-- FROM comment
-- WHERE article_id = @article_id AND parent_id = 0 AND deleted_at IS NULL
-- ORDER BY created_at DESC;

-- 某根评论下的回复
-- SELECT id, parent_id, root_id, reply_to_name, content, created_at
-- FROM comment
-- WHERE root_id = @c1 AND parent_id > 0 AND deleted_at IS NULL
-- ORDER BY created_at ASC;
