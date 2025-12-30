-- +goose Up
ALTER TABLE `contest` ADD COLUMN `created` TIMESTAMP NOT NULL DEFAULT '1970-01-01 00:00:00' AFTER `grace_period`;

-- +goose Down
ALTER TABLE `contest` DROP COLUMN `created`;