package repository

import (
	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type transaction struct {
	db *gorm.DB
}

func (tx *transaction) Commit() error {
	return errors.Wrap(tx.db.Commit().Error, 0)
}

func (tx *transaction) Rollback() {
	tx.db.Rollback()
}
