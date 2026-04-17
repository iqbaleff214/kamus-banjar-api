-- Migration 006: community features — votes, comments, bookmarks, word-of-the-day.
-- All tables reference word by text (VARCHAR) to avoid BIGINT/UUID type conflicts
-- with the existing words.id (AUTO_INCREMENT BIGINT).

CREATE TABLE IF NOT EXISTS word_votes (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    word       VARCHAR(255) NOT NULL,
    user_id    CHAR(36)     NOT NULL,
    vote       TINYINT      NOT NULL COMMENT '1 = upvote, -1 = downvote',
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_vote (word, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS word_comments (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    word       VARCHAR(255) NOT NULL,
    user_id    CHAR(36)     NOT NULL,
    parent_id  CHAR(36)     NULL,
    body       TEXT         NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)   REFERENCES users(id)         ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES word_comments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS word_bookmarks (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    word       VARCHAR(255) NOT NULL,
    user_id    CHAR(36)     NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_bookmark (word, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS word_of_the_day (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    word       VARCHAR(255) NOT NULL,
    date       DATE         NOT NULL UNIQUE,
    set_by     CHAR(36)     NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
