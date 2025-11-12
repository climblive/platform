ALTER TABLE `color`
    CHANGE COLUMN `rgb` `rgb_primary` VARCHAR(7) NOT NULL,
    ADD COLUMN `rgb_secondary` VARCHAR(7) NULL AFTER `rgb_primary`;

UPDATE `color` AS c1 LEFT JOIN `color` AS c2 ON c1.id = c2.id SET c1.`rgb_primary` = CONCAT('#', c2.`rgb_primary`);
