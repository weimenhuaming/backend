-- 黑名单表（存储 refresh token）
CREATE TABLE `token_blacklist` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `refresh_token` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_refresh_token` (`refresh_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

