-- +goose Up
ALTER TABLE `contest`
    ADD COLUMN `use_points` TINYINT(1) NOT NULL DEFAULT 1 AFTER `finalists`,
    ADD COLUMN `pooled_points` TINYINT(1) NOT NULL DEFAULT 0 AFTER `use_points`;

-- +goose Down
ALTER TABLE `contest`
    DROP COLUMN `use_points`,
    DROP COLUMN `pooled_points`;
