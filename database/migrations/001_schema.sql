-- Migration: initial schema
-- Creates the letters and words tables for Kamus Banjar.

CREATE TABLE IF NOT EXISTS `letters` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `letter`     char(1) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `letters_letter_unique` (`letter`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `words` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `word`       varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `letter`     char(1) COLLATE utf8mb4_unicode_ci NOT NULL,
  `verified`   tinyint(1) NOT NULL DEFAULT '0',
  `source`     varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'official',
  `data`       json DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `words_word_unique` (`word`),
  KEY `words_letter_foreign` (`letter`),
  CONSTRAINT `words_letter_foreign` FOREIGN KEY (`letter`) REFERENCES `letters` (`letter`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
