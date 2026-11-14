-- +goose Up
SET @problem_number_sequence_sql = (
    SELECT CONCAT(
        'CREATE SEQUENCE `problem_number_seq` START WITH ',
        GREATEST(COALESCE(MAX(`number`), 0) + 1, 1),
        ' MINVALUE 1 MAXVALUE 2147483647 NOCYCLE'
    ) FROM `problem`
);
PREPARE problem_number_sequence_stmt FROM @problem_number_sequence_sql;
EXECUTE problem_number_sequence_stmt;
DEALLOCATE PREPARE problem_number_sequence_stmt;

ALTER TABLE `problem`
    ADD COLUMN `label` VARCHAR(6) NULL AFTER `number`,
    ALTER COLUMN `number` SET DEFAULT (NEXT VALUE FOR `problem_number_seq`);

UPDATE `problem` SET `label` = CAST(`number` AS CHAR);

ALTER TABLE `problem` MODIFY COLUMN `label` VARCHAR(6) NOT NULL;

CREATE UNIQUE INDEX `label_contest_unique` ON `problem` (`label` ASC, `contest_id` ASC);

-- +goose Down
DROP INDEX `label_contest_unique` ON `problem`;

ALTER TABLE `problem`
    DROP COLUMN `label`,
    ALTER COLUMN `number` DROP DEFAULT;

DROP SEQUENCE `problem_number_seq`;
