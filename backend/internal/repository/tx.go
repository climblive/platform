package repository

import (
	"gorm.io/gorm"
)

type transaction struct {
	db     *gorm.DB
	active bool
}

func (tx *transaction) Commit() error {
	tx.active = false
	return tx.db.Commit().Error
}

func (tx *transaction) Rollback() {
	if !tx.active {
		return
	}

	tx.active = false
	tx.db.Rollback()
}
