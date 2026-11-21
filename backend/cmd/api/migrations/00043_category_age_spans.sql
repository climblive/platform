-- +goose Up
ALTER TABLE comp_class
    ADD COLUMN min_age INT NULL AFTER description,
    ADD COLUMN max_age INT NULL AFTER min_age;

UPDATE comp_class SET min_age = 18, max_age = 99;

ALTER TABLE comp_class
    MODIFY COLUMN min_age INT NOT NULL,
    MODIFY COLUMN max_age INT NOT NULL;

-- +goose Down
ALTER TABLE comp_class
    DROP COLUMN max_age,
    DROP COLUMN min_age;
