-- +goose Up
ALTER TABLE `tick`
    ADD COLUMN `revision` INT UNSIGNED NOT NULL DEFAULT 0 AFTER `timestamp`;

-- +goose Down
ALTER TABLE `tick`
    DROP COLUMN `revision`;
