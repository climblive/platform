ALTER TABLE `color`
    ADD COLUMN `shared` TINYINT(1) NOT NULL AFTER `rgb_secondary`,
    DROP INDEX `name_UNIQUE`,
    ADD UNIQUE INDEX `name_UNIQUE` (`name` ASC, `organizer_id` ASC);