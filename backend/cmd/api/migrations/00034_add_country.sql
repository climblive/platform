-- +goose Up
ALTER TABLE `contest` ADD COLUMN `country` VARCHAR(2) NULL DEFAULT 'SE';

-- +goose Down
ALTER TABLE `contest` DROP COLUMN `country`;
