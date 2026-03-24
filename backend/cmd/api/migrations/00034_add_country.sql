-- +goose Up
ALTER TABLE `contest` ADD COLUMN `country` VARCHAR(2) NOT NULL DEFAULT 'AQ' AFTER `location`;

-- +goose Down
ALTER TABLE `contest` DROP COLUMN `country`;
