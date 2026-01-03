-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `unlock_request` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `contest_id` INT NOT NULL,
  `organizer_id` INT NOT NULL,
  `requested_by_user_id` INT NOT NULL,
  `status` ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
  `reviewed_by_user_id` INT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reviewed_at` TIMESTAMP NULL DEFAULT NULL,
  `reason` TEXT NULL,
  `review_note` TEXT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_unlock_request_contest`
    FOREIGN KEY (`contest_id`)
    REFERENCES `contest` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_unlock_request_organizer`
    FOREIGN KEY (`organizer_id`)
    REFERENCES `organizer` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_unlock_request_user`
    FOREIGN KEY (`requested_by_user_id`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_unlock_request_reviewer`
    FOREIGN KEY (`reviewed_by_user_id`)
    REFERENCES `user` (`id`)
    ON DELETE SET NULL
    ON UPDATE RESTRICT
) ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;

CREATE INDEX `fk_unlock_request_contest_idx` ON `unlock_request` (`contest_id` ASC);
CREATE INDEX `fk_unlock_request_organizer_idx` ON `unlock_request` (`organizer_id` ASC);
CREATE INDEX `fk_unlock_request_status_idx` ON `unlock_request` (`status` ASC);
CREATE INDEX `fk_unlock_request_created_idx` ON `unlock_request` (`created_at` DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `unlock_request`;
-- +goose StatementEnd
