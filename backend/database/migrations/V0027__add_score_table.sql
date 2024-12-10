CREATE TABLE IF NOT EXISTS `score` (
  `contender_id` INT NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  `score` INT NOT NULL,
  `placement` INT NOT NULL,
  `finalist` TINYINT(1) NOT NULL,
  `rank_order` INT NOT NULL,
  PRIMARY KEY (`contender_id`),
  CONSTRAINT `fk_score_1`
    FOREIGN KEY (`contender_id`)
    REFERENCES `contender` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;