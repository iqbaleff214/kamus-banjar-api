-- Migration 005: contribution audit-trail table.
-- word_id is BIGINT UNSIGNED to match words.id (AUTO_INCREMENT).

CREATE TABLE IF NOT EXISTS contributions (
    id             CHAR(36)         NOT NULL PRIMARY KEY,
    word_id        BIGINT UNSIGNED  NOT NULL,
    contributor_id CHAR(36)         NOT NULL,
    reviewer_id    CHAR(36)         NULL,
    action         ENUM('submitted','approved','rejected','revised') NOT NULL,
    notes          TEXT             NULL,
    created_at     DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id)         REFERENCES words(id)  ON DELETE CASCADE,
    FOREIGN KEY (contributor_id)  REFERENCES users(id)  ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id)     REFERENCES users(id)  ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
