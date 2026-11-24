-- +goose Up
ALTER TABLE comp_class
    ADD COLUMN eligibility_basis ENUM('age', 'birth_year') NOT NULL DEFAULT 'age' AFTER description,
    MODIFY COLUMN min_age INT NULL,
    MODIFY COLUMN max_age INT NULL,
    ADD COLUMN min_birth_year INT NULL AFTER max_age,
    ADD COLUMN max_birth_year INT NULL AFTER min_birth_year;

-- +goose Down
UPDATE comp_class SET min_age = 18, max_age = 99 WHERE eligibility_basis = 'birth_year';
UPDATE comp_class SET min_age = COALESCE(min_age, 0), max_age = COALESCE(max_age, 2147483647);

ALTER TABLE comp_class
    DROP COLUMN eligibility_basis,
    DROP COLUMN min_birth_year,
    DROP COLUMN max_birth_year,
    MODIFY COLUMN min_age INT NOT NULL,
    MODIFY COLUMN max_age INT NOT NULL;
