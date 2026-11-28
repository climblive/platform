-- +goose Up
ALTER TABLE contest
    ADD COLUMN duplicated_by TIMESTAMP NULL DEFAULT NULL AFTER created,
    ADD COLUMN duplicated_from INT NULL DEFAULT NULL AFTER duplicated_by;

-- +goose Down
ALTER TABLE contest
    DROP COLUMN duplicated_by,
    DROP COLUMN duplicated_from;
