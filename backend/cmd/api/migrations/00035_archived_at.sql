-- +goose Up
ALTER TABLE contest ADD COLUMN `archived_at` DATETIME NULL DEFAULT NULL AFTER `archived`;
UPDATE contest SET archived_at = NOW() WHERE archived = TRUE;
ALTER TABLE contest DROP COLUMN `archived`;

-- +goose Down
ALTER TABLE contest ADD COLUMN `archived` TINYINT(1) NOT NULL DEFAULT 0 AFTER `organizer_id`;
UPDATE contest SET archived = TRUE WHERE archived_at IS NOT NULL;
ALTER TABLE contest DROP COLUMN `archived_at`;
