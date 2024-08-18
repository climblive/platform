package repository

import (
	"gorm.io/gorm"
)

type transaction struct {
	db *gorm.DB
}

func (tx *transaction) Commit() error {
	return tx.db.Commit().Error
}

func (tx *transaction) Rollback() {
	tx.db.Rollback()
}
