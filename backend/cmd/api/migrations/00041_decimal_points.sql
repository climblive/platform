-- +goose Up
ALTER TABLE `contest` MODIFY COLUMN `point_deduction` DOUBLE NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE `contest` MODIFY COLUMN `point_deduction` INT NOT NULL DEFAULT 0;
