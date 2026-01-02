-- +goose Up
ALTER TABLE `contest` CHANGE COLUMN `rules` `info` TEXT NULL;

-- +goose Down
ALTER TABLE `contest` CHANGE COLUMN `info` `rules` TEXT NULL;