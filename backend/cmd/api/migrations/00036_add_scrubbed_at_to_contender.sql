-- +goose Up
ALTER TABLE `contender` ADD COLUMN `scrubbed_at` TIMESTAMP NULL DEFAULT NULL;

-- +goose Down
ALTER TABLE `contender` DROP COLUMN `scrubbed_at`;
