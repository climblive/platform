-- +goose Up
ALTER TABLE `problem` ADD COLUMN `zone_1_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `name`;
ALTER TABLE `problem` ADD COLUMN `zone_2_enabled` TINYINT(1) NOT NULL DEFAULT 0 AFTER `zone_1_enabled`;
ALTER TABLE `problem` ADD COLUMN `points_zone_1` INT NULL AFTER `description`;
ALTER TABLE `problem` ADD COLUMN `points_zone_2` INT NULL AFTER `points_zone_1`;
ALTER TABLE `problem` CHANGE `points` `points_top` INT NOT NULL AFTER `points_zone_2`;

ALTER TABLE `tick` ADD COLUMN `zone_1` TINYINT(1) NOT NULL DEFAULT 0 AFTER `timestamp`;
ALTER TABLE `tick` ADD COLUMN `attempts_zone_1` INT NOT NULL DEFAULT 0 AFTER `zone_1`;
ALTER TABLE `tick` ADD COLUMN `zone_2` TINYINT(1) NOT NULL DEFAULT 0 AFTER `attempts_zone_1`;
ALTER TABLE `tick` ADD COLUMN `attempts_zone_2` INT NOT NULL DEFAULT 0 AFTER `zone_2`;
ALTER TABLE `tick` ADD COLUMN `top` TINYINT(1) NOT NULL DEFAULT 0  AFTER `attempts_zone_2`;
ALTER TABLE `tick` ADD COLUMN `attempts_top` INT NOT NULL DEFAULT 0 AFTER `top`;

UPDATE `tick`
    SET
        `zone_1` = TRUE,
        `zone_2` = TRUE,
        `top` = TRUE,
        `attempts_zone_1` = CASE WHEN `flash` = 1 THEN 1 ELSE 999 END,
        `attempts_zone_2` = CASE WHEN `flash` = 1 THEN 1 ELSE 999 END,
        `attempts_top` = CASE WHEN `flash` = 1 THEN 1 ELSE 999 END;

ALTER TABLE `tick` DROP COLUMN `flash`;

-- +goose Down
ALTER TABLE `problem` DROP COLUMN `zone_1_enabled`;
ALTER TABLE `problem` DROP COLUMN `zone_2_enabled`;
ALTER TABLE `problem` DROP COLUMN `points_zone_1`;
ALTER TABLE `problem` DROP COLUMN `points_zone_2`;
ALTER TABLE `problem` CHANGE `points_top` `points` INT NOT NULL AFTER `description`;

ALTER TABLE `tick` ADD COLUMN `flash` TINYINT(1) NOT NULL DEFAULT 0 AFTER `problem_id`;

UPDATE `tick`
    SET `flash` = CASE WHEN top = 1 AND attempts_top = 1 THEN 1 ELSE 0 END;

ALTER TABLE `tick` DROP COLUMN `top`;
ALTER TABLE `tick` DROP COLUMN `attempts_top`;
ALTER TABLE `tick` DROP COLUMN `zone_1`;
ALTER TABLE `tick` DROP COLUMN `attempts_zone_1`;
ALTER TABLE `tick` DROP COLUMN `zone_2`;
ALTER TABLE `tick` DROP COLUMN `attempts_zone_2`;