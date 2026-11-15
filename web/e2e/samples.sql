INSERT INTO organizer VALUES (NULL, 'Test Organizer');
INSERT INTO series VALUES (NULL, 1, 'Test series');
SET @first_competition_id = UNHEX('019a85f0000070008000000000000001');
SET @second_competition_id = UNHEX('019a85f0000070008000000000000002');
INSERT INTO competition VALUES (@first_competition_id, 1, NULL, 1, 'World Testing Championships', 'The world\'s number one competition for testing', 'On the web', "AQ", 10, 7, TRUE, FALSE, 0, 0, '<strong>Lorem ipsum</strong> dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.', 5, 20160, NOW());
INSERT INTO comp_class VALUES (NULL, 1, @first_competition_id, 'Males', '16 years and older', '#ff0000', '2024-01-01 00:00:00', '2026-12-31 23:59:59');
INSERT INTO comp_class VALUES (NULL, 1, @first_competition_id, 'Females', '16 years and older', '#ff0000', '2024-01-01 00:00:00', '2026-12-31 23:59:59');
INSERT INTO contender VALUES (NULL, 1, @first_competition_id, 'ABCD0001', 'Albert Einstein', 1, '2024-01-01 00:00:00', FALSE, FALSE, NULL, NULL);
INSERT INTO contender VALUES (NULL, 1, @first_competition_id, 'ABCD0002', NULL, NULL, NULL, FALSE, FALSE, NULL, NULL);
INSERT INTO contender VALUES (NULL, 1, @first_competition_id, 'ABCD0003', 'Michael Scott', 1, '2024-01-01 00:00:00', FALSE, FALSE, NULL, NULL);
INSERT INTO contender VALUES (NULL, 1, @first_competition_id, 'ABCD0004', NULL, NULL, NULL, FALSE, FALSE, NULL, NULL);
INSERT INTO contender VALUES (NULL, 1, @first_competition_id, 'ABCD0005', NULL, NULL, NULL, FALSE, FALSE, NULL, NULL);
INSERT INTO problem VALUES (NULL, 1, @first_competition_id, 1, '#ef4444', '#eab308', TRUE, TRUE, 'Very slippery', 10, 20, 100, 10);
INSERT INTO problem VALUES (NULL, 1, @first_competition_id, 2, '#f97316', NULL, TRUE, TRUE, NULL, 20, 40, 200, 10);
INSERT INTO problem VALUES (NULL, 1, @first_competition_id, 3, '#84cc16', NULL, TRUE, TRUE, NULL, 30, 60, 300, 10);
INSERT INTO problem VALUES (NULL, 1, @first_competition_id, 4, '#0ea5e9', NULL, TRUE, TRUE, NULL, 40, 80, 400, 10);
INSERT INTO problem VALUES (NULL, 1, @first_competition_id, 5, '#8b5cf6', NULL, TRUE, TRUE, NULL, 50, 100, 500, NULL);
INSERT INTO tick VALUES (NULL, 1, @first_competition_id, 1, 1, '2024-01-01 00:00:00', 1, TRUE, 999, TRUE, 999, TRUE, 999);

INSERT INTO competition VALUES (@second_competition_id, 1, NULL, 1, 'World Attempts Championships', 'A competition scored by tops, zones and attempts', 'On the web', 'AQ', 0, 0, FALSE, FALSE, 0, 0, NULL, 5, 20160, NOW());
INSERT INTO comp_class VALUES (NULL, 1, @second_competition_id, 'Open', NULL, '#ff0000', '2024-01-01 00:00:00', '2026-12-31 23:59:59');
INSERT INTO contender VALUES (NULL, 1, @second_competition_id, 'ABCD0006', 'Jim Halpert', 3, '2024-01-01 00:00:00', FALSE, FALSE, NULL, NULL);
INSERT INTO problem VALUES (NULL, 1, @second_competition_id, 1, '#ef4444', NULL, TRUE, TRUE, NULL, 0, 0, 0, NULL);

CREATE TABLE `goose_db_version` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `version_id` bigint(20) NOT NULL,
  `is_applied` tinyint(1) NOT NULL,
  `tstamp` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO goose_db_version SELECT NULL, seq, 1, NOW() FROM seq_1_to_43;
