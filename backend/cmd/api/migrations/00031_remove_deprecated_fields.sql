-- +goose Up
ALTER TABLE `contest` DROP COLUMN `protected`;
ALTER TABLE `contest` DROP COLUMN `final_enabled`;

ALTER TABLE `contender` DROP COLUMN `club`;

UPDATE `problem`
    SET `description` = `name`
    WHERE `description` IS NULL OR `description` = '';

ALTER TABLE `problem` DROP COLUMN `name`;

ALTER TABLE `organizer` DROP COLUMN `homepage`;

ALTER TABLE `user` DROP COLUMN `name`;

ALTER TABLE `raffle` DROP COLUMN `active`;

-- +goose Down
ALTER TABLE `raffle` ADD COLUMN `active` TINYINT(1) NOT NULL DEFAULT 0 AFTER `contest_id`;

ALTER TABLE `user` ADD COLUMN `name` VARCHAR(32) NOT NULL AFTER `id`;

ALTER TABLE `organizer` ADD COLUMN `homepage` VARCHAR(255) NULL AFTER `name`;

ALTER TABLE `problem` ADD COLUMN `name` VARCHAR(64) NULL AFTER `hold_color_secondary`;

ALTER TABLE `contender` ADD COLUMN `club` VARCHAR(128) NULL DEFAULT NULL AFTER `name`;

ALTER TABLE `contest`
  ADD COLUMN `protected` TINYINT(1) NOT NULL DEFAULT 0 AFTER `archived`,
  ADD COLUMN `final_enabled` TINYINT(1) NOT NULL DEFAULT 1 AFTER `location`;
