-- +goose Up
CREATE TABLE IF NOT EXISTS `organizer_invite` (
  `id` VARCHAR(36) NOT NULL,
  `organizer_id` INT NOT NULL,
  `expires_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_organizer_invite_1`
    FOREIGN KEY (`organizer_id`)
    REFERENCES `organizer` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;

CREATE INDEX `fk_organizer_invite_1_idx` ON `organizer_invite` (`organizer_id` ASC);

-- +goose Down
DROP TABLE `organizer_invite`;