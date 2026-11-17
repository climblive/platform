-- +goose Up
ALTER TABLE `contest`
    ADD COLUMN `favor_top_attempts` TINYINT(1) NOT NULL DEFAULT 0 AFTER `use_points`;

UPDATE `contest` SET `favor_top_attempts` = TRUE WHERE `use_points` = FALSE;

-- +goose Down
ALTER TABLE `contest`
    DROP COLUMN `favor_top_attempts`;
