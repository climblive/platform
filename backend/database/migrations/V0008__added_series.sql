CREATE TABLE `series` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `organizer_id` INT NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_series_1`
    FOREIGN KEY (`organizer_id`)
    REFERENCES `organizer` (`id`)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT)
ENGINE = InnoDB;

ALTER TABLE `contest`
    ADD COLUMN `series_id` INT NULL AFTER `organizer_id`,
    ADD CONSTRAINT `fk_contest_3`
        FOREIGN KEY (`series_id`)
        REFERENCES `series` (`id`)
        ON DELETE RESTRICT
        ON UPDATE RESTRICT;