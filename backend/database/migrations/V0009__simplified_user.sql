ALTER TABLE `user`
    CHANGE COLUMN `email` `username` VARCHAR(64) NOT NULL,
    DROP COLUMN `password`;