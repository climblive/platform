-- +goose Up
ALTER TABLE contest ADD COLUMN evaluation_mode TINYINT(1) NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE contest DROP COLUMN evaluation_mode;
