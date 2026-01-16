package migrations

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upConvertIdsToUuid, downConvertIdsToUuid)
}

func upConvertIdsToUuid(tx *sql.Tx) error {
	if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return fmt.Errorf("disabling foreign key checks: %w", err)
	}
	defer func() {
		tx.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}()

	tables := []string{
		"organizer",
		"series",
		"contest",
		"comp_class",
		"contender",
		"problem",
		"tick",
		"user",
		"raffle",
		"raffle_winner",
	}

	type idMapping struct {
		oldID int
		newID string
	}

	mappings := make(map[string][]idMapping)

	for _, table := range tables {
		if err := addUUIDColumn(tx, table); err != nil {
			return fmt.Errorf("adding UUID column to %s: %w", table, err)
		}

		if err := generateAndStoreUUIDs(tx, table, mappings); err != nil {
			return fmt.Errorf("generating UUIDs for %s: %w", table, err)
		}
	}

	fkUpdates := []struct {
		table   string
		column  string
		refTable string
	}{
		{"series", "organizer_id", "organizer"},
		{"contest", "organizer_id", "organizer"},
		{"contest", "series_id", "series"},
		{"comp_class", "organizer_id", "organizer"},
		{"comp_class", "contest_id", "contest"},
		{"contender", "organizer_id", "organizer"},
		{"contender", "contest_id", "contest"},
		{"contender", "class_id", "comp_class"},
		{"problem", "organizer_id", "organizer"},
		{"problem", "contest_id", "contest"},
		{"tick", "organizer_id", "organizer"},
		{"tick", "contest_id", "contest"},
		{"tick", "contender_id", "contender"},
		{"tick", "problem_id", "problem"},
		{"user_organizer", "organizer_id", "organizer"},
		{"user_organizer", "user_id", "user"},
		{"raffle", "organizer_id", "organizer"},
		{"raffle", "contest_id", "contest"},
		{"raffle_winner", "organizer_id", "organizer"},
		{"raffle_winner", "raffle_id", "raffle"},
		{"raffle_winner", "contender_id", "contender"},
		{"score", "contender_id", "contender"},
		{"organizer_invite", "organizer_id", "organizer"},
	}

	for _, fk := range fkUpdates {
		if err := addUUIDFKColumn(tx, fk.table, fk.column); err != nil {
			return fmt.Errorf("adding UUID FK column %s.%s: %w", fk.table, fk.column, err)
		}

		if err := updateFKReferences(tx, fk.table, fk.column, fk.refTable, mappings[fk.refTable]); err != nil {
			return fmt.Errorf("updating FK references %s.%s: %w", fk.table, fk.column, err)
		}
	}

	if err := dropOldConstraints(tx); err != nil {
		return fmt.Errorf("dropping old constraints: %w", err)
	}

	if err := replaceColumns(tx, tables, fkUpdates); err != nil {
		return fmt.Errorf("replacing columns: %w", err)
	}

	if err := recreateConstraints(tx); err != nil {
		return fmt.Errorf("recreating constraints: %w", err)
	}

	return nil
}

func addUUIDColumn(tx *sql.Tx, table string) error {
	_, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `id_new` BINARY(16) NULL", table))
	return err
}

func generateAndStoreUUIDs(tx *sql.Tx, table string, mappings map[string][]idMapping) error {
	rows, err := tx.Query(fmt.Sprintf("SELECT id FROM `%s`", table))
	if err != nil {
		return err
	}
	defer rows.Close()

	mappings[table] = []idMapping{}

	for rows.Next() {
		var oldID int
		if err := rows.Scan(&oldID); err != nil {
			return err
		}

		newUUID := uuid.New()
		uuidBytes, err := newUUID.MarshalBinary()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(
			fmt.Sprintf("UPDATE `%s` SET id_new = ? WHERE id = ?", table),
			uuidBytes, oldID,
		); err != nil {
			return err
		}

		mappings[table] = append(mappings[table], idMapping{
			oldID: oldID,
			newID: newUUID.String(),
		})
	}

	return rows.Err()
}

func addUUIDFKColumn(tx *sql.Tx, table, column string) error {
	_, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s_new` BINARY(16) NULL", table, column))
	return err
}

func updateFKReferences(tx *sql.Tx, table, column, refTable string, mapping []idMapping) error {
	for _, m := range mapping {
		newUUID, err := uuid.Parse(m.newID)
		if err != nil {
			return err
		}
		uuidBytes, err := newUUID.MarshalBinary()
		if err != nil {
			return err
		}

		_, err = tx.Exec(
			fmt.Sprintf("UPDATE `%s` SET `%s_new` = ? WHERE `%s` = ?", table, column, column),
			uuidBytes, m.oldID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropOldConstraints(tx *sql.Tx) error {
	constraints := []struct {
		table      string
		constraint string
	}{
		{"series", "fk_series_1"},
		{"contest", "fk_contest_2"},
		{"contest", "fk_contest_3"},
		{"comp_class", "fk_comp_class_1"},
		{"contender", "fk_contender_1"},
		{"contender", "fk_contender_2"},
		{"problem", "fk_problem_1"},
		{"tick", "fk_tick_1"},
		{"tick", "fk_tick_2"},
		{"user_organizer", "fk_user_organizer_1"},
		{"user_organizer", "fk_user_organizer_2"},
		{"raffle", "fk_raffle_1"},
		{"raffle_winner", "fk_raffle_winner_1"},
		{"raffle_winner", "fk_raffle_winner_2"},
		{"score", "fk_score_1"},
		{"organizer_invite", "fk_organizer_invite_1"},
	}

	for _, c := range constraints {
		_, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", c.table, c.constraint))
		if err != nil {
			return fmt.Errorf("dropping constraint %s.%s: %w", c.table, c.constraint, err)
		}
	}

	return nil
}

func replaceColumns(tx *sql.Tx, tables []string, fkUpdates []struct {
	table   string
	column  string
	refTable string
}) error {
	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `id`", table)); err != nil {
			return fmt.Errorf("dropping old id column from %s: %w", table, err)
		}

		if _, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` CHANGE COLUMN `id_new` `id` BINARY(16) NOT NULL", table)); err != nil {
			return fmt.Errorf("renaming id_new to id in %s: %w", table, err)
		}

		if _, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)", table)); err != nil {
			return fmt.Errorf("recreating primary key for %s: %w", table, err)
		}
	}

	for _, fk := range fkUpdates {
		var checkColumnExists string
		err := tx.QueryRow(fmt.Sprintf(
			"SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = '%s' AND COLUMN_NAME = '%s'",
			fk.table, fk.column,
		)).Scan(&checkColumnExists)
		if err != nil {
			return err
		}

		if checkColumnExists != "0" {
			if _, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", fk.table, fk.column)); err != nil {
				return fmt.Errorf("dropping old FK column %s.%s: %w", fk.table, fk.column, err)
			}
		}

		if _, err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` CHANGE COLUMN `%s_new` `%s` BINARY(16) NOT NULL", fk.table, fk.column, fk.column)); err != nil {
			return fmt.Errorf("renaming FK column %s.%s: %w", fk.table, fk.column, err)
		}
	}

	return nil
}

func recreateConstraints(tx *sql.Tx) error {
	constraints := []string{
		"ALTER TABLE `series` ADD CONSTRAINT `fk_series_1` FOREIGN KEY (`organizer_id`) REFERENCES `organizer` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT",
		"ALTER TABLE `contest` ADD CONSTRAINT `fk_contest_2` FOREIGN KEY (`organizer_id`) REFERENCES `organizer` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT",
		"ALTER TABLE `contest` ADD CONSTRAINT `fk_contest_3` FOREIGN KEY (`series_id`, `organizer_id`) REFERENCES `series` (`id`, `organizer_id`) ON DELETE RESTRICT ON UPDATE RESTRICT",
		"ALTER TABLE `comp_class` ADD CONSTRAINT `fk_comp_class_1` FOREIGN KEY (`contest_id`, `organizer_id`) REFERENCES `contest` (`id`, `organizer_id`) ON DELETE CASCADE ON UPDATE RESTRICT",
		"ALTER TABLE `contender` ADD CONSTRAINT `fk_contender_1` FOREIGN KEY (`class_id`, `contest_id`) REFERENCES `comp_class` (`id`, `contest_id`) ON DELETE RESTRICT ON UPDATE RESTRICT",
		"ALTER TABLE `contender` ADD CONSTRAINT `fk_contender_2` FOREIGN KEY (`contest_id`, `organizer_id`) REFERENCES `contest` (`id`, `organizer_id`) ON DELETE RESTRICT ON UPDATE RESTRICT",
		"ALTER TABLE `problem` ADD CONSTRAINT `fk_problem_1` FOREIGN KEY (`contest_id`, `organizer_id`) REFERENCES `contest` (`id`, `organizer_id`) ON DELETE RESTRICT ON UPDATE RESTRICT",
		"ALTER TABLE `tick` ADD CONSTRAINT `fk_tick_1` FOREIGN KEY (`problem_id`, `organizer_id`, `contest_id`) REFERENCES `problem` (`id`, `organizer_id`, `contest_id`) ON DELETE RESTRICT ON UPDATE RESTRICT",
		"ALTER TABLE `tick` ADD CONSTRAINT `fk_tick_2` FOREIGN KEY (`contender_id`, `organizer_id`, `contest_id`) REFERENCES `contender` (`id`, `organizer_id`, `contest_id`) ON DELETE RESTRICT ON UPDATE RESTRICT",
		"ALTER TABLE `user_organizer` ADD CONSTRAINT `fk_user_organizer_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT",
		"ALTER TABLE `user_organizer` ADD CONSTRAINT `fk_user_organizer_2` FOREIGN KEY (`organizer_id`) REFERENCES `organizer` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT",
		"ALTER TABLE `raffle` ADD CONSTRAINT `fk_raffle_1` FOREIGN KEY (`contest_id`, `organizer_id`) REFERENCES `contest` (`id`, `organizer_id`) ON DELETE NO ACTION ON UPDATE NO ACTION",
		"ALTER TABLE `raffle_winner` ADD CONSTRAINT `fk_raffle_winner_1` FOREIGN KEY (`raffle_id`, `organizer_id`) REFERENCES `raffle` (`id`, `organizer_id`) ON DELETE NO ACTION ON UPDATE NO ACTION",
		"ALTER TABLE `raffle_winner` ADD CONSTRAINT `fk_raffle_winner_2` FOREIGN KEY (`contender_id`, `organizer_id`) REFERENCES `contender` (`id`, `organizer_id`) ON DELETE NO ACTION ON UPDATE NO ACTION",
		"ALTER TABLE `score` ADD CONSTRAINT `fk_score_1` FOREIGN KEY (`contender_id`) REFERENCES `contender` (`id`) ON DELETE CASCADE ON UPDATE CASCADE",
		"ALTER TABLE `organizer_invite` ADD CONSTRAINT `fk_organizer_invite_1` FOREIGN KEY (`organizer_id`) REFERENCES `organizer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE",
	}

	for _, constraint := range constraints {
		if _, err := tx.Exec(constraint); err != nil {
			return fmt.Errorf("recreating constraint: %s: %w", constraint, err)
		}
	}

	return nil
}

func downConvertIdsToUuid(tx *sql.Tx) error {
	return fmt.Errorf("downgrade not supported for UUID conversion")
}
