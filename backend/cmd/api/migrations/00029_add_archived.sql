-- +goose Up
ALTER TABLE contest ADD COLUMN `archived` TINYINT(1) NOT NULL DEFAULT 0 AFTER `organizer_id`;

-- +goose Down
ALTER TABLE contest DROP COLUMN `archived`;