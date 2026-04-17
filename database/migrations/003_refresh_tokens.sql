-- Migration: refresh_tokens table for JWT rotation

CREATE TABLE IF NOT EXISTS `refresh_tokens` (
    `id`         CHAR(36)     NOT NULL,
    `user_id`    CHAR(36)     NOT NULL,
    `token_hash` VARCHAR(64)  NOT NULL,
    `expires_at` DATETIME     NOT NULL,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `refresh_tokens_hash_unique` (`token_hash`),
    KEY `refresh_tokens_user_id` (`user_id`),
    CONSTRAINT `fk_refresh_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
