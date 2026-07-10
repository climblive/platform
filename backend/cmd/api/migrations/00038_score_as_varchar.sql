-- +goose Up
ALTER TABLE `score`
    CHANGE COLUMN `score` `score` VARCHAR(128) NOT NULL;

UPDATE score SET `score` = CONCAT(score, "p");

-- +goose Down
UPDATE score SET `score` = TRIM(TRAILING, "p", FROM score)

ALTER TABLE `score`
    CHANGE COLUMN `score` `score` INT NOT NULL;
