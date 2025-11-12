ALTER TABLE `location`
    ADD COLUMN `organizer_id` INT NOT NULL AFTER `id`;

UPDATE `location` SET `organizer_id` = (SELECT MIN(id) FROM `organizer`);

ALTER TABLE `location`
    ADD CONSTRAINT `fk_location_1`
        FOREIGN KEY (`organizer_id`)
        REFERENCES `organizer` (`id`)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT;

ALTER TABLE `color`
    ADD COLUMN `organizer_id` INT NOT NULL AFTER `id`;

UPDATE `color` SET `organizer_id` = (SELECT MIN(id) FROM `organizer`);

ALTER TABLE `color`
    ADD CONSTRAINT `fk_color_1`
        FOREIGN KEY (`organizer_id`)
        REFERENCES `organizer` (`id`)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT;