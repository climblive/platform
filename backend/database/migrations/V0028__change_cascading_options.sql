ALTER TABLE series
    DROP CONSTRAINT `fk_series_1`;

ALTER TABLE series
    ADD CONSTRAINT `fk_series_1` FOREIGN KEY `fk_series_1`
    FOREIGN KEY (`organizer_id`)
    REFERENCES `organizer` (`id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE);

ALTER TABLE contest
    DROP CONSTRAINT `fk_contest_2`;

ALTER TABLE contest
    DROP CONSTRAINT `fk_contest_3`;

ALTER TABLE contest
    ADD CONSTRAINT `fk_contest_2`
    CONSTRAINT `fk_contest_2`
    FOREIGN KEY (`organizer_id`)
    REFERENCES `organizer` (`id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE contest
    ADD CONSTRAINT `fk_contest_3`
    FOREIGN KEY (`series_id` , `organizer_id`)
    REFERENCES `series` (`id` , `organizer_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE contender
    DROP CONSTRAINT `fk_contender_1`;

ALTER TABLE contender
    DROP CONSTRAINT `fk_contender_2`;

ALTER TABLE contender
    ADD CONSTRAINT `fk_contender_1`
    FOREIGN KEY (`class_id` , `contest_id`)
    REFERENCES `comp_class` (`id` , `contest_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE contender
    ADD CONSTRAINT `fk_contender_2`
    FOREIGN KEY (`contest_id` , `organizer_id`)
    REFERENCES `contest` (`id` , `organizer_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE problem
    DROP CONSTRAINT fk_problem_1;

ALTER TABLE problem
    ADD CONSTRAINT fk_problem_1
    FOREIGN KEY (`contest_id` , `organizer_id`)
    REFERENCES `contest` (`id` , `organizer_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE tick
    DROP CONSTRAINT `fk_tick_1`;

ALTER TABLE tick
    DROP CONSTRAINT `fk_tick_2`;

ALTER TABLE tick
    ADD CONSTRAINT `fk_tick_1`
    FOREIGN KEY (`problem_id` , `organizer_id` , `contest_id`)
    REFERENCES `problem` (`id` , `organizer_id` , `contest_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE tick
    ADD CONSTRAINT `fk_tick_2`
    FOREIGN KEY (`contender_id` , `organizer_id` , `contest_id`)
    REFERENCES `contender` (`id` , `organizer_id` , `contest_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE user_organizer
    DROP CONSTRAINT `fk_user_organizer_1`;

ALTER TABLE user_organizer
    DROP CONSTRAINT `fk_user_organizer_2`;

ALTER TABLE user_organizer
    ADD CONSTRAINT `fk_user_organizer_1`
    FOREIGN KEY (`user_id`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE user_organizer
    ADD CONSTRAINT `fk_user_organizer_2`
    FOREIGN KEY (`organizer_id`)
    REFERENCES `organizer` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE raffle
    DROP CONSTRAINT `fk_raffle_1`;

ALTER TABLE raffle
    ADD CONSTRAINT `fk_raffle_1`
    FOREIGN KEY (`contest_id` , `organizer_id`)
    REFERENCES `contest` (`id` , `organizer_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE raffle_winner
    DROP CONSTRAINT `fk_raffle_winner_1`;

ALTER TABLE raffle_winner
    DROP CONSTRAINT `fk_raffle_winner_2`;

ALTER TABLE raffle_winner
    ADD CONSTRAINT `fk_raffle_winner_1`
    FOREIGN KEY (`raffle_id` , `organizer_id`)
    REFERENCES `raffle` (`id` , `organizer_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

ALTER TABLE raffle_winner
    ADD CONSTRAINT `fk_raffle_winner_2`
    FOREIGN KEY (`contender_id` , `organizer_id`)
    REFERENCES `contender` (`id` , `organizer_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;