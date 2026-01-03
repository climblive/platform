-- +goose Up
-- +goose StatementBegin
ALTER TABLE contest ADD COLUMN evaluation_mode TINYINT(1) NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE contest DROP COLUMN evaluation_mode;
-- +goose StatementEnd
