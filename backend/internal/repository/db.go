package repository

import (
	"database/sql"
	"fmt"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"

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

	queries := database.New(db)

	//	var logLevel logger.LogLevel = logger.Warn
	//
	//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//		Logger:         logger.Default.LogMode(logLevel),
	//		TranslateError: true,
	//	})
	//	if err != nil {
	//		return nil, errors.Wrap(err, 0)
	//	}
	//
	//	sqlDB, _ := db.DB()
	//
	//	sqlDB.SetMaxIdleConns(10)
	//	sqlDB.SetMaxOpenConns(100)
	//	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{
		db:      db,
		queries: queries,
	}, nil
}

func (d *Database) Begin() domain.Transaction {
	tx, err := d.db.Begin()
	if err != nil {
		panic("not handled")
	}

	return &transaction{
		tx: tx,
	}
}

func (d *Database) WithTx(tx domain.Transaction) *database.Queries {
	transaction, ok := tx.(*transaction)
	if ok {
		return d.queries.WithTx(transaction.tx)
	} else {
		return d.queries
	}
}
