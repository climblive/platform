ALTER TABLE `contest`
    ADD COLUMN `grace_period` INT NOT NULL DEFAULT 300 AFTER `rules`;