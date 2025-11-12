ALTER TABLE `contest` ADD COLUMN `final_enabled` TINYINT(1) NOT NULL AFTER `location_id`;
ALTER TABLE `comp_class` ADD COLUMN `color` VARCHAR(7) NULL AFTER `description`;
ALTER TABLE `problem` ADD COLUMN `name` VARCHAR(64) NULL AFTER `color_id`;

UPDATE `contest` SET `final_enabled` = 1;