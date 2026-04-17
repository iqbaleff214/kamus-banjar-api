-- Migration 004: extend words table with source tracking, status, and contributor fields.
-- The original words table has: id, word, letter, verified, source (varchar), data, created_at, updated_at

-- 1. Change source from varchar to ENUM
ALTER TABLE words
    MODIFY COLUMN source ENUM('official','community') NOT NULL DEFAULT 'official';

-- 2. Add status column (replaces verified flag)
ALTER TABLE words
    ADD COLUMN status ENUM('active','pending','rejected') NOT NULL DEFAULT 'active' AFTER source;

-- 3. Add contributor tracking columns
ALTER TABLE words
    ADD COLUMN contributor_id CHAR(36) NULL AFTER status,
    ADD COLUMN approved_by   CHAR(36) NULL AFTER contributor_id,
    ADD COLUMN approved_at   DATETIME NULL AFTER approved_by;

-- 4. Add FK constraints (nullable — no cascade required here)
ALTER TABLE words
    ADD CONSTRAINT fk_words_contributor FOREIGN KEY (contributor_id) REFERENCES users(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_words_approver    FOREIGN KEY (approved_by)    REFERENCES users(id) ON DELETE SET NULL;

-- 5. Drop the old verified column (superseded by status)
ALTER TABLE words
    DROP COLUMN verified;

-- 6. Backfill: all existing rows are official + active
UPDATE words SET source = 'official', status = 'active';
