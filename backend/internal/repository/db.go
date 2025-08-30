package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db      *sql.DB
	queries *database.Queries
}

func NewDatabase(username, password, host string, port int, databaseName string) (*Database, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		databaseName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	db = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter)
	db.Ping()

	queries := database.New(db)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	return &Database{
		db:      db,
		queries: queries,
	}, nil
}

func (d *Database) Begin() (domain.Transaction, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	return &transaction{
		tx: tx,
	}, nil
}

func (d *Database) WithTx(tx domain.Transaction) *database.Queries {
	transaction, ok := tx.(*transaction)
	if ok {
		return d.queries.WithTx(transaction.tx)
	} else {
		return d.queries
	}
}
