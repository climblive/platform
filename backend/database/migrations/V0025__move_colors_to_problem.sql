ALTER TABLE problem ADD COLUMN `hold_color_primary` VARCHAR(7) NOT NULL AFTER number;
ALTER TABLE problem ADD COLUMN `hold_color_secondary` VARCHAR(7) NULL DEFAULT NULL AFTER hold_color_primary;

UPDATE problem
    JOIN color ON problem.color_id = color.id
    SET hold_color_primary = color.rgb_primary, hold_color_secondary = color.rgb_secondary;

UPDATE problem
    JOIN color ON problem.color_id = color.id
    SET problem.name = color.name
    WHERE problem.name IS NULL;

ALTER TABLE problem
    DROP CONSTRAINT `fk_problem_2`,
    DROP KEY `fk_problem_2_idx`,
    DROP COLUMN `color_id`;

DROP TABLE color;