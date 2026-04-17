-- Migration: users and refresh_tokens tables

CREATE TABLE IF NOT EXISTS `users` (
    `id`         CHAR(36)              NOT NULL,
    `name`       VARCHAR(100)          NOT NULL,
    `email`      VARCHAR(255)          NOT NULL,
    `password`   VARCHAR(255)          NOT NULL,
    `role`       ENUM('admin','user')  NOT NULL DEFAULT 'user',
    `is_active`  TINYINT(1)            NOT NULL DEFAULT 1,
    `created_at` DATETIME              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME              NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_email_unique` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
