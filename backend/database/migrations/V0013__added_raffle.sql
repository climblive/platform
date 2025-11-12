CREATE TABLE `raffle` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `contest_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_raffle_1`
    FOREIGN KEY (`contest_id`)
    REFERENCES `contest` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

CREATE INDEX `fk_raffle_1_idx` ON `raffle` (`contest_id` ASC);

CREATE TABLE `raffle_winner` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `raffle_id` INT NOT NULL,
  `contender_id` INT NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_raffle_winner_1`
    FOREIGN KEY (`raffle_id`)
    REFERENCES `raffle` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_raffle_winner_2`
    FOREIGN KEY (`contender_id`)
    REFERENCES `contender` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

CREATE INDEX `fk_raffle_winner_1_idx` ON `raffle_winner` (`raffle_id` ASC);

CREATE INDEX `fk_raffle_winner_2_idx` ON `raffle_winner` (`contender_id` ASC);

CREATE UNIQUE INDEX `index4` ON `raffle_winner` (`raffle_id` ASC, `contender_id` ASC);