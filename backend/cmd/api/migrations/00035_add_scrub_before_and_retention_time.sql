-- +goose Up
ALTER TABLE `contest` ADD COLUMN `name_retention_time` INT NOT NULL DEFAULT 604800;
ALTER TABLE `contender` ADD COLUMN `scrub_before` TIMESTAMP NULL DEFAULT NULL;

UPDATE contender c
JOIN comp_class cc ON c.class_id = cc.id
JOIN contest ct ON c.contest_id = ct.id
SET c.scrub_before = DATE_ADD(cc.time_end, INTERVAL ct.name_retention_time SECOND)
WHERE c.class_id IS NOT NULL AND c.entered IS NOT NULL;

-- +goose Down
ALTER TABLE `contender` DROP COLUMN `scrub_before`;
ALTER TABLE `contest` DROP COLUMN `name_retention_time`;