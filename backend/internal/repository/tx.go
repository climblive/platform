package repository

import (
	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type transaction struct {
	db     *gorm.DB
	active bool
}

func (tx *transaction) Commit() error {
	tx.active = false

	err := tx.db.Commit().Error
	return errors.New(err)
}

func (tx *transaction) Rollback() {
	if !tx.active {
		return
	}

	tx.active = false
	tx.db.Rollback()
}
