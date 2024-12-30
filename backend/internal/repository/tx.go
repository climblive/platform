package repository

import (
	"database/sql"

	"github.com/go-errors/errors"
)

type transaction struct {
	tx *sql.Tx
}

func (tx *transaction) Commit() error {
	return errors.Wrap(tx.tx.Commit().Error, 0)
}

func (tx *transaction) Rollback() {
	_ = tx.tx.Rollback()
}
