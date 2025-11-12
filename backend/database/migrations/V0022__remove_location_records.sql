ALTER TABLE contest ADD COLUMN `location` varchar(1024) DEFAULT NULL AFTER location_id;

UPDATE contest JOIN location ON contest.location_id = location.id SET location = location.name;

ALTER TABLE contest
    DROP CONSTRAINT fk_contest_1,
    DROP KEY fk_contest_1_idx,
    DROP COLUMN `location_id`;

DROP TABLE location;