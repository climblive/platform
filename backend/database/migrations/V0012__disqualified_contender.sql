ALTER TABLE `contender`
    ADD COLUMN `disqualified` TINYINT(1) NOT NULL DEFAULT 0 AFTER `entered`;