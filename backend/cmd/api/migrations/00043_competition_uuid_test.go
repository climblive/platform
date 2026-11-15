package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed *.sql
var competitionMigrationSQL embed.FS

func TestCompetitionUUIDMigration(t *testing.T) {
	for _, scenario := range []string{"populated", "empty", "resume", "partial conversion", "invalid reference", "disabled checks"} {
		t.Run(scenario, func(t *testing.T) {
			ctx := t.Context()
			db := competitionMigrationDatabase(t)
			if scenario != "empty" {
				seedCompetitionMigration(t, db)
			}
			var expectedID domain.CompetitionID

			if scenario == "partial conversion" {
				conn, err := db.Conn(ctx)
				require.NoError(t, err)
				t.Cleanup(func() { _ = conn.Close() })
				require.NoError(t, mapCompetitionIDs(ctx, conn, "contest"))
				_, err = conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0")
				require.NoError(t, err)
				for _, fk := range competitionForeignKeys {
					if fk.table == "competition" {
						continue
					}
					_, err := conn.ExecContext(ctx, fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", fk.table, fk.name))
					require.NoError(t, err)
				}
				_, err = conn.ExecContext(ctx, "ALTER TABLE contest MODIFY COLUMN id BINARY(16) NOT NULL")
				require.NoError(t, err)
				_, err = conn.ExecContext(ctx, "UPDATE contest c JOIN competition_uuid_mapping m ON c.id = CAST(m.old_id AS BINARY(16)) SET c.id = m.new_id")
				require.NoError(t, err)
				_, err = conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1")
				require.NoError(t, err)
				require.NoError(t, conn.Close())
				require.NoError(t, db.QueryRowContext(ctx, "SELECT new_id FROM competition_uuid_mapping WHERE old_id = 1").Scan(&expectedID))
			}

			if scenario == "resume" {
				_, err := db.ExecContext(ctx, "ALTER TABLE contest ADD INDEX fk_competition_2_idx (organizer_id)")
				require.NoError(t, err)
				require.Error(t, upCompetitionUUID(ctx, db))
				require.NoError(t, db.QueryRowContext(ctx, "SELECT new_id FROM competition_uuid_mapping WHERE old_id = 1").Scan(&expectedID))
				var checks int
				require.NoError(t, db.QueryRowContext(ctx, "SELECT @@SESSION.FOREIGN_KEY_CHECKS").Scan(&checks))
				assert.Equal(t, 1, checks)
				_, err = db.ExecContext(ctx, "ALTER TABLE competition DROP INDEX fk_competition_2_idx")
				require.NoError(t, err)
			}

			if scenario == "invalid reference" {
				_, err := db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0")
				require.NoError(t, err)
				_, err = db.ExecContext(ctx, "INSERT INTO raffle (id, organizer_id, contest_id) VALUES (99, 1, 999)")
				require.NoError(t, err)
				_, err = db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1")
				require.NoError(t, err)
				require.Error(t, upCompetitionUUID(ctx, db))
				var checks int
				require.NoError(t, db.QueryRowContext(ctx, "SELECT @@SESSION.FOREIGN_KEY_CHECKS").Scan(&checks))
				assert.Equal(t, 1, checks)
				_, err = db.ExecContext(ctx, "DELETE FROM raffle WHERE id = 99")
				require.NoError(t, err)
			}

			if scenario == "disabled checks" {
				_, err := db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0")
				require.NoError(t, err)
			}
			require.NoError(t, goose.UpContext(ctx, db, "."))
			if scenario == "disabled checks" {
				var checks int
				require.NoError(t, db.QueryRowContext(ctx, "SELECT @@SESSION.FOREIGN_KEY_CHECKS").Scan(&checks))
				assert.Zero(t, checks)
				_, err := db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1")
				require.NoError(t, err)
			}
			version, err := goose.GetDBVersionContext(ctx, db)
			require.NoError(t, err)
			assert.EqualValues(t, 43, version)
			require.NoError(t, upCompetitionUUID(ctx, db))
			require.ErrorContains(t, downCompetitionUUID(ctx, db), "irreversible")
			if expectedID != (domain.CompetitionID{}) {
				var actualID domain.CompetitionID
				require.NoError(t, db.QueryRowContext(ctx, "SELECT id FROM competition WHERE name = 'First'").Scan(&actualID))
				assert.Equal(t, expectedID, actualID)
			}

			var legacyNames int
			require.NoError(t, db.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND (table_name = 'contest' OR column_name = 'contest_id')").Scan(&legacyNames))
			assert.Zero(t, legacyNames)
			var binaryColumns int
			require.NoError(t, db.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND (column_name = 'competition_id' OR (table_name = 'competition' AND column_name = 'id')) AND column_type = 'binary(16)'").Scan(&binaryColumns))
			assert.Equal(t, 6, binaryColumns)
			var numericPrimaryKeys int
			require.NoError(t, db.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND column_name = 'id' AND data_type = 'int' AND extra = 'auto_increment'").Scan(&numericPrimaryKeys))
			assert.Equal(t, 9, numericPrimaryKeys)

			if scenario != "empty" {
				verifyCompetitionMigration(t, db)
			}
		})
	}
}

func competitionMigrationDatabase(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("MIGRATION_TEST_DSN")
	if dsn == "" {
		t.Skip("set MIGRATION_TEST_DSN to a MySQL/MariaDB server with CREATE DATABASE privileges")
	}
	cfg, err := mysql.ParseDSN(dsn)
	require.NoError(t, err)
	cfg.DBName = ""
	cfg.ParseTime = true
	admin, err := sql.Open("mysql", cfg.FormatDSN())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, admin.Close()) })
	name := "competition_uuid_test_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	_, err = admin.ExecContext(t.Context(), "CREATE DATABASE `"+name+"` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := admin.ExecContext(context.Background(), "DROP DATABASE `"+name+"`")
		assert.NoError(t, err)
	})
	cfg.DBName = name
	db, err := sql.Open("mysql", cfg.FormatDSN())
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	goose.SetBaseFS(competitionMigrationSQL)
	t.Cleanup(func() { goose.SetBaseFS(nil) })
	require.NoError(t, goose.SetDialect("mysql"))
	require.NoError(t, goose.UpToContext(t.Context(), db, ".", 42))
	return db
}

func seedCompetitionMigration(t *testing.T, db *sql.DB) {
	t.Helper()
	queries := []string{
		"INSERT INTO organizer (id, name) VALUES (1, 'Organizer')",
		"INSERT INTO series (id, organizer_id, name) VALUES (1, 1, 'Series')",
		"INSERT INTO contest (id, organizer_id, series_id, name, qualifying_problems, finalists) VALUES (1, 1, 1, 'First', 2, 1), (1337, 1, NULL, 'Second', 2, 1)",
		"INSERT INTO comp_class (id, organizer_id, contest_id, name) VALUES (1, 1, 1, 'First category'), (2, 1, 1337, 'Second category')",
		"INSERT INTO contender (id, organizer_id, contest_id, registration_code, class_id, name) VALUES (1, 1, 1, 'FIRST', 1, 'First contender'), (2, 1, 1337, 'SECOND', 2, 'Second contender'), (3, 1, 1337, 'UNREGISTERED', NULL, NULL)",
		"INSERT INTO problem (id, organizer_id, contest_id, number, hold_color_primary, points_top) VALUES (1, 1, 1, 1, '#ffffff', 100), (2, 1, 1337, 1, '#ffffff', 100)",
		"INSERT INTO tick (id, organizer_id, contest_id, contender_id, problem_id) VALUES (1, 1, 1, 1, 1), (2, 1, 1337, 2, 2)",
		"INSERT INTO raffle (id, organizer_id, contest_id) VALUES (1, 1, 1), (2, 1, 1337)",
		"INSERT INTO raffle_winner (id, organizer_id, raffle_id, contender_id, timestamp) VALUES (1, 1, 1, 1, NOW()), (2, 1, 2, 2, NOW())",
		"INSERT INTO result (contender_id, timestamp, score, placement, finalist, rank_order) VALUES (1, NOW(), '100p', 1, 1, 1)",
	}
	for _, query := range queries {
		_, err := db.ExecContext(t.Context(), query)
		require.NoError(t, err)
	}
}

func verifyCompetitionMigration(t *testing.T, db *sql.DB) {
	t.Helper()
	ctx := t.Context()
	queries := database.New(db)
	competitions, err := queries.GetAllCompetitions(ctx)
	require.NoError(t, err)
	require.Len(t, competitions, 2)
	assert.NotEqual(t, competitions[0].Competition.ID, competitions[1].Competition.ID)
	for _, record := range competitions {
		competition := record.Competition
		assert.NotEqual(t, domain.CompetitionID{}, competition.ID)
		_, err := uuid.Parse(competition.ID.String())
		require.NoError(t, err)
		for _, table := range competitionReferenceTables {
			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE competition_id = ?", table)
			require.NoError(t, db.QueryRowContext(ctx, query, competition.ID).Scan(&count))
			expected := 1
			if table == "contender" && competition.Name == "Second" {
				expected = 2
			}
			assert.Equal(t, expected, count, table)
		}
	}
	contender, err := queries.GetCompetitor(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "First contender", contender.Competitor.Name.String)
	assert.Equal(t, "100p", contender.Score.String)
	var winners int
	require.NoError(t, db.QueryRowContext(ctx, "SELECT COUNT(*) FROM raffle_winner rw JOIN contender c ON c.id = rw.contender_id JOIN raffle r ON r.id = rw.raffle_id AND r.competition_id = c.competition_id").Scan(&winners))
	assert.Equal(t, 2, winners)

	params := database.UpsertCompetitionParams{
		ID:          domain.CompetitionID(uuid.New()),
		OrganizerID: 1,
		Name:        "Created with UUID",
		Country:     "SE",
		Created:     time.Now(),
	}
	require.NoError(t, queries.UpsertCompetition(ctx, params))
	params.Name = "Updated with UUID"
	require.NoError(t, queries.UpsertCompetition(ctx, params))
	created, err := queries.GetCompetition(ctx, params.ID)
	require.NoError(t, err)
	assert.Equal(t, params.ID, created.Competition.ID)
	assert.Equal(t, params.Name, created.Competition.Name)
	_, err = db.ExecContext(ctx, "INSERT INTO contender (organizer_id, competition_id, registration_code, class_id) VALUES (1, ?, 'INVALID', 2)", contender.Competitor.CompetitionID)
	require.Error(t, err)
	var checks int
	require.NoError(t, db.QueryRowContext(ctx, "SELECT @@SESSION.FOREIGN_KEY_CHECKS").Scan(&checks))
	assert.Equal(t, 1, checks)
}
