package repository

import (
	"fmt"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(username, password, host string, port int, database string) (*Database, error) {
	var db *gorm.DB

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database)

	var logLevel logger.LogLevel = logger.Warn

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logLevel),
		TranslateError: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{
		db: db,
	}, nil
}

func (d *Database) Begin() domain.Transaction {
	tx := d.db.Begin()

	return &transaction{
		db: tx,
	}
}

func (d *Database) tx(tx domain.Transaction) *gorm.DB {
	transaction, ok := tx.(*transaction)
	if ok {
		return transaction.db
	} else {
		return d.db
	}
}
