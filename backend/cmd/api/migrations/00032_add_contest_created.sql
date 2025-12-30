-- +goose Up
ALTER TABLE `contest` ADD COLUMN `created` TIMESTAMP NOT NULL DEFAULT 0 AFTER `grace_period`;

-- +goose Down
ALTER TABLE `contest` DROP COLUMN `created`;