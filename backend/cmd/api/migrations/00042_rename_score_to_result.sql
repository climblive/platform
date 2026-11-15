-- +goose Up
RENAME TABLE `score` TO `result`;
ALTER TABLE `result`
  DROP FOREIGN KEY `fk_score_1`,
  ADD CONSTRAINT `fk_result_1` FOREIGN KEY (`contender_id`) REFERENCES `contender` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose Down
RENAME TABLE `result` TO `score`;
ALTER TABLE `score`
  DROP FOREIGN KEY `fk_result_1`,
  ADD CONSTRAINT `fk_score_1` FOREIGN KEY (`contender_id`) REFERENCES `contender` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE;
