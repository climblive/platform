-- +goose Up
DROP INDEX `index3` ON `problem`;

ALTER TABLE `problem`
    CHANGE COLUMN `number` `label` VARCHAR(6) NOT NULL;

CREATE UNIQUE INDEX `index3` ON `problem` (`label` ASC, `contest_id` ASC);

-- +goose Down
DROP INDEX `index3` ON `problem`;

UPDATE `problem`
JOIN (
    SELECT `id`, ROW_NUMBER() OVER (PARTITION BY `contest_id` ORDER BY `label`, `id`) AS `number`
    FROM `problem`
) AS `numbered` ON `problem`.`id` = `numbered`.`id`
SET `problem`.`label` = `numbered`.`number`;

ALTER TABLE `problem`
    CHANGE COLUMN `label` `number` INT NOT NULL;

CREATE UNIQUE INDEX `index3` ON `problem` (`number` ASC, `contest_id` ASC);
