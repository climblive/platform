package migrations

import (
	"context"
	"database/sql"
	"database/sql/driver"
	stderrors "errors"
	"fmt"
	"time"
	"uuid"

	"github.com/go-errors/errors"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(upCompetitionUUID, downCompetitionUUID)
}

var competitionReferenceTables = []string{"comp_class", "contender", "problem", "tick", "raffle"}

var competitionForeignKeys = []struct {
	table      string
	name       string
	definition string
}{
	{"competition", "fk_competition_2", "FOREIGN KEY (`organizer_id`) REFERENCES `organizer` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT"},
	{"competition", "fk_competition_3", "FOREIGN KEY (`series_id`, `organizer_id`) REFERENCES `series` (`id`, `organizer_id`) ON DELETE RESTRICT ON UPDATE RESTRICT"},
	{"comp_class", "fk_comp_class_1", "FOREIGN KEY (`competition_id`, `organizer_id`) REFERENCES `competition` (`id`, `organizer_id`) ON DELETE CASCADE ON UPDATE RESTRICT"},
	{"contender", "fk_contender_1", "FOREIGN KEY (`class_id`, `competition_id`) REFERENCES `comp_class` (`id`, `competition_id`) ON DELETE RESTRICT ON UPDATE RESTRICT"},
	{"contender", "fk_contender_2", "FOREIGN KEY (`competition_id`, `organizer_id`) REFERENCES `competition` (`id`, `organizer_id`) ON DELETE RESTRICT ON UPDATE RESTRICT"},
	{"problem", "fk_problem_1", "FOREIGN KEY (`competition_id`, `organizer_id`) REFERENCES `competition` (`id`, `organizer_id`) ON DELETE RESTRICT ON UPDATE RESTRICT"},
	{"tick", "fk_tick_1", "FOREIGN KEY (`problem_id`, `organizer_id`, `competition_id`) REFERENCES `problem` (`id`, `organizer_id`, `competition_id`) ON DELETE RESTRICT ON UPDATE RESTRICT"},
	{"tick", "fk_tick_2", "FOREIGN KEY (`contender_id`, `organizer_id`, `competition_id`) REFERENCES `contender` (`id`, `organizer_id`, `competition_id`) ON DELETE RESTRICT ON UPDATE RESTRICT"},
	{"raffle", "fk_raffle_1", "FOREIGN KEY (`competition_id`, `organizer_id`) REFERENCES `competition` (`id`, `organizer_id`) ON DELETE NO ACTION ON UPDATE NO ACTION"},
}

func upCompetitionUUID(ctx context.Context, db *sql.DB) (retErr error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			retErr = errors.Wrap(stderrors.Join(retErr, err), 0)
		}
	}()

	var foreignKeyChecks int
	if err := conn.QueryRowContext(ctx, "SELECT @@SESSION.FOREIGN_KEY_CHECKS").Scan(&foreignKeyChecks); err != nil {
		return errors.Wrap(err, 0)
	}

	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := conn.ExecContext(cleanupCtx, "SET SESSION FOREIGN_KEY_CHECKS = ?", foreignKeyChecks); err != nil {
			retErr = errors.Wrap(stderrors.Join(retErr, err), 0)
			_ = conn.Raw(func(any) error { return driver.ErrBadConn })
		}
	}()

	if _, err := conn.ExecContext(ctx, "SET SESSION FOREIGN_KEY_CHECKS = 0"); err != nil {
		return errors.Wrap(err, 0)
	}

	table := "contest"
	idType, err := competitionColumnType(ctx, conn, table, "id")
	if errors.Is(err, sql.ErrNoRows) {
		table = "competition"
		idType, err = competitionColumnType(ctx, conn, table, "id")
	}
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if idType == "int" {
		if err := mapCompetitionIDs(ctx, conn, table); err != nil {
			return err
		}
	}

	for _, fk := range competitionForeignKeys {
		fkTable, name := fk.table, fk.name
		if fkTable == "competition" {
			fkTable = table
			if table == "contest" {
				name = "fk_contest_" + name[len("fk_competition_"):]
			}
		}
		exists, err := competitionConstraintExists(ctx, conn, fkTable, name)
		if err != nil {
			return err
		}
		if exists {
			if _, err := conn.ExecContext(ctx, fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", fkTable, name)); err != nil {
				return errors.Wrap(err, 0)
			}
		}
	}

	if idType != "binary" {
		if _, err := conn.ExecContext(ctx, fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `id` BINARY(16) NOT NULL", table)); err != nil {
			return errors.Wrap(err, 0)
		}
	}

	for _, reference := range competitionReferenceTables {
		_, err := competitionColumnType(ctx, conn, reference, "contest_id")
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return errors.Wrap(err, 0)
		}
		if _, err := conn.ExecContext(ctx, fmt.Sprintf("ALTER TABLE `%s` CHANGE COLUMN `contest_id` `competition_id` BINARY(16) NOT NULL", reference)); err != nil {
			return errors.Wrap(err, 0)
		}
	}

	var hasMapping bool
	if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'competition_uuid_mapping'").Scan(&hasMapping); err != nil {
		return errors.Wrap(err, 0)
	}
	if hasMapping {
		for _, target := range append([]string{table}, competitionReferenceTables...) {
			column := "competition_id"
			if target == table {
				column = "id"
			}
			query := fmt.Sprintf("UPDATE `%s` t JOIN competition_uuid_mapping m ON t.`%s` = CAST(m.old_id AS BINARY(16)) SET t.`%s` = m.new_id", target, column, column)
			if _, err := conn.ExecContext(ctx, query); err != nil {
				return errors.Wrap(err, 0)
			}
		}
	}

	if table == "contest" {
		if _, err := conn.ExecContext(ctx, "RENAME TABLE `contest` TO `competition`"); err != nil {
			return errors.Wrap(err, 0)
		}
	}

	for _, index := range []string{"fk_contest_2_idx", "fk_contest_3"} {
		var exists bool
		if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) > 0 FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'competition' AND index_name = ?", index).Scan(&exists); err != nil {
			return errors.Wrap(err, 0)
		}
		if exists {
			name := "fk_competition_" + index[len("fk_contest_"):]
			if _, err := conn.ExecContext(ctx, fmt.Sprintf("ALTER TABLE competition RENAME INDEX `%s` TO `%s`", index, name)); err != nil {
				return errors.Wrap(err, 0)
			}
		}
	}

	if _, err := conn.ExecContext(ctx, "SET SESSION FOREIGN_KEY_CHECKS = 1"); err != nil {
		return errors.Wrap(err, 0)
	}
	for _, fk := range competitionForeignKeys {
		if _, err := conn.ExecContext(ctx, fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `%s` %s", fk.table, fk.name, fk.definition)); err != nil {
			return errors.Wrap(err, 0)
		}
	}

	if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS competition_uuid_mapping"); err != nil {
		return errors.Wrap(err, 0)
	}
	return nil
}

func mapCompetitionIDs(ctx context.Context, conn *sql.Conn, table string) error {
	if _, err := conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS competition_uuid_mapping (old_id INT NOT NULL PRIMARY KEY, new_id BINARY(16) NOT NULL UNIQUE) ENGINE=InnoDB"); err != nil {
		return errors.Wrap(err, 0)
	}
	rows, err := conn.QueryContext(ctx, fmt.Sprintf("SELECT id FROM `%s` WHERE id NOT IN (SELECT old_id FROM competition_uuid_mapping)", table))
	if err != nil {
		return errors.Wrap(err, 0)
	}
	defer func() { _ = rows.Close() }()
	var ids []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return errors.Wrap(err, 0)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return errors.Wrap(err, 0)
	}
	if err := rows.Close(); err != nil {
		return errors.Wrap(err, 0)
	}
	for _, id := range ids {
		newID := uuid.NewV7()
		if _, err := conn.ExecContext(ctx, "INSERT INTO competition_uuid_mapping (old_id, new_id) VALUES (?, ?)", id, newID[:]); err != nil {
			return errors.Wrap(err, 0)
		}
	}
	return nil
}

func competitionColumnType(ctx context.Context, conn *sql.Conn, table, column string) (string, error) {
	var dataType string
	err := conn.QueryRowContext(ctx, "SELECT data_type FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?", table, column).Scan(&dataType)
	if err != nil {
		return "", errors.Wrap(err, 0)
	}
	return dataType, nil
}

func competitionConstraintExists(ctx context.Context, conn *sql.Conn, table, name string) (bool, error) {
	var exists bool
	err := conn.QueryRowContext(ctx, "SELECT COUNT(*) > 0 FROM information_schema.table_constraints WHERE constraint_schema = DATABASE() AND table_name = ? AND constraint_name = ? AND constraint_type = 'FOREIGN KEY'", table, name).Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, 0)
	}
	return exists, nil
}

func downCompetitionUUID(context.Context, *sql.DB) error {
	return errors.New("competition UUID migration is irreversible; restore a pre-migration backup to downgrade")
}
