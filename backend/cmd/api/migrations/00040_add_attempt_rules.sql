-- +goose Up
ALTER TABLE `contest`
    ADD COLUMN `max_attempts` INT NOT NULL DEFAULT 0 AFTER `pooled_points`,
    ADD COLUMN `point_deduction` INT NOT NULL DEFAULT 0 AFTER `max_attempts`;

-- +goose Down
ALTER TABLE `contest`
    DROP COLUMN `max_attempts`,
    DROP COLUMN `point_deduction`;
