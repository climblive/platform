package repository

import (
	"database/sql"

	"github.com/go-errors/errors"
)

type transaction struct {
	tx *sql.Tx
}

func (tx *transaction) Commit() error {
	err := tx.tx.Commit()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (tx *transaction) Rollback() {
	_ = tx.tx.Rollback()
}
