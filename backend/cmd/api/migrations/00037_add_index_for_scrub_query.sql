-- +goose Up
CREATE INDEX `idx_contender_scrub` ON `contender` (`name`, `scrub_before`);

-- +goose Down
DROP INDEX `idx_contender_scrub` ON `contender`;
