ALTER TABLE `raffle`
    ADD COLUMN `active` TINYINT(1) NOT NULL DEFAULT 0 AFTER `contest_id`;
