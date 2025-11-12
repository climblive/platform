ALTER TABLE `contest` ADD COLUMN `protected` TINYINT(1) NOT NULL DEFAULT 0 AFTER `organizer_id`;
